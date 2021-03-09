package statefulset

import (
	"context"

	"k8s-cluster-comparator/internal/config"
)

func getBatchLimit(ctx context.Context) int64 {
	cfg := config.FromContext(ctx)

	if limit := cfg.Workloads.PodControllers.StatefulSets.BatchSize; limit != 0 {
		return limit
	}

	if limit := cfg.Common.DefaultBatchSize; limit != 0 {
		return limit
	}

	return 25
}
