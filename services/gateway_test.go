package services

import (
	"context"
	gen "dig-inv/gen/go"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	ts := http.Server{
		Addr: "127.0.0.1:0",
	}

	gwServer := Server{
		server: &ts,
	}

	serviceRunning := make(chan struct{})

	go func() {
		close(serviceRunning)
		if err := gwServer.Run(); err != nil {
			t.Errorf("Run failed: %v", err)
		}
	}()

	<-serviceRunning

	server, err := gwServer.GetServer()
	if err != nil {
		t.Errorf("GetServer failed: %v", err)
		return
	}

	if server == nil {
		t.Error("GetServer returned nil server")
		return
	}

	if err = server.Shutdown(context.Background()); err != nil {
		t.Errorf("Server shutdown failed: %v", err)
		return
	}
}

func TestRunWithInitializerError(t *testing.T) {
	gwServer := Server{
		initializer: []func(ctx context.Context, mux *runtime.ServeMux) error{
			func(ctx context.Context, mux *runtime.ServeMux) error {
				return fmt.Errorf("test error")
			},
		},
		server: nil,
	}

	// https://stackoverflow.com/a/55710364
	serviceRunning := make(chan struct{})
	serviceDone := make(chan struct{})

	go func() {
		close(serviceRunning)
		err := gwServer.Run()

		if err == nil || !strings.Contains(err.Error(), "failed to initialize handler") {
			t.Errorf("Run did not return expected error: %v", err)
		}

		close(serviceDone)
	}()

	<-serviceRunning
	<-serviceDone
}

func TestRunWithListenError(t *testing.T) {
	ts := http.Server{
		Addr: "thisdoesnot:exist:0:1234",
	}

	gwServer := Server{
		server: &ts,
	}

	serviceRunning := make(chan struct{})
	serviceDone := make(chan struct{})

	go func() {
		close(serviceRunning)
		err := gwServer.Run()
		if err == nil || !strings.Contains(err.Error(), "failed to start HTTP server") {
			t.Errorf("Run did not return expected error: %v", err)
		}
		close(serviceDone)
	}()

	<-serviceRunning
	<-serviceDone
}

func TestGetServer(t *testing.T) {
	_, err := NewGatewayServer().GetServer()
	if err != nil {
		t.Errorf("GetServer failed: %v", err)
		return
	}
}

func TestHttpCookieResponseModifier(t *testing.T) {
	ctx := context.Background()
	w := httptest.NewRecorder()

	ctx = AddMockServerTransportStreamToContext(ctx, "set-cookie", "test_cookie=test_value")

	err := httpCookieResponseModifier(ctx, w, &gen.UserSubjectMessage{})
	if err != nil {
		t.Errorf("httpCookieResponseModifier failed: %v", err)
	} else {
		t.Log("httpCookieResponseModifier executed successfully")
	}
}

func TestHttpCookieResponseModifierNoContext(t *testing.T) {
	ctx := context.Background()
	w := httptest.NewRecorder()

	err := httpCookieResponseModifier(ctx, w, &gen.UserSubjectMessage{})
	if err != nil {
		t.Errorf("httpCookieResponseModifier failed: %v", err)
	} else {
		t.Log("httpCookieResponseModifier executed successfully")
	}
}

func TestHttpCookieResponseModifierBrokenCookie(t *testing.T) {
	ctx := context.Background()
	w := httptest.NewRecorder()

	ctx = AddMockServerTransportStreamToContext(ctx, "set-cookie", "broken_cookie")

	err := httpCookieResponseModifier(ctx, w, &gen.UserSubjectMessage{})
	if err != nil {
		t.Errorf("httpCookieResponseModifier failed: %v", err)
	} else {
		t.Log("httpCookieResponseModifier executed successfully with broken cookie")
	}
}

func TestInitializeHandler(t *testing.T) {
	ctx := context.Background()

	handler, err := initializeHandler(ctx, []func(ctx context.Context, mux *runtime.ServeMux) error{
		func(ctx context.Context, mux *runtime.ServeMux) error {
			return nil
		},
	})
	if err != nil {
		t.Errorf("initializeHandler failed: %v", err)
		return
	}

	if handler == nil {
		t.Error("initializeHandler returned nil handler")
		return
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status OK, got %d", w.Code)
	}
}

func TestInitializeHandlerWithError(t *testing.T) {
	ctx := context.Background()

	_, err := initializeHandler(ctx, []func(ctx context.Context, mux *runtime.ServeMux) error{
		func(ctx context.Context, mux *runtime.ServeMux) error {
			return fmt.Errorf("test error")
		},
	})
	if err == nil {
		t.Error("Expected error, got nil")
	} else {
		t.Logf("initializeHandler returned expected error: %v", err)
	}
}
