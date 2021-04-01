package volumes

import (
	"context"
	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumeHostPath(ctx context.Context, hostPath1, hostPath2 *v1.HostPathVolumeSource) {
	var (
		diffsBatch = diff.DiffBatchFromContext(ctx)
	)

	if hostPath1.Path != hostPath2.Path {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s: %s vs %s", ErrorHostPath.Error(), hostPath1.Path, hostPath2.Path)
	}

	if hostPath1.Type != nil && hostPath2.Type != nil {
		if *hostPath1.Type != *hostPath2.Type {
			diffsBatch.Add(ctx, false, zap.WarnLevel, "%s: %s vs %s", ErrorHostPathType.Error(), *hostPath1.Type, *hostPath2.Type)
		}
	} else if hostPath1.Type != nil || hostPath2.Type != nil {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s", ErrorMissingHostPathType.Error())
	}

	return
}
