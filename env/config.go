package env

import "os"

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func GetPort() string {
	return getEnv("PORT", "8080")
}

func GetListenAddress() string {
	return getEnv("LISTEN_ADDRESS", "0.0.0.0")
}

func GetIsDevelopmentMode() bool {
	value := getEnv("DEVELOPMENT", "false")

	if value == "true" || value == "1" {
		return true
	}

	return false
}

func getOidcEnv(key, defaultValue string) string {
	value := os.Getenv("OIDC_" + key)
	if value == "" {
		return defaultValue
	}

	return value
}

func GetOidcClientID() string {
	return getOidcEnv("CLIENT_ID", "")
}

func GetOidcClientSecret() string {
	return getOidcEnv("CLIENT_SECRET", "")
}

func GetOidcRedirectURL() string {
	return getOidcEnv("REDIRECT_URL", "http://localhost:8080/auth/callback")
}

func GetOidcIssuerURL() string {
	return getOidcEnv("ISSUER_URL", "https://dig-inv.localhost")
}

func GetOidcScopes() []string {
	scopes := getOidcEnv("SCOPES", "openid profile email offline_access")
	if scopes == "" {
		return []string{}
	}

	return []string{scopes}
}
