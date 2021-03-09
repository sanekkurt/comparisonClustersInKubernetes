package metadata

import (
	"context"

	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/common"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

func CompareLabels(ctx context.Context, map1, map2 types.KVMap, skipKeys map[string]struct{}) bool {
	var (
		cfg = config.FromContext(ctx)
	)

	return common.AreKVMapsEqual(ctx, map1, map2, skipKeys, cfg.Common.MetadataCompareConfiguration.DumpDifferentValues)
}

func CompareAnnotations(ctx context.Context, map1, map2 types.KVMap, skipKeys map[string]struct{}) bool {
	var (
		cfg = config.FromContext(ctx)
	)

	return common.AreKVMapsEqual(ctx, map1, map2, skipKeys, cfg.Common.MetadataCompareConfiguration.DumpDifferentValues)
}

func IsMetadataDiffers(ctx context.Context, meta1, meta2 v1.ObjectMeta) bool {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)
	)

	if !CompareLabels(logging.WithLogger(ctx, log.With(zap.String("objectComponent", "labels"))), meta1.Labels, meta2.Labels, cfg.Common.MetadataCompareConfiguration.SkipLabelsMap) {
		return false
	}

	if !CompareAnnotations(logging.WithLogger(ctx, log.With(zap.String("objectComponent", "annotations"))), meta1.Annotations, meta2.Annotations, cfg.Common.MetadataCompareConfiguration.SkipAnnotationsMap) {
		return false
	}

	return true
}
