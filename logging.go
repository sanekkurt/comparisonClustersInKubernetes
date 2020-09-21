package main

import (
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	log    *zap.SugaredLogger
)

func SetupLogging() error {
	var err error

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	if logger, err = zapConfig.Build(); err != nil {
		log = zap.NewNop().Sugar()
	}

	log = logger.Sugar()

	zap.ReplaceGlobals(logger)

	return nil
}
