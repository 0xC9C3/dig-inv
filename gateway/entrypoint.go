package gateway

import (
	"context"
	"dig-inv/env"
	gw "dig-inv/gen/go"
	"dig-inv/services"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"net/http"
)

func Run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	err := gw.RegisterYourServiceHandlerServer(ctx, mux, services.NewEchoServer())

	if err != nil {
		return err
	}

	// @todo config
	corsHandler := cors.Default().Handler(mux)
	if env.GetIsDevelopmentMode() {
		corsHandler = cors.Default().Handler(mux)
	}

	return http.ListenAndServe(
		fmt.Sprintf("%s:%s", env.GetListenAddress(), env.GetPort()),
		corsHandler,
	)
}
