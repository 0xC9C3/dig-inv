package services

import (
	"context"
	"dig-inv/env"
	gw "dig-inv/gen/go"
	"dig-inv/log"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/protobuf/proto"
	"net/http"
	"strings"
)

var serviceInitializer = []func(ctx context.Context, mux *runtime.ServeMux) error{
	func(ctx context.Context, mux *runtime.ServeMux) error {
		return gw.RegisterOpenIdAuthServiceHandlerServer(ctx, mux, NewOpenIdAuthServer())
	},
}

type Server struct {
	initializer []func(ctx context.Context, mux *runtime.ServeMux) error
	server      *http.Server
}

func (gateway *Server) Run() error {
	server, err := gateway.GetServer()
	if err != nil {
		log.S.Errorw("Failed to get server", "error", err)
		return fmt.Errorf("failed to get server: %w", err)
	}

	log.S.Infow("Starting HTTP server", "address", server.Addr)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.S.Errorw("Failed to start HTTP server", "error", err)
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

func (gateway *Server) GetServer() (*http.Server, error) {
	if gateway.server != nil {
		return gateway.server, nil
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	handler, err := initializeHandler(ctx, gateway.initializer)
	if err != nil {
		log.S.Errorw("Failed to initialize handler", "error", err)
		return nil, fmt.Errorf("failed to initialize handler: %w", err)
	}

	gateway.server = &http.Server{Addr: fmt.Sprintf("%s:%s", env.GetListenAddress(), env.GetPort()), Handler: handler}

	return gateway.server, nil
}

func initializeHandler(ctx context.Context, serviceInitializer []func(ctx context.Context, mux *runtime.ServeMux) error) (http.Handler, error) {
	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(httpCookieResponseModifier),
	)

	for _, initFunc := range serviceInitializer {
		if err := initFunc(ctx, mux); err != nil {
			log.S.Errorw("Failed to initialize service handler", "error", err)
			return nil, fmt.Errorf("failed to initialize service handler: %w", err)
		}
	}

	// @todo config
	corsHandler := cors.Default().Handler(mux)
	if env.GetIsDevelopmentMode() {
		corsOption := cors.Options{
			AllowedOrigins:   []string{"http://localhost:5173"},
			AllowCredentials: true,
		}
		corsHandler = cors.New(corsOption).Handler(mux)
	}

	return corsHandler, nil
}

// map cookie headers from gRPC metadata to HTTP headers
func httpCookieResponseModifier(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	if vals := md.HeaderMD.Get("set-cookie"); len(vals) > 0 {
		log.S.Debugw("Setting cookies from gRPC metadata", "cookies", vals)

		delete(md.HeaderMD, "set-cookie")
		for _, cookie := range vals {
			// split by = for name and value
			parts := strings.SplitN(cookie, "=", 2)
			if len(parts) != 2 {
				log.S.Errorw("Invalid cookie format", "cookie", cookie)
				continue
			}

			name := parts[0]
			value := parts[1]

			// set cookie in HTTP response
			http.SetCookie(w, &http.Cookie{
				Name:     name,
				Value:    value,
				HttpOnly: true,
				Secure:   !env.GetIsDevelopmentMode(),
			})
		}

		delete(w.Header(), "Grpc-Metadata-Set-Cookie")
		delete(md.HeaderMD, "set-cookie")
	}

	return nil
}

func NewGatewayServer() *Server {
	return &Server{
		initializer: serviceInitializer,
		server:      nil,
	}
}
