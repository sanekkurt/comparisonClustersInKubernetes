package containers

import (
	"context"

	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/diff"

	"k8s-cluster-comparator/internal/utils"
	v1 "k8s.io/api/core/v1"

	"k8s-cluster-comparator/internal/logging"
)

func compareContainerExecParams(ctx context.Context, container1, container2 v1.Container) error {
	var (
		log = logging.FromContext(ctx)

		diffsBatch = diff.BatchFromContext(ctx)
	)

	log.Debugf("ComparePodSpecs: start checking commands in container - %s", container1.Name)

	if bDiff, dif := utils.AreStringListsEqual(ctx, container1.Command, container2.Command); !bDiff {
		//log.Warnf("%s. container '%s': %s", ErrorContainerCommandsDifferent.Error(), container1.Name, dif)
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. container '%s': %s", ErrorContainerCommandsDifferent.Error(), container1.Name, dif)
	}

	log.Debugf("ComparePodSpecs: started")
	if bDiff, dif := utils.AreStringListsEqual(ctx, container1.Args, container2.Args); !bDiff {
		//log.Warnf("%s. container '%s': %s", ErrorContainerArgumentsDifferent.Error(), container1.Name, dif)
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s. container '%s': %s", ErrorContainerArgumentsDifferent.Error(), container1.Name, dif)
	}

	return nil
}
