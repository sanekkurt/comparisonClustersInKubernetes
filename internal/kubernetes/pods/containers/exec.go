package containers

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"

	"k8s-cluster-comparator/internal/logging"
)

func compareContainerExecParams(ctx context.Context, container1, container2 v1.Container) (bool, error) {
	var (
		log = logging.FromContext(ctx)
	)

	log.Debugf("ComparePodSpecs: start checking commands in container - %s", container1.Name)
	if err := CompareMassStringsInContainers(ctx, container1.Command, container2.Command); err != nil {
		return true, fmt.Errorf("%w. Name container: %s. %s", ErrorContainerCommandsDifferent, container1.Name, err)
	}

	log.Debugf("ComparePodSpecs: start checking args in container - %s", container1.Name)
	if err := CompareMassStringsInContainers(ctx, container1.Args, container2.Args); err != nil {
		return true, fmt.Errorf("%w. Name container: %s. %s", ErrorContainerArgumentsDifferent, container1.Name, err)
	}

	return false, nil
}
