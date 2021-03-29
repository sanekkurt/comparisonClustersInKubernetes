package volumes

import (
	"context"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumeHostPath(ctx context.Context, hostPath1, hostPath2 *v1.HostPathVolumeSource) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx)

		//diffs = make([]types.KubeObjectsDifference, 0)
	)

	if hostPath1.Path != hostPath2.Path {
		log.Warnf("%s: %s vs %s", ErrorHostPath.Error(), hostPath1.Path, hostPath2.Path)
	}

	if hostPath1.Type != nil && hostPath2.Type != nil {
		if *hostPath1.Type != *hostPath2.Type {
			log.Warnf("%s: %s vs %s", ErrorHostPathType.Error(), *hostPath1.Type, *hostPath2.Type)
		}
	} else if hostPath1.Type != nil || hostPath2.Type != nil {
		log.Warnf("%s", ErrorMissingHostPathType.Error())
	}

	return nil, nil
}
