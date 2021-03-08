package deployment

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/config"
	kubectx "k8s-cluster-comparator/internal/kubernetes/context"
	"k8s-cluster-comparator/internal/kubernetes/pod_controllers/common"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	deploymentKind = "deployment"
)

func deploymentsRetrieveBatchLimit(ctx context.Context) int64 {
	cfg := config.FromContext(ctx)

	if limit := cfg.Workloads.PodControllers.Deployments.BatchSize; limit != 0 {
		return limit
	}

	if limit := cfg.Common.DefaultBatchSize; limit != 0 {
		return limit
	}

	return 25
}

func fillInComparisonMap(ctx context.Context, namespace string, limit int64) (*v1.DeploymentList, error) {
	var (
		log       = logging.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		batch       *v1.DeploymentList
		deployments = &v1.DeploymentList{
			Items: make([]v1.Deployment, 0),
		}

		continueToken string

		err error
	)

	log.Debugf("fillInComparisonMap started")
	defer log.Debugf("fillInComparisonMap completed")

forLoop:
	for {
		select {
		case <-ctx.Done():
			return nil, context.Canceled
		default:
			batch, err = clientSet.AppsV1().Deployments(namespace).List(metav1.ListOptions{
				Limit:    limit,
				Continue: continueToken,
			})
			if err != nil {
				return nil, err
			}

			log.Debugf("fillInComparisonMap: %d objects received", len(batch.Items))

			deployments.Items = append(deployments.Items, batch.Items...)

			deployments.TypeMeta = batch.TypeMeta
			deployments.ListMeta = batch.ListMeta

			if batch.Continue == "" {
				break forLoop
			}

			continueToken = batch.Continue
		}
	}

	deployments.Continue = ""

	return deployments, err
}

type DeploymentsComparator struct {
}

func NewDeploymentsComparator(ctx context.Context, namespace string) DeploymentsComparator {
	return DeploymentsComparator{}
}

func (cmp DeploymentsComparator) Compare(ctx context.Context, namespace string) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", deploymentKind))
		cfg = config.FromContext(ctx)
	)

	ctx = logging.WithLogger(ctx, log)

	if !cfg.Workloads.Enabled ||
		!cfg.Workloads.PodControllers.Enabled ||
		!cfg.Workloads.PodControllers.Deployments.Enabled {
		log.Infof("'%s' kind skipped from comparison due to configuration", deploymentKind)
		return nil, nil
	}

	deployments1, err := fillInComparisonMap(kubectx.WithClientSet(ctx, cfg.Connections.Cluster1.ClientSet), namespace, deploymentsRetrieveBatchLimit(ctx))
	if err != nil {
		return nil, fmt.Errorf("cannot obtain deployments list from 1st cluster: %w", err)
	}

	deployments2, err := fillInComparisonMap(kubectx.WithClientSet(ctx, cfg.Connections.Cluster2.ClientSet), namespace, deploymentsRetrieveBatchLimit(ctx))
	if err != nil {
		return nil, fmt.Errorf("cannot obtain deployments list from 2st cluster: %w", err)
	}

	apc1List, deploymentStatuses1, map1, apc2List, deploymentStatuses2, map2 := prepareDeploymentMaps(ctx, deployments1, deployments2)

	for _, value := range deploymentStatuses1 {
		if value.Replicas != value.ReadyReplicas {
			log.Warn("Some deployment replicas in 1st cluster are not ready. The comparison might be inaccurate")
		}
	}

	for _, value := range deploymentStatuses2 {
		if value.Replicas != value.ReadyReplicas {
			log.Warn("Some deployment replicas in 2nd cluster are not ready. The comparison might be inaccurate")
		}
	}

	_, err = common.ComparePodControllers(ctx, &common.ClusterCompareTask{
		Client:                   cfg.Connections.Cluster1.ClientSet,
		APCList:                  apc1List,
		IsAlreadyCheckedFlagsMap: map1,
	}, &common.ClusterCompareTask{
		Client:                   cfg.Connections.Cluster2.ClientSet,
		APCList:                  apc2List,
		IsAlreadyCheckedFlagsMap: map2,
	}, namespace)

	return nil, err
}

// prepareDeploymentMaps prepare deployment maps for comparison
func prepareDeploymentMaps(ctx context.Context, obj1, obj2 *v1.DeploymentList) ([]common.AbstractPodController, []v1.DeploymentStatus, map[string]types.IsAlreadyComparedFlag, []common.AbstractPodController, []v1.DeploymentStatus, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)

		map1                = make(map[string]types.IsAlreadyComparedFlag)
		apc1List            = make([]common.AbstractPodController, 0)
		deploymentStatuses1 = make([]v1.DeploymentStatus, 0)
		deploymentStatuses2 = make([]v1.DeploymentStatus, 0)

		map2     = make(map[string]types.IsAlreadyComparedFlag)
		apc2List = make([]common.AbstractPodController, 0)

		indexCheck types.IsAlreadyComparedFlag
	)

	for index, value := range obj1.Items {
		if cfg.ExcludesIncludes.IsSkippedEntity(deploymentKind, value.Name) {
			log.With(zap.String("name", value.Name)).Debugf("deployment/%s is skipped from comparison", value.Name)
			continue
		}

		indexCheck.Index = index
		map1[value.Name] = indexCheck

		apc1List = append(apc1List, common.AbstractPodController{
			Metadata: types.AbstractObjectMetadata{
				Type: metav1.TypeMeta{
					Kind:       "deployments",
					APIVersion: "apps/v1",
				},
				Meta: value.ObjectMeta,
			},
			Name:             value.Name,
			Labels:           value.Labels,
			Annotations:      value.Annotations,
			Replicas:         value.Spec.Replicas,
			PodLabelSelector: value.Spec.Selector,
			PodTemplateSpec:  value.Spec.Template,
		})

		deploymentStatuses1 = append(deploymentStatuses1, value.Status)
	}

	for index, value := range obj2.Items {
		if cfg.ExcludesIncludes.IsSkippedEntity(deploymentKind, value.Name) {
			log.With(zap.String("name", value.Name)).Debugf("deployment/%s is skipped from comparison", value.Name)
			continue
		}

		indexCheck.Index = index
		map2[value.Name] = indexCheck

		apc2List = append(apc2List, common.AbstractPodController{
			Metadata: types.AbstractObjectMetadata{
				Type: metav1.TypeMeta{
					Kind:       "deployments",
					APIVersion: "apps/v1",
				},
				Meta: value.ObjectMeta,
			},
			Name:             value.Name,
			Labels:           value.Labels,
			Annotations:      value.Annotations,
			Replicas:         value.Spec.Replicas,
			PodLabelSelector: value.Spec.Selector,
			PodTemplateSpec:  value.Spec.Template,
		})
		deploymentStatuses2 = append(deploymentStatuses2, value.Status)
	}

	return apc1List, deploymentStatuses1, map1, apc2List, deploymentStatuses2, map2
}
