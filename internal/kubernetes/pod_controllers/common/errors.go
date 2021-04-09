package common

import (
	"errors"
)

var (
	ErrorMatchLabelsNotEqual = errors.New("different matchLabels in pod controller specs")

	//ErrorContainersCountInPod      = errors.New("the containers count in pod are different")
	//ErrorContainerImageTemplatePod = errors.New("the container image in the template does not match the actual image in the Pod")

	ErrorDifferentImageIDInPods = errors.New("the ImageID in Pods is different")

	ErrorDifferentNumberReplicas = errors.New("number of replicas is different")
	ErrorMissingReplicas         = errors.New("missing number replicas in one of the apc")
)
