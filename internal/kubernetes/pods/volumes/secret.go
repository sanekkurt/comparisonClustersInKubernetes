package volumes

import (
	"context"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	v1 "k8s.io/api/core/v1"
)

func CompareVolumeSecret(ctx context.Context, secret1, secret2 *v1.SecretVolumeSource) {
	var (
		diffsBatch = diff.BatchFromContext(ctx)
	)

	if secret1.SecretName != secret2.SecretName {
		diffsBatch.Add(ctx, false, "%w. '%s' vs '%s'", ErrorVolumeSecretName, secret1.SecretName, secret2.SecretName)
	}

	if secret1.Optional != nil && secret2.Optional != nil {

		if *secret1.Optional != *secret2.Optional {
			diffsBatch.Add(ctx, false, "%w. '%t' vs '%t'", ErrorVolumeSecretOptional, *secret1.Optional, *secret2.Optional)
		}

	} else if secret1.Optional != nil || secret2.Optional != nil {
		diffsBatch.Add(ctx, false, "%w", ErrorMissingVolumeSecretOptional)
	}

	if secret1.DefaultMode != nil && secret2.DefaultMode != nil {

		if *secret1.DefaultMode != *secret2.DefaultMode {
			diffsBatch.Add(ctx, false, "%w. '%d' vs '%d'", ErrorVolumeSecretDefaultMode, *secret1.DefaultMode, *secret2.DefaultMode)
		}

	} else if secret1.DefaultMode != nil || secret2.DefaultMode != nil {
		diffsBatch.Add(ctx, false, "%w", ErrorMissingVolumeSecretDefaultMode)
	}

	if len(secret1.Items) != len(secret2.Items) {
		diffsBatch.Add(ctx, false, "%w. '%d' vs '%d'", ErrorVolumeSecretItemsLen, len(secret1.Items), len(secret2.Items))
	} else {

		for index, item := range secret1.Items {

			if item.Path != secret2.Items[index].Path {
				diffsBatch.Add(ctx, false, "%w. '%s' vs '%s'", ErrorVolumeSecretItemsPath, item.Path, secret2.Items[index].Path)
			}

			if item.Key != secret2.Items[index].Key {
				diffsBatch.Add(ctx, false, "%w. '%s' vs '%s'", ErrorVolumeSecretItemsKey, item.Key, secret2.Items[index].Key)
			}

			if item.Mode != nil && secret2.Items[index].Mode != nil {
				if *item.Mode != *secret2.Items[index].Mode {
					diffsBatch.Add(ctx, false, "%w. '%d' vs '%d'", ErrorVolumeSecretItemsMode, *item.Mode, *secret2.Items[index].Mode)
				}
			} else if item.Mode != nil || secret2.Items[index].Mode != nil {
				diffsBatch.Add(ctx, false, "%w", ErrorMissingVolumeSecretItemsMode)
			}

		}

	}

	return
}
