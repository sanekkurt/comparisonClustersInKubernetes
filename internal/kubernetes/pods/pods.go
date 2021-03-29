package pods

import (
	"context"
	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/pods/nodeSelectors"
	"k8s-cluster-comparator/internal/kubernetes/pods/volumes"

	"k8s-cluster-comparator/internal/kubernetes/pods/containers"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

//func CompareRunningPods() {
//	if !simplifiedVerification {
//
//		if len(pods1.Items) != len(pods2.Items) {
//			return ErrorPodsCount
//		}
//		for controlledPod1Idx := range pods1.Items {
//			var (
//				flag                       int
//				containerWithSameNameFound bool
//				templateHasAbsolutePath    bool
//
//				containersStatusesInPod1               = GetContainerStatusesInPod(ctx, pods1.Items[controlledPod1Idx].Status.ContainerStatuses)
//				containersStatusesInPod2               = GetContainerStatusesInPod(ctx, pods2.Items[controlledPod1Idx].Status.ContainerStatuses)
//				containersDeploymentTemplateSplitLabel = strings.Split(containersPod1[podTemplate1ContainerIdx].Image, ":")
//			)
//
//			if strings.Contains(containersPod1[podTemplate1ContainerIdx].Image, "/") {
//				log.Debugf("ComparePodSpecs: image in template - %s has absolute path. image: %s", containersPod1[podTemplate1ContainerIdx].Name, containersPod1[podTemplate1ContainerIdx].Image)
//				templateHasAbsolutePath = true
//			} else {
//				log.Debugf("ComparePodSpecs: image in template - %s doesn't have an absolute path. image: %s", containersPod1[podTemplate1ContainerIdx].Name, containersPod1[podTemplate1ContainerIdx].Image)
//			}
//
//			if len(containersStatusesInPod1) != len(containersStatusesInPod2) {
//				log.Debug("ComparePodSpecs: ErrorContainersCountInPod. Pod 1 name - %s, pod 2 name - %s. Count in pod 1 - %d, count in pod 2 - %d", pods1.Items[controlledPod1Idx].Name, pods2.Items[controlledPod1Idx].Name, len(containersStatusesInPod1), len(containersStatusesInPod2))
//				return ErrorContainersCountInPod
//			}
//
//			for controlledPod1ContainerStatusIdx := range containersStatusesInPod1 {
//				if containersPod1[podTemplate1ContainerIdx].Name == containersStatusesInPod1[controlledPod1ContainerStatusIdx].Name && containersPod1[podTemplate1ContainerIdx].Name == containersStatusesInPod2[controlledPod1ContainerStatusIdx].Name { //nolint:gocritic,unused
//					flag++
//
//					var containersStatusesInPod1SplitLabel []string
//					var containersStatusesInPod2SplitLabel []string
//
//					if templateHasAbsolutePath {
//
//						containersStatusesInPod1SplitLabel = strings.Split(containersStatusesInPod1[controlledPod1ContainerStatusIdx].Image, ":")
//						containersStatusesInPod2SplitLabel = strings.Split(containersStatusesInPod2[controlledPod1ContainerStatusIdx].Image, ":")
//
//					} else {
//
//						if strings.Contains(containersStatusesInPod1[controlledPod1ContainerStatusIdx].Image, "/") || strings.Contains(containersStatusesInPod2[controlledPod1ContainerStatusIdx].Image, "/") {
//
//							pathImage1 := strings.Split(containersStatusesInPod1[controlledPod1ContainerStatusIdx].Image, "/")
//							pathImage2 := strings.Split(containersStatusesInPod2[controlledPod1ContainerStatusIdx].Image, "/")
//							containersStatusesInPod1SplitLabel = strings.Split(pathImage1[len(pathImage1)-1], ":")
//							containersStatusesInPod2SplitLabel = strings.Split(pathImage2[len(pathImage2)-1], ":")
//							log.Debugf("ComparePodSpecs: image in pod - %s has `/` so it was divided. containersStatusesInPod1SplitLabel - %s, containersStatusesInPod2SplitLabel - %s", containersStatusesInPod1[controlledPod1ContainerStatusIdx].Name, fmt.Sprintln(containersStatusesInPod1SplitLabel), fmt.Sprintln(containersStatusesInPod2SplitLabel))
//
//						} else {
//
//							containersStatusesInPod1SplitLabel = strings.Split(containersStatusesInPod1[controlledPod1ContainerStatusIdx].Image, ":")
//							containersStatusesInPod2SplitLabel = strings.Split(containersStatusesInPod2[controlledPod1ContainerStatusIdx].Image, ":")
//							log.Debugf("ComparePodSpecs: image in pod - %s doesn't have `/` it was therefore divided as follows. containersStatusesInPod1SplitLabel - %s, containersStatusesInPod2SplitLabel - %s", containersStatusesInPod1[controlledPod1ContainerStatusIdx].Name, fmt.Sprintln(containersStatusesInPod1SplitLabel), fmt.Sprintln(containersStatusesInPod2SplitLabel))
//
//						}
//					}
//
//					if containersDeploymentTemplateSplitLabel[0] != containersStatusesInPod1SplitLabel[0] || containersDeploymentTemplateSplitLabel[0] != containersStatusesInPod2SplitLabel[0] { //nolint:gocritic,unused
//						log.Debugf("ComparePodSpecs: Image not equal in containersDeploymentTemplate and in containersStatusesInPod. containersDeploymentTemplateSplitLabel - %s, containersStatusesInPod1SplitLabel - %s, containersStatusesInPod2SplitLabel - %s", containersDeploymentTemplateSplitLabel[0], containersStatusesInPod1SplitLabel[0], containersStatusesInPod2SplitLabel[0])
//						return ErrorContainerImageTemplatePod
//					}
//
//					if len(containersDeploymentTemplateSplitLabel) > 1 {
//						if containersDeploymentTemplateSplitLabel[1] != containersStatusesInPod1SplitLabel[1] || containersDeploymentTemplateSplitLabel[1] != containersStatusesInPod2SplitLabel[1] {
//
//							log.Infof("Container name - %s. the container image tag in the template does not match the actual image tag in the pod: template image tag - %s, pod1 image tag - %s, pod2 image tag - %s", containersPod1[podTemplate1ContainerIdx].Name, containersDeploymentTemplateSplitLabel[1], containersStatusesInPod1SplitLabel[1], containersStatusesInPod2SplitLabel[1])
//
//							if switchFatalDifferentTag {
//								log.Debug("ComparePodSpecs: Container name - %s. the container image tag in the template does not match the actual image tag in the pod: template image tag - %s, pod1 image tag - %s, pod2 image tag - %s", containersPod1[podTemplate1ContainerIdx].Name, containersDeploymentTemplateSplitLabel[1], containersStatusesInPod1SplitLabel[1], containersStatusesInPod2SplitLabel[1])
//								return ErrorContainerImageTagTemplatePod
//
//							}
//						}
//					}
//
//					for _, value := range containersStatusesInPod2 {
//						if containersStatusesInPod1[controlledPod1ContainerStatusIdx].Name == value.Name {
//
//							containerWithSameNameFound = true
//
//							if containersStatusesInPod1[controlledPod1ContainerStatusIdx].Image != value.Image {
//								return fmt.Errorf("%w. \nPods name: '%s'. Image name on pod1: '%s'. Image name on pod2: '%s'", ErrorDifferentImageInPods, value.Name, containersStatusesInPod1[controlledPod1Idx].Image, value.Image)
//							}
//							if containersStatusesInPod1[controlledPod1ContainerStatusIdx].ImageID != value.ImageID {
//								return fmt.Errorf("%w. Pods name: '%s'. ImageID on pod1: '%s'. ImageID on pod2: '%s'", ErrorDifferentImageIDInPods, value.Name, containersStatusesInPod1[controlledPod1Idx].ImageID, value.ImageID)
//							}
//						}
//					}
//					if !containerWithSameNameFound {
//						return fmt.Errorf("%w. Name container: %s", ErrorContainerNotFound, containersStatusesInPod1[controlledPod1Idx].Name)
//					}
//				}
//			}
//		}
//	}
//}

// ComparePodSpecs compares pod templates of two abstract pod controllers
func ComparePodSpecs(ctx context.Context, spec1, spec2 types.InformationAboutObject) ([]types.KubeObjectsDifference, error) {
	var (
		log   = logging.FromContext(ctx)
		diffs = make([]types.KubeObjectsDifference, 0)
	)

	log.Debugf("ComparePodSpecs (pod/%s, pod/%s): started", spec1.Template.Name, spec2.Template.Name)
	defer func() {
		log.Debug("ComparePodSpecs: completed")
	}()

	var (
		containersPod1   = spec1.Template.Spec.Containers
		containersPod2   = spec2.Template.Spec.Containers
		nodeSelectorPod1 = spec1.Template.Spec.NodeSelector
		nodeSelectorPod2 = spec2.Template.Spec.NodeSelector
		volumesPod1      = spec1.Template.Spec.Volumes
		volumesPod2      = spec2.Template.Spec.Volumes
	)

	if len(containersPod1) != len(containersPod2) {
		log.Warnf("%s: %d vs %d", ErrorDiffersContainersNumberInTemplates.Error(), len(containersPod1), len(containersPod2))
		return nil, nil
	}

	//pods1, err := common.GetPodsListOnMatchLabels(ctx, spec1.Selector.MatchLabels, namespace, clientSet1)
	//if err != nil {
	//	return false, err
	//}
	//pods2, err := common.GetPodsListOnMatchLabels(ctx, spec1.Selector.MatchLabels, namespace, clientSet2)
	//if err != nil {
	//	return false, err
	//}

	for podTemplate1ContainerIdx := range containersPod1 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			log := log.With(zap.String("containerName", containersPod1[podTemplate1ContainerIdx].Name))
			ctx := logging.WithLogger(ctx, log)

			diff, err := containers.CompareContainerSpecs(ctx, containersPod1[podTemplate1ContainerIdx], containersPod2[podTemplate1ContainerIdx])
			if err != nil {
				return nil, err
			}
			diffs = append(diffs, diff...)
		}
	}

	if nodeSelectorPod1 != nil && nodeSelectorPod2 != nil {

		if len(nodeSelectorPod1) != len(nodeSelectorPod2) {
			log.Warnf("%s", ErrorDiffersNodeSelectorsNumberInTemplates.Error())
			return nil, nil
		}

	} else if nodeSelectorPod1 != nil || nodeSelectorPod2 != nil {

		log.Warnf("%s", ErrorPodMissingNodeSelectors.Error())
		return nil, nil

	} else {
		nodeSelectors.CompareNodeSelectors(ctx, nodeSelectorPod1, nodeSelectorPod2)
	}

	if volumesPod1 != nil && volumesPod2 != nil {
		if len(volumesPod1) != len(volumesPod2) {
			log.Warnf("%s", ErrorDiffersVolumesNumberInTemplates.Error())
			return nil, nil
		}
	} else if volumesPod1 != nil && volumesPod2 != nil {
		log.Warnf("%s", ErrorPodMissingVolumes.Error())
		return nil, nil
	}

	for podTemplate1VolumeIdx := range volumesPod1 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			log := log.With(zap.String("volumeName", volumesPod1[podTemplate1VolumeIdx].Name))
			ctx := logging.WithLogger(ctx, log)

			diff, err := volumes.CompareVolumes(ctx, volumesPod1[podTemplate1VolumeIdx], volumesPod2[podTemplate1VolumeIdx])
			if err != nil {
				return nil, err
			}
			diffs = append(diffs, diff...)
		}
	}

	return diffs, nil
}
