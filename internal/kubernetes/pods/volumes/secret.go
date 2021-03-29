package volumes

import (
	"context"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumeSecret(ctx context.Context, secret1, secret2 *v1.SecretVolumeSource) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx)

		//diffs = make([]types.KubeObjectsDifference, 0)
	)

	if secret1.SecretName != secret2.SecretName {
		log.Warnf("%s. %s vs %s", ErrorVolumeSecretName.Error(), secret1.SecretName, secret2.SecretName)
	}

	if secret1.Optional != nil && secret2.Optional != nil {

		if *secret1.Optional != *secret2.Optional {
			log.Warnf("%s. %t vs %t", ErrorVolumeSecretOptional.Error(), *secret1.Optional, *secret2.Optional)
		}

	} else if secret1.Optional != nil || secret2.Optional != nil {
		log.Warnf("%s", ErrorMissingVolumeSecretOptional.Error())
	}

	if secret1.DefaultMode != nil && secret2.DefaultMode != nil {

		if *secret1.DefaultMode != *secret2.DefaultMode {
			log.Warnf("%s. %d vs %d", ErrorVolumeSecretDefaultMode.Error(), *secret1.DefaultMode, *secret2.DefaultMode)
		}

	} else if secret1.DefaultMode != nil || secret2.DefaultMode != nil {
		log.Warnf("%s", ErrorMissingVolumeSecretDefaultMode.Error())
	}

	if len(secret1.Items) != len(secret2.Items) {
		log.Warnf("%s. %d vs %d", ErrorVolumeSecretItemsLen.Error(), len(secret1.Items), len(secret2.Items))
	} else {

		for index, item := range secret1.Items {

			if item.Path != secret2.Items[index].Path {
				log.Warnf("%s. %s vs %s", ErrorVolumeSecretItemsPath.Error(), item.Path, secret2.Items[index].Path)
			}

			if item.Key != secret2.Items[index].Key {
				log.Warnf("%s. %s vs %s", ErrorVolumeSecretItemsKey.Error(), item.Key, secret2.Items[index].Key)
			}

			if item.Mode != nil && secret2.Items[index].Mode != nil {
				if *item.Mode != *secret2.Items[index].Mode {
					log.Warnf("%s. %d vs %d", ErrorVolumeSecretItemsMode.Error(), item.Mode, secret2.Items[index].Mode)
				}
			} else if item.Mode != nil || secret2.Items[index].Mode != nil {
				log.Warnf("%s", ErrorMissingVolumeSecretItemsMode.Error())
			}

		}

	}

	return nil, nil
}
