package containers

import (
	"context"
	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/core/v1"
)

func compareContainerVolumeMounts(ctx context.Context, container1, container2 v1.Container) ([]types.KubeObjectsDifference, error) {
	var (
		log   = logging.FromContext(ctx)
		diffs = make([]types.KubeObjectsDifference, 0)
	)

	log.Debugf("ComparePodSpecs: start checking volume mounts in container - %s", container1.Name)

	log = log.With(zap.String("containerName", container1.Name))

	if len(container1.VolumeMounts) != len(container2.VolumeMounts) {
		log.Warnf("%s", ErrorContainerVolumeMountsLength.Error())
	} else {
		for index, volumeMount := range container1.VolumeMounts {

			if volumeMount.Name != container2.VolumeMounts[index].Name {
				log.Warnf("%s. %s vs %s", ErrorContainerVolumeMountsName.Error(), volumeMount.Name, container2.VolumeMounts[index].Name)
			}

			if volumeMount.ReadOnly != container2.VolumeMounts[index].ReadOnly {
				log.Warnf("%s. %t vs %t", ErrorContainerVolumeMountsReadOnly.Error(), volumeMount.ReadOnly, container2.VolumeMounts[index].ReadOnly)
			}

			if volumeMount.MountPath != container2.VolumeMounts[index].MountPath {
				log.Warnf("%s. %s vs %s", ErrorContainerVolumeMountsMountPath.Error(), volumeMount.MountPath, container2.VolumeMounts[index].MountPath)
			}

			if volumeMount.SubPath != container2.VolumeMounts[index].SubPath {
				log.Warnf("%s. %s vs %s", ErrorContainerVolumeMountsSubPath.Error(), volumeMount.SubPath, container2.VolumeMounts[index].SubPath)
			}

			if volumeMount.SubPathExpr != container2.VolumeMounts[index].SubPathExpr {
				log.Warnf("%s. %s vs %s", ErrorContainerVolumeMountsSubPathExpr.Error(), volumeMount.SubPathExpr, container2.VolumeMounts[index].SubPathExpr)
			}

			if volumeMount.MountPropagation != nil && container2.VolumeMounts[index].MountPropagation != nil {
				if *volumeMount.MountPropagation != *container2.VolumeMounts[index].MountPropagation {
					log.Warnf("%s. %s vs %s", ErrorContainerVolumeMountsMountPropagation.Error(), *volumeMount.MountPropagation, *container2.VolumeMounts[index].MountPropagation)
				}

			} else if volumeMount.MountPropagation != nil || container2.VolumeMounts[index].MountPropagation != nil {

				log.Warnf("%s", ErrorMissingVolumeMountsMountPropagation.Error())
			}

		}
	}

	return diffs, nil
}
