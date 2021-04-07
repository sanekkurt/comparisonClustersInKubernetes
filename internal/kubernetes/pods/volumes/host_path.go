package volumes

import (
	"context"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumeHostPath(ctx context.Context, hostPath1, hostPath2 *v1.HostPathVolumeSource) {
	var (
		diffsBatch = diff.BatchFromContext(ctx)
	)

	if hostPath1.Path != hostPath2.Path {
		diffsBatch.Add(ctx, false, "%w: '%s' vs '%s'", ErrorHostPath, hostPath1.Path, hostPath2.Path)
	}

	if hostPath1.Type != nil && hostPath2.Type != nil {
		if *hostPath1.Type != *hostPath2.Type {
			diffsBatch.Add(ctx, false, "%w: '%s' vs '%s'", ErrorHostPathType, *hostPath1.Type, *hostPath2.Type)
		}
	} else if hostPath1.Type != nil || hostPath2.Type != nil {
		diffsBatch.Add(ctx, false, "%w", ErrorMissingHostPathType)
	}

	return
}
