package volumes

import (
	"context"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/core/v1"
	"reflect"
)

func CompareVolumeCSI(ctx context.Context, csi1, csi2 *v1.CSIVolumeSource) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx)

		//diffs = make([]types.KubeObjectsDifference, 0)
	)

	if csi1.Driver != csi2.Driver {
		log.Warnf("%s. %s vs %s", ErrorVolumeCSIDriver.Error(), csi1.Driver, csi2.Driver)
	}

	if csi1.ReadOnly != nil && csi2.ReadOnly != nil {

		if *csi1.ReadOnly != *csi2.ReadOnly {
			log.Warnf("%s. %t vs %t", ErrorVolumeCSIReadOnly.Error(), *csi1.ReadOnly, *csi2.ReadOnly)
		}

	} else if csi1.ReadOnly != nil || csi2.ReadOnly != nil {
		log.Warnf("%s", ErrorMissingCSIReadOnly.Error())
	}

	if csi1.FSType != nil && csi2.FSType != nil {

		if *csi1.FSType != *csi2.FSType {
			log.Warnf("%s. %s vs %s", ErrorVolumeCSIFSType.Error(), *csi1.FSType, *csi2.FSType)
		}

	} else if csi1.FSType != nil || csi2.FSType != nil {
		log.Warnf("%s", ErrorMissingCSIFSType.Error())
	}

	if csi1.NodePublishSecretRef != nil && csi2.NodePublishSecretRef != nil {

		if csi1.NodePublishSecretRef.Name != csi2.NodePublishSecretRef.Name {
			log.Warnf("%s. %s vs %s", ErrorVolumeCSIName.Error(), csi1.NodePublishSecretRef.Name, csi2.NodePublishSecretRef.Name)
		}

	} else if csi1.NodePublishSecretRef != nil || csi2.NodePublishSecretRef != nil {
		log.Warnf("%s", ErrorMissingCSINodePublishSecretRef.Error())
	}

	if !reflect.DeepEqual(csi1.VolumeAttributes, csi2.VolumeAttributes) {
		log.Warnf("%s", ErrorVolumeCSIVolumeAttributes.Error())
	}

	return nil, nil
}
