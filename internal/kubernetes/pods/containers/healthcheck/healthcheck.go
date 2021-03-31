package healthcheck

import (
	"context"
	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	"reflect"
	"runtime"

	"k8s-cluster-comparator/internal/kubernetes/types"
	corev1 "k8s.io/api/core/v1"

	"k8s-cluster-comparator/internal/logging"
)

func Compare(ctx context.Context, container1, container2 corev1.Container) error {
	var (
		log = logging.FromContext(ctx)

		diffsBatch = ctx.Value("diffBatch").(*diff.DiffsBatch)
		meta       = ctx.Value("apcMeta").(types.AbstractObjectMetadata)
	)

	if container1.LivenessProbe != nil && container2.LivenessProbe != nil {
		log.Debugf("Compare liveness probes: started")

		err := compareContainerProbes(ctx, *container1.LivenessProbe, *container2.LivenessProbe)
		if err != nil {
			return err
		}

	} else if container1.LivenessProbe != nil || container2.LivenessProbe != nil {
		//logging.DiffLog(log, ErrorContainerHealthCheckLivenessProbeDifferent, "One of the containers has no liveness probe defined", *container1.LivenessProbe, *container2.LivenessProbe)
		diffsBatch.Add(ctx, &meta.Type, &meta.Meta, false, zap.WarnLevel, "%s", ErrorContainerHealthCheckLivenessProbeDifferent.Error())
	}

	if container1.ReadinessProbe != nil && container2.ReadinessProbe != nil {
		log.Debugf("Compare readiness probes: started")

		err := compareContainerProbes(ctx, *container1.ReadinessProbe, *container2.ReadinessProbe)
		if err != nil {
			return err
		}

	} else if container1.ReadinessProbe != nil || container2.ReadinessProbe != nil {
		//logging.DiffLog(log, ErrorContainerHealthCheckReadinessProbeDifferent, "One of the containers has no readiness probe defined", *container1.ReadinessProbe, *container2.ReadinessProbe)
		diffsBatch.Add(ctx, &meta.Type, &meta.Meta, false, zap.WarnLevel, "%s", ErrorContainerHealthCheckReadinessProbeDifferent.Error())
	}

	return nil
}

func compareContainerProbes(ctx context.Context, probe1, probe2 corev1.Probe) error {
	var (
		log = logging.FromContext(ctx)
	)

	comparisons := []func(context.Context, corev1.Probe, corev1.Probe){
		compareCommonProbeParams,
		compareTCPSocketProbes,
		compareHTTPGetProbes,
		compareExecProbes,
	}

	for _, cmp := range comparisons {
		log.Debugf("%s started", runtime.FuncForPC(reflect.ValueOf(cmp).Pointer()).Name())

		cmp(ctx, probe1, probe2)

	}

	return nil
}
