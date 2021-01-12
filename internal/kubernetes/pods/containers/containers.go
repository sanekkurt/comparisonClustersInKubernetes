package containers

import (
	"context"
	"errors"

	v1 "k8s.io/api/core/v1"

	"k8s-cluster-comparator/internal/logging"
)

var (
	ErrorContainerDifferentNames  = errors.New("different container names in Pod specs")
	ErrorContainerDifferentImages = errors.New("different container images in Pod specs")
)

func CompareContainerSpecs(ctx context.Context, container1, container2 v1.Container) (bool, error) {
	var (
		log = logging.FromContext(ctx)
	)

	if container1.Name != container2.Name {
		log.Warnf("%s: %s vs %s", ErrorContainerDifferentNames.Error(), container1.Name, container2.Name)
		return true, ErrorContainerDifferentNames
	}

	if container1.Image != container2.Image {
		log.Warnf("%s: %s vs %s", ErrorContainerDifferentImages.Error(), container1.Image, container2.Image)
		return true, ErrorContainerDifferentImages
	}

	bDiff, err := compareContainerEnvVars(ctx, container1.Env, container2.Env)
	if err != nil {
		return false, err
	}
	if bDiff {
		return true, err
	}

	bDiff, err = compareContainerExecParams(ctx, container1, container2)
	if err != nil {
		return false, err
	}
	if bDiff {
		return true, err
	}

	bDiff, err = compareContainerHealthCheckParams(ctx, container1, container2)
	if err != nil {
		return false, err
	}
	if bDiff {
		return true, err
	}

	return false, nil
}
