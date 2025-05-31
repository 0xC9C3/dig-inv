package main

import (
	"dig-inv/cli"
	"os"
	"testing"
)

func TestEntrypoint_RunSuccess(t *testing.T) {
	mockCommandlineArgs(t, func(t *testing.T) {
		entrypoint := NewEntrypoint(
			func() error {
				return nil
			},
			func() error {
				return nil
			},
		)

		entrypoint.Run()
	}, cli.CommandServer)
}

func TestEntrypoint_RunWithError(t *testing.T) {
	mockCommandlineArgs(t, func(t *testing.T) {
		entrypoint := NewEntrypoint(
			func() error {
				return os.ErrPermission
			},
			func() error {
				return nil
			},
		)

		exitCode := entrypoint.Run()
		if exitCode != ErrorExitCode {
			t.Errorf("Expected exit code %d, got %d", ErrorExitCode, exitCode)
		}
	}, cli.CommandServer)
}

func mockCommandlineArgs(t *testing.T, inner func(t *testing.T), args ...string) {
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	os.Args = append([]string{"dig-inv"}, args...)

	inner(t)
}

func TestWorker(t *testing.T) {
	err := worker()
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}

func Test_Run(t *testing.T) {
	mockCommandlineArgs(t, func(t *testing.T) {
		run()
	}, cli.CommandWorker)
}

func Test_RunWithError(t *testing.T) {
	if err := os.Setenv("PORT", "invalid"); err != nil {
		t.Fatalf("Failed to set PORT environment variable: %v", err)
	}

	mockCommandlineArgs(t, func(t *testing.T) {
		exitCode := run()
		if exitCode != ErrorExitCode {
			t.Errorf("Expected exit code %d, got %d", ErrorExitCode, exitCode)
		}
	}, cli.CommandServer)
}
