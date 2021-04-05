package containers

import (
	"context"
	"errors"
	"strings"

	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/diff"

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

func compareContainerSpecImages(ctx context.Context, container1, container2 v1.Container) error {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)

		diffsBatch = diff.BatchFromContext(ctx)
	)

	var (
		imgParts = make([][]string, 2, 2)
	)

	for imgId, container := range []v1.Container{container1, container2} {
		imgParts[imgId] = strings.Split(container.Image, containerImageLabelTagSep)
	}

	for _, mirrorCfg := range cfg.Workloads.Containers.Image.Mirrors {
		if bMatched := strings.HasPrefix(imgParts[1][0], mirrorCfg.To); bMatched {
			log.Infof("using image mirror %s for image %s", mirrorCfg.To, imgParts[0][0])
			imgParts[1][0] = strings.Replace(imgParts[1][0], mirrorCfg.To, mirrorCfg.From, 1)
		}
	}

	if imgParts[0][0] != imgParts[1][0] {
		//log.Warnf("%s: %s vs %s", ErrorContainerDifferentImageLabels.Error(), imgParts[0][0], imgParts[1][0])
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s: %s vs %s", ErrorContainerDifferentImageLabels.Error(), imgParts[0][0], imgParts[1][0])
	}

	if imgParts[0][1] != imgParts[1][1] {
		//log.Warnf("%s: %s vs %s", ErrorContainerDifferentImageTags.Error(), imgParts[0][1], imgParts[1][1])
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s: %s vs %s", ErrorContainerDifferentImageTags.Error(), imgParts[0][1], imgParts[1][1])
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
		//log.Warnf("%s: %s vs %s", ErrorContainerDifferentImagePolicies.Error(), container1.ImagePullPolicy, container2.ImagePullPolicy)
		diffsBatch.Add(ctx, false, zap.WarnLevel, "%s: %s vs %s", ErrorContainerDifferentImagePolicies.Error(), container1.ImagePullPolicy, container2.ImagePullPolicy)
	}

	return nil
}
