package cli

import (
	"errors"
	"os"
	"testing"
)

func TestCLI_RunServer(t *testing.T) {
	mockCommandlineArgs(
		t,
		cliRun,
		CommandServer,
	)
}

func TestCLI_RunWorker(t *testing.T) {
	mockCommandlineArgs(
		t,
		cliRun,
		CommandWorker,
	)
}

func TestCLI_RunWithError(t *testing.T) {
	mockCommandlineArgs(
		t,
		func(t *testing.T) {
			cli := NewCLI(nil, nil)
			if err := cli.Run(); err != nil {
				t.Errorf("Expected nil error, got %v", err)
			}

			cli.ServerHandler = func() error { return os.ErrInvalid }
			if err := cli.Run(); !errors.Is(err, os.ErrInvalid) {
				t.Errorf("Expected error, got nil")
			}
		},
		CommandServer,
	)
}

func cliRun(t *testing.T) {
	cli := NewCLI(nil, nil)
	if err := cli.Run(); err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// Test with server handler
	cli.ServerHandler = func() error { return nil }
	if err := cli.Run(); err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// Test with worker handler
	cli.WorkerHandler = func() error { return nil }
	if err := cli.Run(); err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}

func TestNewCLI(t *testing.T) {
	cli := NewCLI(nil, nil)
	if cli == nil {
		t.Error("NewCLI returned nil")
	}
	if cli != nil && (cli.ServerHandler != nil || cli.WorkerHandler != nil) {
		t.Error("NewCLI did not initialize handlers to nil")
	}
}

func mockCommandlineArgs(t *testing.T, inner func(t *testing.T), args ...string) {
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	os.Args = append([]string{"dig-inv"}, args...)

	inner(t)
}
