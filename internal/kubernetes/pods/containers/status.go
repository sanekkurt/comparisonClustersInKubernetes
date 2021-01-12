package containers

import (
	"context"

	v12 "k8s.io/api/core/v1"

	"k8s-cluster-comparator/internal/kubernetes/types"
)

type PodContainerStatusesList map[int]types.Container

// GetPodContainerStatuses get statuses containers in Pod
func GetPodContainerStatuses(ctx context.Context, containerStatuses []v12.ContainerStatus) PodContainerStatusesList {
	var (
		container types.Container

		containerStatusesList = make(PodContainerStatusesList)
	)

	for index, value := range containerStatuses {
		container.Name = value.Name
		container.Image = value.Image
		container.ImageID = value.ImageID
		containerStatusesList[index] = container
	}

	return containerStatusesList
}
