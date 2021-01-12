package containers

import (
	"errors"
)

var (
	ErrorContainerCommandsDifferent  = errors.New("container commands (entrypoints) in containers are different")
	ErrorContainerArgumentsDifferent = errors.New("container arguments (commands) in containers are different")

	ErrorContainerImageTagTemplatePod = errors.New("the container image tag in the template does not match the actual image tag in the Pod")

	ErrorDifferentImageInPods = errors.New("the Image in Pods is different")

	ErrorContainerNotFound = errors.New("container not found")
)
