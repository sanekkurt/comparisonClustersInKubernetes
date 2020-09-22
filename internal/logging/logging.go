package logging

import (
	"go.uber.org/zap"
)

var (
	Log    *zap.SugaredLogger
)

func SetupLogging() error { //nolint
	var err error
	var logger *zap.Logger

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	if logger, err = zapConfig.Build(); err != nil {
		Log = zap.NewNop().Sugar()
	}

	Log = logger.Sugar()

	zap.ReplaceGlobals(logger)

	return nil
}
