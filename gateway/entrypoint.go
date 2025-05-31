package gateway

import (
	"context"
	"dig-inv/env"
	gw "dig-inv/gen/go"
	"dig-inv/log"
	"dig-inv/services"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/protobuf/proto"
	"net/http"
	"strings"
)

func Run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(httpCookieResponseModifier),
	)

	if err := gw.RegisterOpenIdAuthServiceHandlerServer(ctx, mux, services.NewOpenIdAuthServer()); err != nil {
		log.S.Errorw("Failed to register OpenID Auth service handler", "error", err)
		return fmt.Errorf("failed to register OpenID Auth service handler: %w", err)
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

	return http.ListenAndServe(
		fmt.Sprintf("%s:%s", env.GetListenAddress(), env.GetPort()),
		corsHandler,
	)
}

// map cookies headers from gRPC metadata to HTTP headers
func httpCookieResponseModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
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

func httpCookieRequestModifier(ctx context.Context, r *http.Request) context.Context {
	cookies := r.Cookies()
	if len(cookies) == 0 {
		return ctx
	}

	log.S.Debugw("Setting cookies from HTTP request", "cookies", cookies)

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		md = runtime.ServerMetadata{}
	}

	for _, cookie := range cookies {
		md.HeaderMD.Append("set-cookie", fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}

	return runtime.NewServerMetadataContext(ctx, md)
}
