package common

import (
	"context"
	"fmt"
	"k8s-cluster-comparator/internal/kubernetes/metadata"

	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"k8s.io/client-go/kubernetes"

	//"k8s-cluster-comparator/internal/kubernetes/metadata"
	"k8s-cluster-comparator/internal/kubernetes/pods"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

type ClusterCompareTask struct {
	Client                   kubernetes.Interface
	APCList                  []AbstractPodController
	IsAlreadyCheckedFlagsMap map[string]types.IsAlreadyComparedFlag
}

func ComparePodControllerSpecs(ctx context.Context, name string, apc1, apc2 *AbstractPodController) error {
	var (
		log = logging.FromContext(ctx).With(zap.String("objectName", name))

		kind = types.ObjectKindWrapper(apc1.Metadata.Type.Kind)

		diffsBatch = diff.StorageFromContext(ctx).NewLazyBatch(apc1.Metadata.Type, apc1.Metadata.Meta)
	)

	//ctx = diff.WithDiffBatch(ctx, diffsBatch) !!!!!!!!!!!!!!!!!!!!!!!!
	ctx = diff.WithDiffBatch(ctx, diffsBatch)
	ctx = logging.WithLogger(ctx, log)

	//??????????????????????????????????????????????????????????????????????????????????
	metadata.IsMetadataDiffers(ctx, apc1.Metadata.Meta, apc2.Metadata.Meta) // ?????????????????????????????????
	//??????????????????????????????????????????????????????????????????????????????????

	log.Debugf("%s/%s: check started", kind, apc1.Name)
	defer func() {
		log.Debugf("%s/%s: check completed", kind, apc1.Name)
	}()

	if apc1.Replicas != nil && apc2.Replicas != nil {
		if *apc1.Replicas != *apc2.Replicas {
			//log.Warnf("the number of replicas is different: %d and %d", *apc1.Replicas, *apc2.Replicas)
			//diffsBatch.Add(ctx, false, zap.WarnLevel, "the number of replicas is different: %d and %d", *apc1.Replicas, *apc2.Replicas)
			diffsBatch.Add(ctx, false, "%w: '%d' vs '%d'", ErrorDifferentNumberReplicas, *apc1.Replicas, *apc2.Replicas)
		}
	} else if apc1.Replicas != nil || apc2.Replicas != nil {
		diffsBatch.Add(ctx, false, "%w", ErrorMissingReplicas)
	}

	// fill in the information that will be used for comparison
	objects := []types.InformationAboutObject{{
		Template: apc1.PodTemplateSpec,
		Selector: apc1.PodLabelSelector,
	}, {
		Template: apc2.PodTemplateSpec,
		Selector: apc2.PodLabelSelector,
	}}

	err := pods.ComparePodSpecs(ctx, objects[0], objects[1])
	if err != nil {
		return fmt.Errorf("cannot compare Pod Specs: %w", err)
	}

	return nil
}

func CompareAbstractPodControllerMaps(ctx context.Context, kind string, apcs1, apcs2 map[string]*AbstractPodController) error {
	var (
		log = logging.FromContext(ctx)
	)

	if len(apcs1) != len(apcs2) {
		log.Warnw("object counts are different", zap.Int("objectsCount1st", len(apcs1)), zap.Int("objectsCount2nd", len(apcs2)))
	}

	for name, obj1 := range apcs1 {
		ctx := logging.WithLogger(ctx, log.With(zap.String("objectName", name)))

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if obj2, ok := apcs2[name]; ok {
				err := ComparePodControllerSpecs(ctx, name, obj1, obj2)
				if err != nil {
					return err
				}

				delete(apcs1, name)
				delete(apcs2, name)
			} else {
				log.With(zap.String("objectName", name)).Warnf("%s/%s does not exist in 2nd cluster", kind, name)
			}
		}
	}

	for name, _ := range apcs2 {
		log.With(zap.String("objectName", name)).Warnf("%s/%s does not exist in 1st cluster", kind, name)
	}

	return nil
}
