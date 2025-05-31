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
	newLoggerInitializer().setupLogging()
}

type loggerInitializer struct {
	instanceInitializer func() error
}

func newLoggerInitializer() *loggerInitializer {
	return &loggerInitializer{
		instanceInitializer: zapInit,
	}
}

func (l *loggerInitializer) setupLogging() {
	if err := l.instanceInitializer(); err != nil {
		panic(err)
	}

	S = L.Sugar()

	// Ensure the logger is flushed before the program exits @todo
	/*defer func() {
		if err := logger.Sync(); err != nil {
			sugar.Errorw("Error syncing logger", "error", err)
		}
	}*/

	grpclog.SetLoggerV2(
		zapgrpc.NewLogger(L),
	)
}

func zapInit() error {
	var err error

	if env.GetIsDevelopmentMode() {
		L, err = zap.NewDevelopment()
		return err
	}

	L, err = zap.NewProduction()
	return err
}
