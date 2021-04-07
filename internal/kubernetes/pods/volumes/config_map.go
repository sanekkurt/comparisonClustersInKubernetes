package volumes

import (
	"context"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumeConfigMap(ctx context.Context, configMap1, configMap2 *v1.ConfigMapVolumeSource) {
	var (
		diffsBatch = diff.BatchFromContext(ctx)
	)

	if configMap1.Name != configMap2.Name {
		diffsBatch.Add(ctx, false, "%w: '%s' vs '%s'", ErrorVolumeConfigMapName, configMap1.Name, configMap2.Name)
	}

	if len(configMap1.Items) != len(configMap2.Items) {
		diffsBatch.Add(ctx, false, "%w", ErrorVolumeConfigMapItemsLen)
	} else {
		for i, item := range configMap1.Items {
			if item.Path != configMap2.Items[i].Path {
				diffsBatch.Add(ctx, false, "%w. '%s' vs '%s'", ErrorVolumeConfigMapPath, item.Path, configMap2.Items[i].Path)
			}

			if item.Key != configMap2.Items[i].Key {
				diffsBatch.Add(ctx, false, "%w. '%s' vs '%s'", ErrorVolumeConfigMapKey, item.Key, configMap2.Items[i].Key)
			}

			if item.Mode != nil && configMap2.Items[i].Mode != nil {
				if *item.Mode != *configMap2.Items[i].Mode {
					diffsBatch.Add(ctx, false, "%w. '%d' vs '%d'", ErrorVolumeConfigMapMode, *item.Mode, *configMap2.Items[i].Mode)
				}
			} else if item.Mode != nil || configMap2.Items[i].Mode != nil {
				diffsBatch.Add(ctx, false, "%w", ErrorVolumeConfigMapMode)
			}

		}
	}

	if configMap1.DefaultMode != nil && configMap2.DefaultMode != nil {
		if *configMap1.DefaultMode != *configMap2.DefaultMode {
			diffsBatch.Add(ctx, false, "%s: %d vs %d", ErrorVolumeConfigMapDefaultMode.Error(), *configMap1.DefaultMode, *configMap2.DefaultMode)
		}

	} else if configMap1.DefaultMode != nil || configMap2.DefaultMode != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorVolumeConfigMapDefaultMode.Error())
	}

	if configMap1.Optional != nil && configMap2.Optional != nil {
		if *configMap1.Optional != *configMap2.Optional {
			diffsBatch.Add(ctx, false, "%s: %t vs %t", ErrorVolumeConfigMapOptional.Error(), *configMap1.Optional, *configMap2.Optional)
		}

	} else if configMap1.DefaultMode != nil || configMap2.DefaultMode != nil {
		diffsBatch.Add(ctx, false, "%s", ErrorVolumeConfigMapOptional.Error())
	}

	return
}
