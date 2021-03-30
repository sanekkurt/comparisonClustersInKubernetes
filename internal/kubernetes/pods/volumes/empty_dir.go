package volumes

import (
	"context"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumeEmptyDir(ctx context.Context, emptyDir1, emptyDir2 *v1.EmptyDirVolumeSource) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx)

		//diffs = make([]types.KubeObjectsDifference, 0)
	)

	if emptyDir1.Medium != emptyDir2.Medium {
		log.Warnf("%s. %s vs %s", ErrorVolumeEmptyDirMedium.Error(), emptyDir1.Medium, emptyDir2.Medium)
	}

	if emptyDir1.SizeLimit != nil && emptyDir2.SizeLimit != nil {

		if emptyDir1.SizeLimit.Format != emptyDir2.SizeLimit.Format {
			log.Warnf("%s. %s vs %s", ErrorVolumeEmptyDirSizeLimitFormat.Error(), emptyDir1.SizeLimit.Format, emptyDir2.SizeLimit.Format)
		}

	} else if emptyDir1.SizeLimit != nil || emptyDir2.SizeLimit != nil {
		log.Warnf("%s", ErrorVolumeEmptyDirSizeLimit.Error())
	}

	return nil, nil
}
