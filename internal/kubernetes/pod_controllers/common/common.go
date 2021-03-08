package common

import (
	"context"

	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/kubernetes/common"
	"k8s-cluster-comparator/internal/kubernetes/metadata"

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

func ComparePodControllerSpecs(ctx context.Context, name string, apc1, apc2 *AbstractPodController) {
	var (
		log = logging.FromContext(ctx).With(zap.String("objectName", name))
	)
	ctx = logging.WithLogger(ctx, log)

	kind := types.ObjectKindWrapper(apc1.Metadata.Type.Kind)

	if metadata.IsMetadataDiffers(ctx, apc1.Metadata.Meta, apc2.Metadata.Meta) {
		return
	}

	log.Debugf("----- Start checking %s/%s pod controller spec -----", kind, apc1.Name)

	if apc1.Replicas != nil || apc2.Replicas != nil {
		if *apc1.Replicas != *apc2.Replicas {
			log.Warnf("the number of replicas is different: %d and %d", *apc1.Replicas, *apc2.Replicas)
		}
	}

	if (apc1.Replicas != nil && apc2.Replicas == nil) || (apc2.Replicas != nil && apc1.Replicas == nil) {
		log.Warnf("strange replicas specification difference: %#v and %#v", apc1.Replicas, apc2.Replicas)
	}

	// fill in the information that will be used for comparison
	object1 := types.InformationAboutObject{
		Template: apc1.PodTemplateSpec,
		Selector: apc1.PodLabelSelector,
	}
	object2 := types.InformationAboutObject{
		Template: apc2.PodTemplateSpec,
		Selector: apc2.PodLabelSelector,
	}

	matchLabelsString1 := common.ConvertMatchLabelsToString(ctx, object1.Selector.MatchLabels)
	matchLabelsString2 := common.ConvertMatchLabelsToString(ctx, object1.Selector.MatchLabels)

	if matchLabelsString1 != matchLabelsString2 {
		log.Warnf("%s: %s vs %s", ErrorMatchLabelsNotEqual.Error(), matchLabelsString1, matchLabelsString2)
	}

	bDiff, err := pods.ComparePodSpecs(ctx, object1, object2)
	if err != nil || bDiff {
		log.Warnw(err.Error())
	}

	log.Debugf("----- End checking %s/%s -----", kind, name)
}

// ComparePodControllers compares abstracted pod controller specifications in two k8s clusters
func ComparePodControllers(ctx context.Context, c1, c2 *ClusterCompareTask, namespace string) (bool, error) {
	var (
		log = logging.FromContext(ctx)

		flag bool
	)

	if len(c1.IsAlreadyCheckedFlagsMap) != len(c2.IsAlreadyCheckedFlagsMap) {
		log.Warnw("object counts are different", zap.Int("objectsCount1st", len(c1.IsAlreadyCheckedFlagsMap)), zap.Int("objectsCount2nd", len(c2.IsAlreadyCheckedFlagsMap)))
		flag = true
	}

	for name, index1 := range c1.IsAlreadyCheckedFlagsMap {
		ctx = logging.WithLogger(ctx, log.With(zap.String("objectName", name)))

		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
			if index2, ok := c2.IsAlreadyCheckedFlagsMap[name]; ok {
				index1.Check = true
				c1.IsAlreadyCheckedFlagsMap[name] = index1

				index2.Check = true
				c2.IsAlreadyCheckedFlagsMap[name] = index2

				apc1 := c1.APCList[index1.Index]
				apc2 := c2.APCList[index2.Index]

				// TODO: migrate to a goroutine
				ComparePodControllerSpecs(ctx, name, &apc1, &apc2)
			} else {
				log.With(zap.String("objectName", name)).Warn("object does not exist in 2nd cluster")
			}

		}
	}

	for name, index := range c2.IsAlreadyCheckedFlagsMap {
		if !index.Check {
			log.With(zap.String("objectName", name)).Warn("object does not exist in 1st cluster")
			flag = true
		}
	}

	return flag, nil
}
