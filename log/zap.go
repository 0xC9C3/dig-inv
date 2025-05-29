package log

import (
	"dig-inv/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc/grpclog"
)

var logger *zap.Logger
var sugar *zap.SugaredLogger

func init() {
	err := initLogger()

	if err != nil {
		panic(err)
	}
	sugar = logger.Sugar()

	// Ensure the logger is flushed before the program exits
	/*defer func() {
		if err := logger.Sync(); err != nil {
			sugar.Errorw("Error syncing logger", "error", err)
		}
	}*/

	grpclog.SetLoggerV2(
		zapgrpc.NewLogger(logger),
	)
}

func initLogger() error {
	var err error

	if env.GetIsDevelopmentMode() {
		logger, err = zap.NewDevelopment()
		return err
	}

	logger, err = zap.NewProduction()
	return err
}

func L() *zap.Logger {
	return logger
}

func S() *zap.SugaredLogger {
	return sugar
}
