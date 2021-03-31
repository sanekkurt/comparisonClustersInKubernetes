package containers

import (
	"context"
	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/diff"

	"k8s-cluster-comparator/internal/kubernetes/pods/containers/env"
	"k8s-cluster-comparator/internal/kubernetes/pods/containers/healthcheck"
	"k8s-cluster-comparator/internal/kubernetes/types"

	v1 "k8s.io/api/core/v1"
)

func CompareContainerSpecs(ctx context.Context, container1, container2 v1.Container) error {
	var (
		//log = logging.FromContext(ctx)

		diffsBatch = ctx.Value("diffBatch").(*diff.DiffsBatch)
		meta       = ctx.Value("apcMeta").(types.AbstractObjectMetadata)
	)

	if container1.Name != container2.Name {
		diffsBatch.Add(ctx, &meta.Type, &meta.Meta, true, zap.WarnLevel, "%s: %s vs %s", ErrorContainerDifferentNames.Error(), container1.Name, container2.Name)
		//log.Warnf("%s: %s vs %s", ErrorContainerDifferentNames.Error(), container1.Name, container2.Name)
		return nil
	}

	err := compareContainerSpecImages(ctx, container1, container2)
	if err != nil {
		return err
	}

	err = env.Compare(ctx, container1.Env, container2.Env)
	if err != nil {
		return err
	}

	err = compareContainerExecParams(ctx, container1, container2)
	if err != nil {
		return err
	}

	err = healthcheck.Compare(ctx, container1, container2)
	if err != nil {
		return err
	}

	return nil
}
