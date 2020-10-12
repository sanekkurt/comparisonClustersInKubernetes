package networking

import (
	"context"

	"go.uber.org/zap"

	"k8s-cluster-comparator/internal/logging"
)

var (
	log *zap.SugaredLogger
)

func Init(ctx context.Context) error {
	log = logging.FromContext(ctx)
	return nil
}
