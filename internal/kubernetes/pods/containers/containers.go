package containers

import (
	"context"
	"errors"

	"k8s-cluster-comparator/internal/kubernetes/pods/containers/env"
	"k8s-cluster-comparator/internal/kubernetes/pods/containers/healthcheck"
	"k8s-cluster-comparator/internal/kubernetes/types"
	v1 "k8s.io/api/core/v1"

	"k8s-cluster-comparator/internal/logging"
)

var (
	ErrorContainerDifferentNames = errors.New("different container names in Pod specs")
)

func CompareContainerSpecs(ctx context.Context, container1, container2 v1.Container) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx)

		diffs = make([]types.KubeObjectsDifference, 0)
	)

	if container1.Name != container2.Name {
		log.Warnf("%s: %s vs %s", ErrorContainerDifferentNames.Error(), container1.Name, container2.Name)
	}

	bDiff, err := compareContainerSpecImages(ctx, container1, container2)
	if err != nil {
		return nil, err
	}
	diffs = append(diffs, bDiff...)

	bDiff, err = env.Compare(ctx, container1.Env, container2.Env)
	if err != nil {
		return nil, err
	}
	diffs = append(diffs, bDiff...)

	bDiff, err = compareContainerExecParams(ctx, container1, container2)
	if err != nil {
		return nil, err
	}
	diffs = append(diffs, bDiff...)

	bDiff, err = healthcheck.Compare(ctx, container1, container2)
	if err != nil {
		return nil, err
	}
	diffs = append(diffs, bDiff...)

	return diffs, nil
}
