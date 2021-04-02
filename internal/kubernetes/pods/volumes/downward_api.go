package volumes

import (
	"context"
	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumeDownwardAPI(ctx context.Context, downwardAPI1, downwardAPI2 *v1.DownwardAPIVolumeSource) {
	var (
		diffsBatch = diff.BatchFromContext(ctx)
	)

	if downwardAPI1.DefaultMode != nil && downwardAPI2.DefaultMode != nil {

		if *downwardAPI1.DefaultMode != *downwardAPI2.DefaultMode {
			diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %d vs %d", ErrorDownwardAPIDefaultMode.Error(), *downwardAPI1.DefaultMode, *downwardAPI2.DefaultMode)
		}

	} else if downwardAPI1.DefaultMode != nil || downwardAPI2.DefaultMode != nil {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s", ErrorMissingDownwardAPIDefaultMode.Error())
	}

	if len(downwardAPI1.Items) != len(downwardAPI2.Items) {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %d vs %d", ErrorVolumeDownwardAPIItemsLen.Error(), len(downwardAPI1.Items), len(downwardAPI2.Items))
	} else {

		for index, item := range downwardAPI1.Items {

			if item.Path != downwardAPI2.Items[index].Path {
				diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %s vs %s", ErrorDownwardAPIItemsPath.Error(), item.Path, downwardAPI2.Items[index].Path)
			}

			if item.Mode != nil && downwardAPI2.Items[index].Mode != nil {

				if *item.Mode != *downwardAPI2.Items[index].Mode {
					diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %d vs %d", ErrorDownwardAPIItemsMode.Error(), *item.Mode, *downwardAPI2.Items[index].Mode)
				}

			} else if item.Mode != nil || downwardAPI2.Items[index].Mode != nil {
				diffsBatch.Add(ctx, false, zap.WarnLevel, "%s", ErrorMissingDownwardAPIItemsMode.Error())
			}

			if item.ResourceFieldRef != nil && downwardAPI2.Items[index].ResourceFieldRef != nil {

				if item.ResourceFieldRef.Resource != downwardAPI2.Items[index].ResourceFieldRef.Resource {
					diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %s vs %s", ErrorDownwardAPIItemsResFieldRefResource.Error(), item.ResourceFieldRef.Resource, downwardAPI2.Items[index].ResourceFieldRef.Resource)
				}

				if item.ResourceFieldRef.ContainerName != downwardAPI2.Items[index].ResourceFieldRef.ContainerName {
					diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %s vs %s", ErrorDownwardAPIItemsResFieldRefContainerName.Error(), item.ResourceFieldRef.ContainerName, downwardAPI2.Items[index].ResourceFieldRef.ContainerName)
				}

				if item.ResourceFieldRef.Divisor.Format != downwardAPI2.Items[index].ResourceFieldRef.Divisor.Format {
					diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %s vs %s", ErrorDownwardAPIItemsResFieldRefFormat.Error(), item.ResourceFieldRef.Divisor.Format, downwardAPI2.Items[index].ResourceFieldRef.Divisor.Format)
				}

			} else if item.ResourceFieldRef != nil || downwardAPI2.Items[index].ResourceFieldRef != nil {
				diffsBatch.Add(ctx, false, zap.WarnLevel, "%s", ErrorMissingDownwardAPIItemsResourceFieldRef.Error())
			}

			if item.FieldRef != nil && downwardAPI2.Items[index].FieldRef != nil {

				if item.FieldRef.APIVersion != downwardAPI2.Items[index].FieldRef.APIVersion {
					diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %s vs %s", ErrorDownwardAPIItemsFieldRefAPIVersion.Error(), item.FieldRef.APIVersion, downwardAPI2.Items[index].FieldRef.APIVersion)
				}

				if item.FieldRef.FieldPath != downwardAPI2.Items[index].FieldRef.FieldPath {
					diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %s vs %s", ErrorDownwardAPIItemsFieldRefFieldPath.Error(), item.FieldRef.FieldPath, downwardAPI2.Items[index].FieldRef.FieldPath)
				}

			} else if item.FieldRef != nil || downwardAPI2.Items[index].FieldRef != nil {
				diffsBatch.Add(ctx, false, zap.WarnLevel, "%s", ErrorMissingDownwardAPIItemsFieldRef.Error())
			}

		}
	}

	return
}
