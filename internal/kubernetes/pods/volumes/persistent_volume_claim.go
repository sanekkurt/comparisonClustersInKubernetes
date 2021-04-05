package volumes

import (
	"context"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumePersistentVolumeClaim(ctx context.Context, persistentVolumeClaim1, persistentVolumeClaim2 *v1.PersistentVolumeClaimVolumeSource) {
	var (
		diffsBatch = diff.BatchFromContext(ctx)
	)

	if persistentVolumeClaim1.ReadOnly != persistentVolumeClaim2.ReadOnly {
		diffsBatch.Add(ctx, false, "%s. %t vs %t", ErrorPersistentVolumeClaimReadOnly.Error(), persistentVolumeClaim1.ReadOnly, persistentVolumeClaim2.ReadOnly)
	}

	if persistentVolumeClaim1.ClaimName != persistentVolumeClaim2.ClaimName {
		diffsBatch.Add(ctx, false, "%s. %s vs %s", ErrorPersistentVolumeClaimName.Error(), persistentVolumeClaim1.ClaimName, persistentVolumeClaim2.ClaimName)
	}

	return
}
