package metadata

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s-cluster-comparator/internal/kubernetes/common"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

func CompareLabels(ctx context.Context, map1, map2 types.KVMap, skipKeys map[string]struct{}) bool {
	return common.AreKVMapsEqual(ctx, map1, map2, skipKeys)
}

func CompareAnnotations(ctx context.Context, map1, map2 types.KVMap, skipKeys map[string]struct{}) bool {
	return common.AreKVMapsEqual(ctx, map1, map2, skipKeys)
}

func IsMetadataDiffers(ctx context.Context, meta1, meta2 v1.ObjectMeta) bool {
	log := logging.FromContext(ctx)

	if !CompareLabels(ctx, meta1.Labels, meta2.Labels, common.SkippedKubeLabels) {
		log.Warnw("metadata differs: different labels")
		return false
	}

	if !CompareAnnotations(ctx, meta1.Annotations, meta2.Annotations, common.SkippedKubeAnnotations) {
		log.Warnw("metadata differs: different annotations")
		return false
	}

	return true
}
