package containers

import (
	"context"

	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/utils"
	v1 "k8s.io/api/core/v1"

	"k8s-cluster-comparator/internal/logging"
)

func compareContainerExecParams(ctx context.Context, container1, container2 v1.Container) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx)
		diffs = make([]types.KubeObjectsDifference, 0)
	)

	log.Debugf("ComparePodSpecs: start checking commands in container - %s", container1.Name)

	if bDiff, diff := utils.AreStringListsEqual(ctx, container1.Command, container2.Command); !bDiff {
		log.Warnf("%s. container '%s': %s", ErrorContainerCommandsDifferent, container1.Name, diff)
	}

	log.Debugf("ComparePodSpecs: started")
	if bDiff, diff := utils.AreStringListsEqual(ctx, container1.Args, container2.Args); !bDiff {
		log.Warnf("%s. container '%s': %s", ErrorContainerArgumentsDifferent, container1.Name, diff)
	}

	return diffs, nil
}
