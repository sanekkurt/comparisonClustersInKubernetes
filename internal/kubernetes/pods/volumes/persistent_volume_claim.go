package volumes

import (
	"context"
	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumePersistentVolumeClaim(ctx context.Context, persistentVolumeClaim1, persistentVolumeClaim2 *v1.PersistentVolumeClaimVolumeSource) {
	var (
		diffsBatch = diff.DiffBatchFromContext(ctx)
	)

	if persistentVolumeClaim1.ReadOnly != persistentVolumeClaim2.ReadOnly {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %t vs %t", ErrorPersistentVolumeClaimReadOnly.Error(), persistentVolumeClaim1.ReadOnly, persistentVolumeClaim2.ReadOnly)
	}

	if persistentVolumeClaim1.ClaimName != persistentVolumeClaim2.ClaimName {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %s vs %s", ErrorPersistentVolumeClaimName.Error(), persistentVolumeClaim1.ClaimName, persistentVolumeClaim2.ClaimName)
	}

	return
}
