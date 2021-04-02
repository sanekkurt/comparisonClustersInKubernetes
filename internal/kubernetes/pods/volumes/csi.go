package volumes

import (
	"context"
	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
	"reflect"
)

func CompareVolumeCSI(ctx context.Context, csi1, csi2 *v1.CSIVolumeSource) {
	var (
		diffsBatch = diff.BatchFromContext(ctx)
	)

	if csi1.Driver != csi2.Driver {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %s vs %s", ErrorVolumeCSIDriver.Error(), csi1.Driver, csi2.Driver)
	}

	if csi1.ReadOnly != nil && csi2.ReadOnly != nil {

		if *csi1.ReadOnly != *csi2.ReadOnly {
			diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %t vs %t", ErrorVolumeCSIReadOnly.Error(), *csi1.ReadOnly, *csi2.ReadOnly)
		}

	} else if csi1.ReadOnly != nil || csi2.ReadOnly != nil {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s", ErrorMissingCSIReadOnly.Error())
	}

	if csi1.FSType != nil && csi2.FSType != nil {

		if *csi1.FSType != *csi2.FSType {
			diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %s vs %s", ErrorVolumeCSIFSType.Error(), *csi1.FSType, *csi2.FSType)
		}

	} else if csi1.FSType != nil || csi2.FSType != nil {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s", ErrorMissingCSIFSType.Error())
	}

	if csi1.NodePublishSecretRef != nil && csi2.NodePublishSecretRef != nil {

		if csi1.NodePublishSecretRef.Name != csi2.NodePublishSecretRef.Name {
			diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %s vs %s", ErrorVolumeCSIName.Error(), csi1.NodePublishSecretRef.Name, csi2.NodePublishSecretRef.Name)
		}

	} else if csi1.NodePublishSecretRef != nil || csi2.NodePublishSecretRef != nil {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s", ErrorMissingCSINodePublishSecretRef.Error())
	}

	if !reflect.DeepEqual(csi1.VolumeAttributes, csi2.VolumeAttributes) {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s", ErrorVolumeCSIVolumeAttributes.Error())
	}

	return
}
