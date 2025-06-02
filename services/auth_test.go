package services

import (
	"context"
	gw "dig-inv/gen/go"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"testing"
)

func setMockEnv(url string) {
	_ = os.Setenv("OIDC_CLIENT_ID", "test_client_id")
	_ = os.Setenv("OIDC_CLIENT_SECRET", "test_client_secret")
	_ = os.Setenv("OIDC_REDIRECT_URL", "http://testing.localhost/login")
	_ = os.Setenv("OIDC_ISSUER_URL", url)
	_ = os.Setenv("OIDC_SCOPES", "openid email profile offline_access")
}

func getMockOpenIdAuthServer() openidAuthServer {
	return openidAuthServer{
		getOAuth2ContextImpl: func(ctx context.Context) (*oauth2.Config, *oidc.Provider, string, error) {
			return &oauth2.Config{}, &oidc.Provider{}, "test_verifier", nil
		},
	}
}

func getMockOpenIdAuthServerError() openidAuthServer {
	return openidAuthServer{
		getOAuth2ContextImpl: func(ctx context.Context) (*oauth2.Config, *oidc.Provider, string, error) {
			return nil, nil, "", errors.New("test error")
		},
	}
}

func startMockOIDCServer(
	errorAction ...string,
) httptest.Server {
	return *httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(errorAction) > 0 && slices.Contains(errorAction, r.URL.Path) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			_, _ = fmt.Fprintf(w, `{
				"issuer": "http://%s",
				"authorization_endpoint": "http://%s/auth",
				"token_endpoint": "http://%s/token",
				"userinfo_endpoint": "http://%s/userinfo",
				"jwks_uri": "http://%s/jwks",
				"response_types_supported": ["code"],
				"subject_types_supported": ["public"],
				"id_token_signing_alg_values_supported": ["RS256"],
				"scopes_supported": ["openid", "email", "profile", "offline_access"],
				"token_endpoint_auth_methods_supported": ["client_secret_basic", "client_secret_post"],
				"claims_supported": ["sub", "email", "email_verified", "name", "given_name", "family_name"]
			}`, r.Host, r.Host, r.Host, r.Host, r.Host)

		case "/auth":
			http.Redirect(w, r, "https://mock-oidc-provider.com/callback?code=test_code&state=test_state", http.StatusFound)
		case "/token":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			resp := `{
				"access_token": "test_access_token",
				"token_type": "Bearer",
				"expires_in": 3600,
				"refresh_token": "test_refresh_token"
			}`

			_, _ = w.Write([]byte(resp))
		case "/userinfo":
			_, _ = w.Write([]byte(`{
				"sub": "test_subject",
				"email": "test@test.mail",
				"email_verified": true,
				"name": "Test User",
				"given_name": "Test",
				"family_name": "User",
				"picture": "http://example.com/test.jpg",
				"locale": "en-US"
			}`))
		}
	}))
}

func expectError(t *testing.T, err error) {
	if err == nil {
		t.Error("Expected an error, but got nil")
	} else {
		t.Logf("Received expected error: %v", err)
	}
}

func expectNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	} else {
		t.Log("No error as expected")
	}
}

func TestNewOpenIdAuthServer(t *testing.T) {
	server := NewOpenIdAuthServer()
	if server == nil {
		t.Error("NewOpenIdAuthServer returned nil")
	} else {
		t.Log("NewOpenIdAuthServer executed successfully")
	}
}

func TestOpenidAuthServer_BeginAuth(t *testing.T) {
	server := getMockOpenIdAuthServer()

	ctx := context.Background()
	ctx = AddMockServerTransportStreamToContext(ctx)

	authURL, err := server.BeginAuth(ctx, &gw.EmptyMessage{})
	expectNoError(t, err)

	if authURL == nil || authURL.Url == "" {
		t.Error("BeginAuth returned empty URL")
	} else {
		t.Logf("BeginAuth executed successfully, URL: %s", authURL.Url)
	}
}

func TestOpenidAuthServer_BeginAuthErr(t *testing.T) {
	server := getMockOpenIdAuthServerError()

	ctx := context.Background()
	ctx = AddMockServerTransportStreamToContext(ctx)

	_, err := server.BeginAuth(ctx, &gw.EmptyMessage{})
	expectError(t, err)
}

func TestOpenidAuthServer_BeginAuthErrSetCookies(t *testing.T) {
	server := getMockOpenIdAuthServer()

	ctx := context.Background()

	_, err := server.BeginAuth(ctx, &gw.EmptyMessage{})
	expectError(t, err)
}

func TestOpenidAuthServer_ExchangeCodeA(t *testing.T) {
	mockServer := startMockOIDCServer()

	fmt.Println("Mock OIDC server started at:", mockServer.URL)

	setMockEnv(mockServer.URL)
	server := NewOpenIdAuthServer()

	ctx := context.Background()
	ctx = AddMockServerTransportStreamToContext(ctx)
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("grpcgateway-cookie", "verifier=test_verifier; state=test_state"))

	msg := &gw.ExchangeCodeMessage{
		Code:  "test_code",
		State: "test_state",
	}

	res, err := server.ExchangeCode(ctx, msg)
	expectNoError(t, err)

	if res == nil {
		t.Error("ExchangeCode returned nil response")
	} else {
		t.Log("ExchangeCode executed successfully, response:", res)
	}

	// check if token, refresh token, verifier, and state cookies are set
	stream := grpc.ServerTransportStreamFromContext(ctx)

	if stream == nil {
		t.Error("Expected server transport stream in context, but got nil")
		return
	}

	// cast to MockServerTransportStream
	mockStream, ok := stream.(*MockServerTransportStream)
	if !ok {
		t.Error("Expected MockServerTransportStream, but got a different type")
		return
	}

	headers := mockStream.GetHeaders()
	if headers == nil {
		t.Error("Expected headers in MockServerTransportStream, but got nil")
		return
	}

	expectedValues := map[string]bool{
		fmt.Sprintf("%s=%s", TokenCookieName, "test_access_token"):         false,
		fmt.Sprintf("%s=%s", RefreshTokenCookieName, "test_refresh_token"): false,
		fmt.Sprintf("%s=", VerifierCookieName):                             false,
		fmt.Sprintf("%s=", StateCookieName):                                false,
	}

	for cookieName := range expectedValues {
		for _, header := range headers {
			if header.Get("set-cookie") == nil {
				continue
			}

			for _, cookie := range header.Get("set-cookie") {

				if cookie == cookieName {
					if expectedValues[cookieName] {
						t.Errorf("Unexpected cookie found: %s", cookieName)
						continue
					}

					expectedValues[cookieName] = true
					t.Logf("Found expected cookie: %s", cookieName)
					break
				}
			}
		}
	}

	for cookieName, expectedValue := range expectedValues {
		if !expectedValue {
			t.Errorf("Expected cookie not found: %s", cookieName)
		} else {
			t.Logf("Verified expected cookie: %s", cookieName)
		}
	}
}

func TestOpenidAuthServer_ExchangeCodeErr(t *testing.T) {
	server := getMockOpenIdAuthServerError()

	ctx := context.Background()
	ctx = AddMockServerTransportStreamToContext(ctx)
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("grpcgateway-cookie", "verifier=test_verifier; state=test_state"))

	msg := &gw.ExchangeCodeMessage{
		Code:  "test_code",
		State: "test_state",
	}

	_, err := server.ExchangeCode(ctx, msg)
	expectError(t, err)
}

func TestOpenidAuthServer_ExchangeCodeErrVerifier(t *testing.T) {
	server := getMockOpenIdAuthServer()

	ctx := context.Background()
	ctx = AddMockServerTransportStreamToContext(ctx)

	msg := &gw.ExchangeCodeMessage{
		Code:  "test_code",
		State: "test_state",
	}

	_, err := server.ExchangeCode(ctx, msg)
	expectError(t, err)
}

func TestOpenidAuthServer_ExchangeCodeErrState(t *testing.T) {
	server := getMockOpenIdAuthServer()

	ctx := context.Background()
	ctx = AddMockServerTransportStreamToContext(ctx)
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("grpcgateway-cookie", "verifier=test_verifier"))

	msg := &gw.ExchangeCodeMessage{
		Code:  "test_code",
		State: "test_state",
	}

	_, err := server.ExchangeCode(ctx, msg)
	expectError(t, err)
}

func TestOpenidAuthServer_ExchangeCodeErrStateVerify(t *testing.T) {
	server := getMockOpenIdAuthServer()

	ctx := context.Background()
	ctx = AddMockServerTransportStreamToContext(ctx)
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("grpcgateway-cookie", "verifier=test_verifier; state=test_state"))

	msg := &gw.ExchangeCodeMessage{
		Code:  "test_code",
		State: "test_state_wrong",
	}

	_, err := server.ExchangeCode(ctx, msg)
	expectError(t, err)
}

func TestOpenidAuthServer_ExchangeCodeServerError(t *testing.T) {
	mockServer := startMockOIDCServer("/token")

	setMockEnv(mockServer.URL)
	server := NewOpenIdAuthServer()
	ctx := context.Background()
	ctx = AddMockServerTransportStreamToContext(ctx)
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("grpcgateway-cookie", "verifier=test_verifier; state=test_state"))
	msg := &gw.ExchangeCodeMessage{
		Code:  "test_code",
		State: "test_state",
	}
	_, err := server.ExchangeCode(ctx, msg)
	expectError(t, err)
}

func TestOpenidAuthServer_ExchangeCodeServerSetTokenCookieError(t *testing.T) {
	mockServer := startMockOIDCServer()

	setMockEnv(mockServer.URL)
	server := NewOpenIdAuthServer()
	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("grpcgateway-cookie", "verifier=test_verifier; state=test_state"))
	msg := &gw.ExchangeCodeMessage{
		Code:  "test_code",
		State: "test_state",
	}
	_, err := server.ExchangeCode(ctx, msg)
	expectError(t, err)
}

func TestOpenidAuthServer_GetUserInfo(t *testing.T) {
	mockServer := startMockOIDCServer()

	setMockEnv(mockServer.URL)
	server := NewOpenIdAuthServer()
	ctx := context.Background()
	ctx = AddMockServerTransportStreamToContext(ctx)
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("grpcgateway-cookie", "token=test_access_token; verifier=test_verifier; state=test_state"))

	userInfo, err := server.GetUserInfo(ctx, &gw.EmptyMessage{})
	expectNoError(t, err)

	if userInfo == nil || userInfo.Subject == "" {
		t.Error("GetUserInfo returned empty user info")
	} else if userInfo.Subject != "test_subject" {
		t.Errorf("GetUserInfo returned unexpected subject: %s", userInfo.Subject)
	} else {
		t.Logf("GetUserInfo executed successfully, Subject: %s", userInfo.Subject)
	}
}

func TestOpenidAuthServer_GetUserInfoServerError(t *testing.T) {
	mockServer := startMockOIDCServer("/userinfo")

	setMockEnv(mockServer.URL)
	server := NewOpenIdAuthServer()
	ctx := context.Background()
	ctx = AddMockServerTransportStreamToContext(ctx)
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("grpcgateway-cookie", "token=test_access_token; verifier=test_verifier; state=test_state"))

	_, err := server.GetUserInfo(ctx, &gw.EmptyMessage{})
	expectError(t, err)
}

func TestOpenidAuthServer_GetUserInfoErr(t *testing.T) {
	server := getMockOpenIdAuthServerError()
	ctx := context.Background()
	_, err := server.GetUserInfo(ctx, &gw.EmptyMessage{})

	expectError(t, err)
}

func TestOpenidAuthServer_GetUserInfoErrAccessToken(t *testing.T) {
	server := getMockOpenIdAuthServer()
	ctx := context.Background()
	_, err := server.GetUserInfo(ctx, &gw.EmptyMessage{})

	expectError(t, err)
}

func TestOpenidAuthServer_Logout(t *testing.T) {
	server := getMockOpenIdAuthServer()

	ctx := context.Background()
	ctx = AddMockServerTransportStreamToContext(ctx)

	_, err := server.Logout(ctx, &gw.EmptyMessage{})

	expectNoError(t, err)

	// Check if cookies are cleared

}

func TestOpenidAuthServer_LogoutNoContext(t *testing.T) {
	server := getMockOpenIdAuthServer()

	ctx := context.Background()

	_, err := server.Logout(ctx, &gw.EmptyMessage{})
	expectError(t, err)
}

func TestOpenidAuthServer_GetOAuth2Context(t *testing.T) {
	ctx := context.Background()

	oauth2Config, provider, verifier, err := getOAuth2Context(ctx)

	expectNoError(t, err)

	if oauth2Config == nil || provider == nil || verifier == "" {
		t.Error("getOAuth2Context returned incomplete data")
	} else {
		t.Logf("getOAuth2Context executed successfully, Verifier: %s", verifier)
	}
}

func TestOpenidAuthServer_GetOAuth2ContextError(t *testing.T) {
	ctx := context.Background()

	_ = os.Setenv("OIDC_ISSUER_URL", "http://invalid-url")

	// Simulate an error by passing a nil context
	_, _, _, err := getOAuth2Context(ctx)
	expectError(t, err)
}

func TestOpenidAuthServer_OAuth2ConfigErrorResponse(t *testing.T) {
	err := oauth2ConfigErrorResponse(errors.New("test error"))
	expectError(t, err)
}

func TestOpenidAuthServer_SetCookie(t *testing.T) {
	ctx := context.Background()
	ctx = AddMockServerTransportStreamToContext(ctx)

	err := setCookie(ctx, "test_cookie", "test_value")
	expectNoError(t, err)
}

func TestOpenidAuthServer_GetCookie(t *testing.T) {
	ctx := context.Background()
	ctx = AddMockServerTransportStreamToContext(ctx, "test_cookie", "test_value")
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("grpcgateway-cookie", "test_cookie=test_value"))

	_, err := getCookie(ctx, "test_cookie")
	expectNoError(t, err)
}

func TestOpenidAuthServer_GetCookieNoMetadata(t *testing.T) {
	ctx := context.Background()
	ctx = AddMockServerTransportStreamToContext(ctx, "test_cookie", "test_value")
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("grpcgateway-cookie-ex", "test_cookie=test_value"))

	_, err := getCookie(ctx, "test_cookie")
	expectError(t, err)
}

func TestOpenidAuthServer_GetCookieBrokenCookie(t *testing.T) {
	ctx := context.Background()
	ctx = AddMockServerTransportStreamToContext(ctx, "test_cookie", "test_value")
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("grpcgateway-cookie", "12453#1  @⁴56"))
	_, err := getCookie(ctx, "test_cookie")
	expectError(t, err)
}

func Test_setCookiesEmpty(t *testing.T) {
	ctx := context.Background()
	err := setCookies(ctx)
	expectNoError(t, err)
}
