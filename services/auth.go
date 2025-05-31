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

type openidAuthServer struct {
	gw.UnimplementedOpenIdAuthServiceServer
}

func NewOpenIdAuthServer() gw.OpenIdAuthServiceServer {
	return new(openidAuthServer)
}

func (s *openidAuthServer) GetUserInfo(ctx context.Context, msg *gw.EmptyMessage) (*gw.UserSubjectMessage, error) {
	_, provider, _, err := getOAuth2Context(ctx)

	if err != nil {
		return nil, oauth2ConfigErrorResponse(err)
	}

	accessToken, err := getCookie(ctx, "token")
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
	oauth2Config, _, _, err := getOAuth2Context(ctx)
	if err != nil {
		return nil, oauth2ConfigErrorResponse(err)
	}

	verifier, err := getCookie(ctx, "verifier")
	if err != nil {
		grpclog.Errorf("Failed to get verifier cookie: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "failed to get verifier cookie")
	}

	state, err := getCookie(ctx, "state")
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

	if token.RefreshToken != "" {
		if err := setCookie(ctx, "refresh_token", token.RefreshToken); err != nil {
			grpclog.Errorf("Failed to set refresh token cookie: %v", err)
			return nil, status.Errorf(codes.Internal, "failed to set refresh token cookie")
		}
	}
	log.S.Debugw("Exchanged token", "access_token", token.AccessToken, "refresh_token", token.RefreshToken)

	if err := setCookie(ctx, "token", token.AccessToken); err != nil {
		grpclog.Errorf("Failed to set access token cookie: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to set access token cookie")
	}

	return &gw.EmptyMessage{}, nil
}

func (s *openidAuthServer) BeginAuth(ctx context.Context, _ *gw.EmptyMessage) (*gw.AuthUrlMessage, error) {
	oauth2Config, _, verifier, err := getOAuth2Context(ctx)

	if err != nil {
		return nil, oauth2ConfigErrorResponse(err)
	}

	state := make([]byte, 16)
	if _, err := rand.Read(state); err != nil {
		grpclog.Errorf("Failed to generate random state: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to generate random state")
	}

	stateString := base64.RawURLEncoding.EncodeToString(state)

	authURL := oauth2Config.AuthCodeURL(
		stateString,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
		oauth2.S256ChallengeOption(verifier),
	)

	grpclog.Infof("Generated auth URL: %s", authURL)

	if err := setCookie(ctx, "verifier", verifier); err != nil {
		return nil, err
	}

	if err := setCookie(ctx, "state", stateString); err != nil {
		return nil, err
	}

	return &gw.AuthUrlMessage{Url: authURL}, nil
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

	// parse header cookie string

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
