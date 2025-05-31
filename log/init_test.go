package log

import (
	"os"
	"testing"
)

func TestInitializeLoggerDevelopment(t *testing.T) {
	if err := os.Setenv("DEVELOPMENT", "true"); err != nil {
		t.Fatalf("Failed to set DEVELOPMENT environment variable: %v", err)
	}

	newLoggerInitializer().setupLogging()

	if L == nil {
		t.Error("Logger L is nil")
	}

	if S == nil {
		t.Error("Logger S is nil")
	}
}

func TestInitializeLoggerProduction(t *testing.T) {
	if err := os.Setenv("DEVELOPMENT", "false"); err != nil {
		t.Fatalf("Failed to set DEVELOPMENT environment variable: %v", err)
	}

	newLoggerInitializer().setupLogging()

	if L == nil {
		t.Error("Logger L is nil in production mode")
	}

	if S == nil {
		t.Error("Logger S is nil in production mode")
	}
}

func TestInitLoggerError(t *testing.T) {
	li := &loggerInitializer{
		instanceInitializer: func() error {
			return os.ErrInvalid // Simulate an error
		},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic, but did not panic")
		}
	}()

	li.setupLogging()
}
