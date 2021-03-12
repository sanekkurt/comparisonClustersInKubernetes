package secret

import (
	"context"

	"k8s-cluster-comparator/internal/config"
)

const (
	defaultSecretBatchLimit = 25
)

func getBatchLimit(ctx context.Context) int64 {
	cfg := config.FromContext(ctx)

	if limit := cfg.Configs.Secrets.BatchSize; limit != 0 {
		return limit
	}

	if limit := cfg.Common.DefaultBatchSize; limit != 0 {
		return limit
	}

	return defaultSecretBatchLimit
}

