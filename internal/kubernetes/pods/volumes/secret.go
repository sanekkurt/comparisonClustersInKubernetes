package volumes

import (
	"context"
	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumeSecret(ctx context.Context, secret1, secret2 *v1.SecretVolumeSource) {
	var (
		diffsBatch = diff.DiffBatchFromContext(ctx)
	)

	if secret1.SecretName != secret2.SecretName {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %s vs %s", ErrorVolumeSecretName.Error(), secret1.SecretName, secret2.SecretName)
	}

	if secret1.Optional != nil && secret2.Optional != nil {

		if *secret1.Optional != *secret2.Optional {
			diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %t vs %t", ErrorVolumeSecretOptional.Error(), *secret1.Optional, *secret2.Optional)
		}

	} else if secret1.Optional != nil || secret2.Optional != nil {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s", ErrorMissingVolumeSecretOptional.Error())
	}

	if secret1.DefaultMode != nil && secret2.DefaultMode != nil {

		if *secret1.DefaultMode != *secret2.DefaultMode {
			diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %d vs %d", ErrorVolumeSecretDefaultMode.Error(), *secret1.DefaultMode, *secret2.DefaultMode)
		}

	} else if secret1.DefaultMode != nil || secret2.DefaultMode != nil {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s", ErrorMissingVolumeSecretDefaultMode.Error())
	}

	if len(secret1.Items) != len(secret2.Items) {
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %d vs %d", ErrorVolumeSecretItemsLen.Error(), len(secret1.Items), len(secret2.Items))
	} else {

		for index, item := range secret1.Items {

			if item.Path != secret2.Items[index].Path {
				diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %s vs %s", ErrorVolumeSecretItemsPath.Error(), item.Path, secret2.Items[index].Path)
			}

			if item.Key != secret2.Items[index].Key {
				diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %s vs %s", ErrorVolumeSecretItemsKey.Error(), item.Key, secret2.Items[index].Key)
			}

			if item.Mode != nil && secret2.Items[index].Mode != nil {
				if *item.Mode != *secret2.Items[index].Mode {
					diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. %d vs %d", ErrorVolumeSecretItemsMode.Error(), item.Mode, secret2.Items[index].Mode)
				}
			} else if item.Mode != nil || secret2.Items[index].Mode != nil {
				diffsBatch.Add(ctx, false, zap.WarnLevel, "%s", ErrorMissingVolumeSecretItemsMode.Error())
			}

		}

	}

	return
}
