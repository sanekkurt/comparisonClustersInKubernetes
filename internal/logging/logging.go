package logging

import (
	"context"

	"go.uber.org/zap"
)

type loggingCtxKey struct{}

var (
	log *zap.SugaredLogger
)

func Configure(debugMode bool) error { //nolint
	var err error
	var logger *zap.Logger

	logLevel := zap.InfoLevel
	if debugMode {
		logLevel = zap.DebugLevel
	}

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(logLevel)

	if logger, err = zapConfig.Build(); err != nil {
		log = zap.NewNop().Sugar()
	} else {
		log = logger.Sugar()
	}

	return nil
}

//func injectLoggerToContext(ctx context.Context) context.Context {
//	return context.WithValue(ctx, loggingCtxKey{}, Log)
//}

//func LoggerFromContext(ctx context.Context) *zap.SugaredLogger {
//	return ctx.Value(loggingCtxKey{}).(*zap.SugaredLogger)
//}

func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggingCtxKey{}, logger)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(loggingCtxKey{}).(*zap.SugaredLogger); ok {
		return logger
	}
	return log
}
