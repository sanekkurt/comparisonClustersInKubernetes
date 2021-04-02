package volumes

import (
	"context"
	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumeEmptyDir(ctx context.Context, emptyDir1, emptyDir2 *v1.EmptyDirVolumeSource) {
	var (
		diffsBatch = diff.BatchFromContext(ctx)
	)

	if emptyDir1.Medium != emptyDir2.Medium {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %s vs %s", ErrorVolumeEmptyDirMedium.Error(), emptyDir1.Medium, emptyDir2.Medium)
	}

	if emptyDir1.SizeLimit != nil && emptyDir2.SizeLimit != nil {

		if emptyDir1.SizeLimit.Format != emptyDir2.SizeLimit.Format {
			diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %s vs %s", ErrorVolumeEmptyDirSizeLimitFormat.Error(), emptyDir1.SizeLimit.Format, emptyDir2.SizeLimit.Format)
		}

	} else if emptyDir1.SizeLimit != nil || emptyDir2.SizeLimit != nil {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s", ErrorVolumeEmptyDirSizeLimit.Error())
	}

	return
}
