package main

import (
	"go.uber.org/zap"
)

var (
	log    *zap.SugaredLogger
)

func SetupLogging() error { //nolint
	var err error
	var logger *zap.Logger

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	if logger, err = zapConfig.Build(); err != nil {
		log = zap.NewNop().Sugar()
	}

	log = logger.Sugar()

	zap.ReplaceGlobals(logger)

	return nil
}
