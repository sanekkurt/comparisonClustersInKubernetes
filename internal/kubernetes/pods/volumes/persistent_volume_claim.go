package volumes

import (
	"context"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumePersistentVolumeClaim(ctx context.Context, persistentVolumeClaim1, persistentVolumeClaim2 *v1.PersistentVolumeClaimVolumeSource) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx)

		//diffs = make([]types.KubeObjectsDifference, 0)
	)

	if persistentVolumeClaim1.ReadOnly != persistentVolumeClaim2.ReadOnly {
		log.Warnf("%s. %t vs %t", ErrorPersistentVolumeClaimReadOnly.Error(), persistentVolumeClaim1.ReadOnly, persistentVolumeClaim2.ReadOnly)
	}

	if persistentVolumeClaim1.ClaimName != persistentVolumeClaim2.ClaimName {
		log.Warnf("%s. %s vs %s", ErrorPersistentVolumeClaimName.Error(), persistentVolumeClaim1.ClaimName, persistentVolumeClaim2.ClaimName)
	}

	return nil, nil
}
