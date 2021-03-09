package logging

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggingCtxKey struct{}

var (
	log *zap.SugaredLogger
)

func Configure(debugMode bool) error { //nolint
	var (
		err    error
		logger *zap.Logger

		logLevel = zap.InfoLevel
	)

	zapConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey: "",
			//LevelKey:      "L",
			LevelKey: "level",
			//NameKey:       "N",
			CallerKey:   "",
			FunctionKey: zapcore.OmitKey,
			MessageKey:  "message",
			//StacktraceKey: "S",
			LineEnding:  zapcore.DefaultLineEnding,
			EncodeLevel: zapcore.CapitalLevelEncoder,
			//EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if debugMode {
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		zapConfig.Development = true
		zapConfig.EncoderConfig.CallerKey = "caller"
		zapConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	if logger, err = zapConfig.Build(); err != nil {
		return fmt.Errorf("cannot create logger: %w", err)
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
