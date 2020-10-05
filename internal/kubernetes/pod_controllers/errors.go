package pod_controllers

import "errors"

var (
	ErrorDiffersTemplatesNumber = errors.New("the number templates of containers differs") //nolint

	ErrorMatchlabelsNotEqual = errors.New("matchLabels are not equal")

	ErrorContainerNamesTemplate  = errors.New("container names in template are not equal")
	ErrorContainerImagesTemplate = errors.New("container name images in template are not equal")

	ErrorPodsCount = errors.New("the pods count are different")

	ErrorContainersCountInPod      = errors.New("the containers count in pod are different")
	ErrorContainerImageTemplatePod = errors.New("the container image in the template does not match the actual image in the Pod")

	ErrorDifferentImageInPods   = errors.New("the Image in Pods is different")
	ErrorDifferentImageIDInPods = errors.New("the ImageID in Pods is different")

	ErrorContainerNotFound = errors.New("container not found")
	ErrorNumberVariables   = errors.New("the number of variables in containers differs")

	ErrorDifferentValueConfigMapKey = errors.New("the value for the ConfigMapKey is different")
	ErrorDifferentValueSecretKey    = errors.New("the value for the SecretKey is different")

	ErrorEnvironmentNotEqual = errors.New("the environment in containers not equal")
)
