package volumes

import (
	"context"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumeConfigMap(ctx context.Context, configMap1, configMap2 *v1.ConfigMapVolumeSource) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx)

		//diffs = make([]types.KubeObjectsDifference, 0)
	)

	if configMap1.Name != configMap2.Name {
		log.Warnf("%s: %s vs %s", ErrorVolumeConfigMapName.Error(), configMap1.Name, configMap2.Name)
	}

	if len(configMap1.Items) != len(configMap2.Items) {
		log.Warnf("%s", ErrorVolumeConfigMapItemsLen.Error())
	} else {
		for i, item := range configMap1.Items {
			if item.Path != configMap2.Items[i].Path {
				log.Warnf("%s. %s vs %s", ErrorVolumeConfigMapPath.Error(), item.Path, configMap2.Items[i].Path)
			}

			if item.Key != configMap2.Items[i].Key {
				log.Warnf("%s. %s vs %s", ErrorVolumeConfigMapKey.Error(), item.Key, configMap2.Items[i].Key)
			}

			if item.Mode != nil && configMap2.Items[i].Mode != nil {
				if *item.Mode != *configMap2.Items[i].Mode {
					log.Warnf("%s. %d vs %d", ErrorVolumeConfigMapMode.Error(), *item.Mode, *configMap2.Items[i].Mode)
				}
			} else if item.Mode != nil || configMap2.Items[i].Mode != nil {
				log.Warnf("%s", ErrorVolumeConfigMapMode.Error())
			}

		}
	}

	if configMap1.DefaultMode != nil && configMap2.DefaultMode != nil {
		if *configMap1.DefaultMode != *configMap2.DefaultMode {
			log.Warnf("%s: %d vs %d", ErrorVolumeConfigMapDefaultMode.Error(), *configMap1.DefaultMode, *configMap2.DefaultMode)
		}

	} else if configMap1.DefaultMode != nil || configMap2.DefaultMode != nil {
		log.Warnf("%s", ErrorVolumeConfigMapDefaultMode.Error())
	}

	if configMap1.Optional != nil && configMap2.Optional != nil {
		if *configMap1.Optional != *configMap2.Optional {
			log.Warnf("%s: %t vs %t", ErrorVolumeConfigMapOptional.Error(), *configMap1.Optional, *configMap2.Optional)
		}

	} else if configMap1.DefaultMode != nil || configMap2.DefaultMode != nil {
		log.Warnf("%s", ErrorVolumeConfigMapOptional.Error())
	}

	return nil, nil
}
