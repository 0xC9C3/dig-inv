package log

import (
	"dig-inv/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc/grpclog"
)

var L *zap.Logger
var S *zap.SugaredLogger

func init() {
	err := initLogger()

	if err != nil {
		panic(err)
	}
	S = L.Sugar()

	// Ensure the logger is flushed before the program exits
	/*defer func() {
		if err := logger.Sync(); err != nil {
			sugar.Errorw("Error syncing logger", "error", err)
		}
	}*/

	grpclog.SetLoggerV2(
		zapgrpc.NewLogger(L),
	)
}

func initLogger() error {
	var err error

	if env.GetIsDevelopmentMode() {
		L, err = zap.NewDevelopment()
		return err
	}

	L, err = zap.NewProduction()
	return err
}
