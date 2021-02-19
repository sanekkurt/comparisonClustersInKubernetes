package metadata

import (
	"context"

	"go.uber.org/zap"
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

	if !CompareLabels(logging.WithLogger(ctx, log.With(zap.String("objectComponent", "labels"))), meta1.Labels, meta2.Labels, common.SkippedKubeLabels) {
		return false
	}

	if !CompareAnnotations(logging.WithLogger(ctx, log.With(zap.String("objectComponent", "annotations"))), meta1.Annotations, meta2.Annotations, common.SkippedKubeAnnotations) {
		return false
	}

	return true
}
