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
