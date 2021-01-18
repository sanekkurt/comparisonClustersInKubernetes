package pod_controllers

import "errors"

var (
	ErrorMatchLabelsNotEqual = errors.New("different matchLabels in pod controller specs")

	//ErrorContainersCountInPod      = errors.New("the containers count in pod are different")
	//ErrorContainerImageTemplatePod = errors.New("the container image in the template does not match the actual image in the Pod")
	ErrorDifferentImageIDInPods = errors.New("the ImageID in Pods is different")

//ErrorDifferentValueConfigMapKey = errors.New("the value for the ConfigMapKey is different")
//ErrorDifferentValueSecretKey    = errors.New("the value for the SecretKey is different")

//ErrorEnvironmentNotEqual = errors.New("the environment in containers not equal")
)
