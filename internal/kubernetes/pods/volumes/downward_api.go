package volumes

import (
	"context"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumeDownwardAPI(ctx context.Context, downwardAPI1, downwardAPI2 *v1.DownwardAPIVolumeSource) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx)

		//diffs = make([]types.KubeObjectsDifference, 0)
	)

	if downwardAPI1.DefaultMode != nil && downwardAPI2.DefaultMode != nil {

		if *downwardAPI1.DefaultMode != *downwardAPI2.DefaultMode {
			log.Warnf("%s. %d vs %d", ErrorDownwardAPIDefaultMode.Error(), *downwardAPI1.DefaultMode, *downwardAPI2.DefaultMode)
		}

	} else if downwardAPI1.DefaultMode != nil || downwardAPI2.DefaultMode != nil {
		log.Warnf("%s", ErrorMissingDownwardAPIDefaultMode.Error())
	}

	if len(downwardAPI1.Items) != len(downwardAPI2.Items) {
		log.Warnf("%s. %d vs %d", ErrorVolumeDownwardAPIItemsLen.Error(), len(downwardAPI1.Items), len(downwardAPI2.Items))
	} else {

		for index, item := range downwardAPI1.Items {

			if item.Path != downwardAPI2.Items[index].Path {
				log.Warnf("%s. %s vs %s", ErrorDownwardAPIItemsPath.Error(), item.Path, downwardAPI2.Items[index].Path)
			}

			if item.Mode != nil && downwardAPI2.Items[index].Mode != nil {

				if *item.Mode != *downwardAPI2.Items[index].Mode {
					log.Warnf("%s. %d vs %d", ErrorDownwardAPIItemsMode.Error(), *item.Mode, *downwardAPI2.Items[index].Mode)
				}

			} else if item.Mode != nil || downwardAPI2.Items[index].Mode != nil {
				log.Warnf("%s", ErrorMissingDownwardAPIItemsMode.Error())
			}

			if item.ResourceFieldRef != nil && downwardAPI2.Items[index].ResourceFieldRef != nil {

				if item.ResourceFieldRef.Resource != downwardAPI2.Items[index].ResourceFieldRef.Resource {
					log.Warnf("%s. %s vs %s", ErrorDownwardAPIItemsResFieldRefResource.Error(), item.ResourceFieldRef.Resource, downwardAPI2.Items[index].ResourceFieldRef.Resource)
				}

				if item.ResourceFieldRef.ContainerName != downwardAPI2.Items[index].ResourceFieldRef.ContainerName {
					log.Warnf("%s. %s vs %s", ErrorDownwardAPIItemsResFieldRefContainerName.Error(), item.ResourceFieldRef.ContainerName, downwardAPI2.Items[index].ResourceFieldRef.ContainerName)
				}

				if item.ResourceFieldRef.Divisor.Format != downwardAPI2.Items[index].ResourceFieldRef.Divisor.Format {
					log.Warnf("%s. %s vs %s", ErrorDownwardAPIItemsResFieldRefFormat.Error(), item.ResourceFieldRef.Divisor.Format, downwardAPI2.Items[index].ResourceFieldRef.Divisor.Format)
				}

			} else if item.ResourceFieldRef != nil || downwardAPI2.Items[index].ResourceFieldRef != nil {
				log.Warnf("%s", ErrorMissingDownwardAPIItemsResourceFieldRef.Error())
			}

			if item.FieldRef != nil && downwardAPI2.Items[index].FieldRef != nil {

				if item.FieldRef.APIVersion != downwardAPI2.Items[index].FieldRef.APIVersion {
					log.Warnf("%s. %s vs %s", ErrorDownwardAPIItemsFieldRefAPIVersion.Error(), item.FieldRef.APIVersion, downwardAPI2.Items[index].FieldRef.APIVersion)
				}

				if item.FieldRef.FieldPath != downwardAPI2.Items[index].FieldRef.FieldPath {
					log.Warnf("%s. %s vs %s", ErrorDownwardAPIItemsFieldRefFieldPath.Error(), item.FieldRef.FieldPath, downwardAPI2.Items[index].FieldRef.FieldPath)
				}

			} else if item.FieldRef != nil || downwardAPI2.Items[index].FieldRef != nil {
				log.Warnf("%s", ErrorMissingDownwardAPIItemsFieldRef.Error())
			}

		}
	}

	return nil, nil
}
