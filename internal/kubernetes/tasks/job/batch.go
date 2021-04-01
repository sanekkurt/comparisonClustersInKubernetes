package job

import (
	"context"

	"k8s-cluster-comparator/internal/config"
)

const (
	defaultJobBatchLimit = 25
)

func getBatchLimit(ctx context.Context) int64 {
	cfg := config.FromContext(ctx)

	if limit := cfg.Tasks.Jobs.BatchSize; limit != 0 {
		return limit
	}

	if limit := cfg.Common.DefaultBatchSize; limit != 0 {
		return limit
	}

	return defaultJobBatchLimit
}