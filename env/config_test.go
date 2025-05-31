package env

import (
	"os"
	"testing"
)

func setEnvDeferrable(
	t *testing.T,
	key, value string,
) func() {
	err := os.Setenv(key, value)
	if err != nil {
		t.Fatalf("Failed to set environment variable '%s': %v", key, err)
	}

	return func() {
		err := os.Unsetenv(key)
		if err != nil {
			t.Errorf("Failed to unset environment variable '%s': %v", key, err)
		}
	}
}

func sanityTestEnvVariables(
	t *testing.T,
	envGetter func() string,
	key,
	expectedValue string,
) {
	defer setEnvDeferrable(t, key, expectedValue)()

	envValue := envGetter()
	if envValue != expectedValue {
		t.Errorf("Expected value '%s', got '%s'", expectedValue, envValue)
	}
}

type EnvironmentVariableTest struct {
	Key           string
	Getter        func() string
	ExpectedValue string
}

func TestEnvironmentVariables(t *testing.T) {
	environmentMapping := []EnvironmentVariableTest{
		{"PORT", GetPort, "8080"},
		{"LISTEN_ADDRESS", GetListenAddress, "1.2.2.4"},
		{"DEVELOPMENT", func() string {
			if GetIsDevelopmentMode() {
				return "true"
			}
			return "false"
		}, "true"},
		{"OIDC_CLIENT_ID", GetOidcClientID, "test-client-id"},
		{"OIDC_CLIENT_SECRET", GetOidcClientSecret, "test-client-secret"},
		{"OIDC_REDIRECT_URL", GetOidcRedirectURL, "http://localhost:8080/auth/callback"},
		{"OIDC_ISSUER_URL", GetOidcIssuerURL, "https://dig-inv.localhost"},
		{"OIDC_SCOPES", func() string {
			scopes := GetOidcScopes()
			if len(scopes) == 0 {
				return ""
			}
			return scopes[0]
		}, "openid profile email offline_access"},
	}

	for _, envTest := range environmentMapping {
		t.Run(envTest.Key, func(t *testing.T) {
			sanityTestEnvVariables(
				t,
				envTest.Getter,
				envTest.Key,
				envTest.ExpectedValue,
			)
		})
	}
}

func TestGetDefaults(t *testing.T) {
	defer setEnvDeferrable(t, "PORT", "")()

	defaultPort := GetPort()
	if defaultPort != "8080" {
		t.Errorf("Expected default port '8080', got '%s'", defaultPort)
	}
}

func TestAlternativeDevelopmentMode(t *testing.T) {
	defer setEnvDeferrable(t, "DEVELOPMENT", "foo")()

	developmentMode := GetIsDevelopmentMode()
	if developmentMode {
		t.Error("Expected development mode to be false when set to 'foo'")
	}
}
