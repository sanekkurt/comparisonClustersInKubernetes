package common

import (
	"context"

	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/pods"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"

	batchv1 "k8s.io/api/batch/v1"
)

func CompareJobSpecInternals(ctx context.Context, obj1, obj2 batchv1.JobSpec) (bool, error) {
	log := logging.FromContext(ctx)

	if obj1.BackoffLimit != nil && obj2.BackoffLimit != nil {
		if *obj1.BackoffLimit != *obj2.BackoffLimit {
			log.Warnw("Job backoff limit is different", zap.Int32("backoffLimit1", *obj1.BackoffLimit), zap.Int32("backoffLimit2", *obj2.BackoffLimit))
			return true, ErrorBackoffLimitDifferent
		}
	} else if obj1.BackoffLimit != nil || obj2.BackoffLimit != nil {
		return true, ErrorBackoffLimitDifferent
	}

	if obj1.Template.Spec.RestartPolicy != obj2.Template.Spec.RestartPolicy {
		log.Warnw("Job restartPolicy limit is different", zap.String("restartPolicy1", string(obj1.Template.Spec.RestartPolicy)), zap.String("restartPolicy2", string(obj1.Template.Spec.RestartPolicy)))
		return true, ErrorRestartPolicyDifferent
	}

	castJob1ForCompareContainers := types.InformationAboutObject{
		Template: obj1.Template,
		Selector: nil,
	}
	castJob2ForCompareContainers := types.InformationAboutObject{
		Template: obj2.Template,
		Selector: nil,
	}

	bDiff, err := pods.ComparePodSpecs(ctx, castJob1ForCompareContainers, castJob2ForCompareContainers)
	if err != nil || bDiff {
		return bDiff, err
	}

	return false, nil
}
