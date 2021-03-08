package containers

import (
	"context"
	"errors"
	"strings"

	v1 "k8s.io/api/core/v1"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/logging"
)

var (
	ErrorContainerDifferentImages = errors.New("different container images in Pod specs")

	ErrorContainerDifferentImageLabels   = errors.New("different container image labels in Pod specs")
	ErrorContainerDifferentImageTags     = errors.New("different container image tags in Pod specs")
	ErrorContainerDifferentImagePolicies = errors.New("different container image pullPolicies in Pod specs")
)

const (
	containerImageLabelTagSep = ":"
)

func compareContainerSpecImages(ctx context.Context, container1, container2 v1.Container) (bool, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)
	)

	var (
		imgParts = make([][]string, 2, 2)
	)

	for imgId, container := range []v1.Container{container1, container2} {
		imgParts[imgId] = strings.Split(container.Image, containerImageLabelTagSep)
	}

	if imgParts[0][0] != imgParts[1][0] {
		log.Warnf("%s: %s vs %s", ErrorContainerDifferentImageLabels.Error(), imgParts[0][0], imgParts[1][0])
		return true, ErrorContainerDifferentImageLabels
	}

	if imgParts[0][1] != imgParts[1][1] {
		log.Warnf("%s: %s vs %s", ErrorContainerDifferentImageTags.Error(), imgParts[0][1], imgParts[1][1])
		return true, ErrorContainerDifferentImageTags
	}

	if cfg.Workloads.Containers.RollingTags.WarnOnRollingTag {
	OUTER:
		for tag := range cfg.Workloads.Containers.RollingTags.TagsListMap {
			for idx := range imgParts {
				if imgParts[idx][1] == tag {
					log.Infof("cluster #%d: rolling tag '%s' detected, comparison might be inaccurate", idx+1, tag)
					break OUTER
				}
			}
		}
	}

	if container1.ImagePullPolicy != container2.ImagePullPolicy {
		log.Warnf("%s: %s vs %s", ErrorContainerDifferentImagePolicies.Error(), container1.ImagePullPolicy, container2.ImagePullPolicy)
		return true, ErrorContainerDifferentImagePolicies
	}

	return false, nil
}
