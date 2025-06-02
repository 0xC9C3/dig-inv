package services

import (
	"context"
	"crypto/rand"
	"dig-inv/env"
	gw "dig-inv/gen/go"
	"dig-inv/log"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
)

const (
	TokenCookieName        = "token"
	RefreshTokenCookieName = "refresh_token"
	VerifierCookieName     = "verifier"
	StateCookieName        = "state"
)

type openidAuthServer struct {
	gw.UnimplementedOpenIdAuthServiceServer
	getOAuth2ContextImpl func(ctx context.Context) (*oauth2.Config, *oidc.Provider, string, error)
}

func NewOpenIdAuthServer() gw.OpenIdAuthServiceServer {
	return &openidAuthServer{
		getOAuth2ContextImpl: getOAuth2Context,
	}
}

func (s *openidAuthServer) GetUserInfo(ctx context.Context, _ *gw.EmptyMessage) (*gw.UserSubjectMessage, error) {
	_, provider, _, err := s.getOAuth2ContextImpl(ctx)

	if err != nil {
		return nil, oauth2ConfigErrorResponse(err)
	}

	accessToken, err := getCookie(ctx, TokenCookieName)
	if err != nil {
		grpclog.Errorf("Failed to get access token cookie: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "failed to get access token cookie")
	}

	userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
	}))

	if err != nil {
		grpclog.Errorf("Failed to get user info: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "Failed to get user info")
	}

	return &gw.UserSubjectMessage{Subject: userInfo.Subject}, nil
}

func (s *openidAuthServer) ExchangeCode(ctx context.Context, msg *gw.ExchangeCodeMessage) (*gw.EmptyMessage, error) {
	oauth2Config, _, _, err := s.getOAuth2ContextImpl(ctx)
	if err != nil {
		return nil, oauth2ConfigErrorResponse(err)
	}

	verifier, err := getCookie(ctx, VerifierCookieName)
	if err != nil {
		grpclog.Errorf("Failed to get verifier cookie: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "failed to get verifier cookie")
	}

	state, err := getCookie(ctx, StateCookieName)
	if err != nil {
		grpclog.Errorf("Failed to get state cookie: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "failed to get state cookie")
	}

	if msg.State != state {
		grpclog.Errorf("State mismatch: expected %s, got %s", state, msg.State)
		return nil, status.Errorf(codes.Unauthenticated, "state mismatch")
	}

	log.S.Debugw("Exchanging code", "code", msg.Code, "state", msg.State, "verifier", verifier)

	token, err := oauth2Config.Exchange(ctx, msg.Code, oauth2.VerifierOption(verifier))
	if err != nil {
		grpclog.Errorf("Failed to exchange token: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "failed to exchange token")
	}

	if err := setCookies(ctx, []string{TokenCookieName, token.AccessToken},
		[]string{RefreshTokenCookieName, token.RefreshToken},
		[]string{VerifierCookieName, ""},
		[]string{StateCookieName, ""}); err != nil {
		grpclog.Errorf("Failed to set cookies: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to set cookies")
	}

	return &gw.EmptyMessage{}, nil
}

func (s *openidAuthServer) BeginAuth(ctx context.Context, _ *gw.EmptyMessage) (*gw.AuthUrlMessage, error) {
	oauth2Config, _, verifier, err := s.getOAuth2ContextImpl(ctx)

	if err != nil {
		return nil, oauth2ConfigErrorResponse(err)
	}

	state := generateState()

	authURL := oauth2Config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
		oauth2.S256ChallengeOption(verifier),
	)

	grpclog.Infof("Generated auth URL: %s", authURL)

	if err := setCookies(ctx, []string{VerifierCookieName, verifier}, []string{StateCookieName, state}); err != nil {
		grpclog.Errorf("Failed to set cookies: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to set cookies")
	}

	return &gw.AuthUrlMessage{Url: authURL}, nil
}

func (s *openidAuthServer) Logout(ctx context.Context, _ *gw.EmptyMessage) (*gw.EmptyMessage, error) {
	cookiesToClear := []string{TokenCookieName, RefreshTokenCookieName, VerifierCookieName, StateCookieName}

	for _, cookie := range cookiesToClear {
		if err := setCookie(ctx, cookie, ""); err != nil {
			grpclog.Errorf("Failed to clear cookie %s: %v", cookie, err)
			return nil, status.Errorf(codes.Internal, "failed to clear cookies")
		}
	}

	return &gw.EmptyMessage{}, nil
}

func getOAuth2Context(
	ctx context.Context,
) (*oauth2.Config, *oidc.Provider, string, error) {
	provider, err := oidc.NewProvider(ctx, env.GetOidcIssuerURL())
	if err != nil {
		grpclog.Errorf("Failed to create OIDC provider: %v", err)
		return nil, nil, "", errors.New("failed to create OIDC provider")
	}

	verifier := oauth2.GenerateVerifier()
	oauth2Config := oauth2.Config{
		ClientID:     env.GetOidcClientID(),
		ClientSecret: env.GetOidcClientSecret(),
		RedirectURL:  env.GetOidcRedirectURL(),
		Endpoint:     provider.Endpoint(),
		Scopes:       env.GetOidcScopes(),
	}

	return &oauth2Config, provider, verifier, nil
}

func oauth2ConfigErrorResponse(err error) error {
	grpclog.Errorf("Failed to get OAuth2 config: %v", err)
	return status.Errorf(codes.Internal, "failed to get OAuth2 config")
}

func setCookie(ctx context.Context, key string, value string) error {
	err := grpc.SetHeader(ctx, metadata.Pairs("set-cookie", fmt.Sprintf("%s=%s", key, value)))
	if err != nil {
		grpclog.Errorf("Failed to set header: %v", err)
		return status.Errorf(codes.Internal, "failed to set header")
	}

	return nil
}

func setCookies(ctx context.Context, cookies ...[]string) error {
	if len(cookies) == 0 {
		return nil
	}

	for _, cookie := range cookies {
		if err := setCookie(ctx, cookie[0], cookie[1]); err != nil {
			grpclog.Errorf("Failed to set cookie %s: %v", cookie[0], err)
			return status.Errorf(codes.Internal, "failed to set cookie %s", cookie[0])
		}
	}

	return nil
}

func getCookie(ctx context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Internal, "failed to get metadata from context")
	}

	cookies := md.Get("grpcgateway-cookie")
	if len(cookies) == 0 {
		grpclog.Errorf("Cookie not found: %s", key)
		return "", status.Errorf(codes.NotFound, "cookie not found")
	}

	log.S.Debugw("Retrieved cookies from metadata", "cookies", cookies)

	for _, cookie := range cookies {
		parsedCookies, err := http.ParseCookie(cookie)
		if err != nil {
			grpclog.Errorf("Failed to parse cookie: %v", err)
			return "", status.Errorf(codes.Internal, "failed to parse cookie")
		}

		for _, c := range parsedCookies {
			if c.Name == key {
				return c.Value, nil
			}
		}
	}

	grpclog.Errorf("Cookie not found: %s", key)
	return "", status.Errorf(codes.NotFound, "cookie not found")
}

func generateState() string {
	state := make([]byte, 16)

	// ignoring the error since
	//  "Read calls [io.ReadFull] on [Reader] and crashes the program irrecoverably if
	//  an error is returned."
	_, _ = rand.Read(state)
	return base64.RawURLEncoding.EncodeToString(state)
}
