package containers

import (
	"errors"
)

var (
	ErrorContainerDifferentNames = errors.New("different container names in Pod specs")

	ErrorContainerCommandsDifferent  = errors.New("container commands (entrypoints) in containers are different")
	ErrorContainerArgumentsDifferent = errors.New("container arguments (commands) in containers are different")

	ErrorContainerVolumeMountsLength           = errors.New("container VolumeMounts in containers have different length")
	ErrorContainerVolumeMountsName             = errors.New("different container VolumeMounts N in Pod specs")
	ErrorContainerVolumeMountsReadOnly         = errors.New("different container VolumeMounts ReadOnly in Pod specs")
	ErrorContainerVolumeMountsMountPath        = errors.New("different container VolumeMounts MountPath in Pod specs")
	ErrorContainerVolumeMountsSubPath          = errors.New("different container VolumeMounts SubPath in Pod specs")
	ErrorContainerVolumeMountsSubPathExpr      = errors.New("different container VolumeMounts SubPathExpr in Pod specs")
	ErrorContainerVolumeMountsMountPropagation = errors.New("different container VolumeMounts MountPropagation in Pod specs")
	ErrorMissingVolumeMountsMountPropagation   = errors.New("missing container VolumeMounts MountPropagation in Pod specs")

	ErrorContainerImageTagTemplatePod = errors.New("the container image tag in the template does not match the actual image tag in the Pod")

	ErrorDifferentImageInPods = errors.New("the Image in Pods is different")

	ErrorContainerNotFound = errors.New("container not found")
)
