package diff

import (
	"context"

	"go.uber.org/zap/zapcore"
	"k8s-cluster-comparator/internal/logging"
)

func diffLog(ctx context.Context, o object, logLevel zapcore.Level, msg string, variables ...interface{}) {
	var (
		log = logging.FromContext(ctx).With(
			"kind", o.Type.Kind,
			"objectName", o.Meta.Name,
		)
	)

	switch logLevel {
	case zapcore.WarnLevel:
		log.Warnf(msg, variables...)
	case zapcore.ErrorLevel:
		log.Errorf(msg, variables...)
	case zapcore.FatalLevel:
		log.Fatalf(msg, variables...)
	case zapcore.PanicLevel:
		log.Panicf(msg, variables...)
	}
}
