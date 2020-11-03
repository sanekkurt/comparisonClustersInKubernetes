package pod_controllers

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"
)

func addItemsToDeploymentList(clientSet kubernetes.Interface, namespace string, limit int64) (*v1.DeploymentList, error) {
	log.Debugf("addItemsToDeploymentList started")
	defer log.Debugf("addItemsToDeploymentList completed")

	var (
		batch       *v1.DeploymentList
		deployments = &v1.DeploymentList{
			Items: make([]v1.Deployment, 0),
		}

		continueToken string

		err error
	)

	for {
		batch, err = clientSet.AppsV1().Deployments(namespace).List(metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		})
		if err != nil {
			return nil, err
		}

		log.Debugf("addItemsToDeploymentList: %d objects received", len(batch.Items))

		deployments.Items = append(deployments.Items, batch.Items...)

		deployments.TypeMeta = batch.TypeMeta
		deployments.ListMeta = batch.ListMeta

		if batch.Continue == "" {
			break
		}

		continueToken = batch.Continue
	}

	deployments.Continue = ""

	return deployments, err
}

func CompareDeployments(clientSet1, clientSet2 kubernetes.Interface, namespace string, skipEntityList skipper.SkipEntitiesList) (bool, error) {
	var (
		isClustersDiffer bool
	)

	deployments1, err := addItemsToDeploymentList(clientSet1, namespace, objectBatchLimit)
	if err != nil {
		return false, fmt.Errorf("cannot obtain deployments list from 1st cluster: %w", err)
	}

	deployments2, err := addItemsToDeploymentList(clientSet2, namespace, objectBatchLimit)
	if err != nil {
		return false, fmt.Errorf("cannot obtain deployments list from 2st cluster: %w", err)
	}

	apc1List, deploymentStatuses1, map1, apc2List, deploymentStatuses2, map2 := prepareDeploymentMaps(deployments1, deployments2, skipEntityList.GetByKind("deployments"))

	for _, value := range deploymentStatuses1 {
		if value.Replicas != value.ReadyReplicas {
			log.Info("Not all pods replicas are ready. The comparison may be incorrect")
		}
	}

	for _, value := range deploymentStatuses2 {
		if value.Replicas != value.ReadyReplicas {
			log.Info("Not all pods replicas are ready. The comparison may be incorrect")
		}
	}

	isClustersDiffer = comparePodControllerSpecs(&clusterCompareTask{
		Client:                   clientSet1,
		APCList:                  apc1List,
		IsAlreadyCheckedFlagsMap: map1,
	}, &clusterCompareTask{
		Client:                   clientSet2,
		APCList:                  apc2List,
		IsAlreadyCheckedFlagsMap: map2,
	}, namespace)

	return isClustersDiffer, nil
}

// prepareDeploymentMaps prepare deployment maps for comparison
func prepareDeploymentMaps(obj1, obj2 *v1.DeploymentList, skipEntities skipper.SkipComponentNames) ([]AbstractPodController, []v1.DeploymentStatus, map[string]types.IsAlreadyComparedFlag, []AbstractPodController, []v1.DeploymentStatus, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	var (
		map1                = make(map[string]types.IsAlreadyComparedFlag)
		apc1List            = make([]AbstractPodController, 0)
		deploymentStatuses1 = make([]v1.DeploymentStatus, 0)
		deploymentStatuses2 = make([]v1.DeploymentStatus, 0)

		map2     = make(map[string]types.IsAlreadyComparedFlag)
		apc2List = make([]AbstractPodController, 0)

		indexCheck types.IsAlreadyComparedFlag
	)

	for index, value := range obj1.Items {
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("deployment %s is skipped from comparison due to its name", value.Name)
			continue
		}

		indexCheck.Index = index
		map1[value.Name] = indexCheck

		apc1List = append(apc1List, AbstractPodController{
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
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("deployment %s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		map2[value.Name] = indexCheck

		apc2List = append(apc2List, AbstractPodController{
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
