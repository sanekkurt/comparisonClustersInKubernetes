package pod_controllers

import "errors"

var (
	ErrorDiffersTemplatesNumber = errors.New("the number templates of containers differs") //nolint

	ErrorMatchlabelsNotEqual = errors.New("matchLabels are not equal")

	ErrorContainerNamesTemplate  = errors.New("container names in template are not equal")
	ErrorContainerImagesTemplate = errors.New("container name images in template are not equal")
	ErrorContainerCommandsDifferent = errors.New("сommands in containers are different")
	ErrorContainerArgumentsDifferent = errors.New("arguments in containers are different")

	ErrorContainerLivenessProbeDifferent = errors.New("livenessProbe in containers are different")
	ErrorContainerReadinessProbeDifferent = errors.New("readinessProbe in containers are different")


	ErrorPodsCount = errors.New("the pods count are different")

	ErrorContainersCountInPod      = errors.New("the containers count in pod are different")
	ErrorContainerImageTemplatePod = errors.New("the container image in the template does not match the actual image in the Pod")
	ErrorContainerImageTagTemplatePod = errors.New("the container image tag in the template does not match the actual image tag in the Pod")

	ErrorDifferentImageInPods   = errors.New("the Image in Pods is different")
	ErrorDifferentImageIDInPods = errors.New("the ImageID in Pods is different")

	ErrorContainerNotFound = errors.New("container not found")
	ErrorNumberVariables   = errors.New("the number of variables in containers differs")

	ErrorDifferentValueConfigMapKey = errors.New("the value for the ConfigMapKey is different")
	ErrorDifferentValueSecretKey    = errors.New("the value for the SecretKey is different")

	ErrorEnvironmentNotEqual = errors.New("the environment in containers not equal")

	ErrorDifferentExec = errors.New("the exec command missing in one probe")
	ErrorDifferentExecCommand = errors.New("the exec command in probe not equal")
	ErrorDifferentTCPSocket = errors.New("the TCPSocket missing in one probe")
	ErrorDifferentTCPSocketHost = errors.New("the TCPSocket.Host in probe not equal")
	ErrorDifferentTCPSocketPort = errors.New("the TCPSocket.Port in probe not equal")
	ErrorDifferentHTTPGet = errors.New("the HTTPGet missing in one probe")
	ErrorDifferentHTTPGetHost = errors.New("the HTTPGet.Host in probe not equal")
	ErrorDifferentHTTPGetHTTPHeaders = errors.New("the HTTPGet.HTTPHeaders in probe not equal")
	ErrorDifferentNameHeader = errors.New("the name header in probe not equal")
	ErrorDifferentValueHeader = errors.New("the value header in probe not equal")
	ErrorMissingHeader = errors.New("one of the containers is missing headers")
	ErrorDifferentHTTPGetPath = errors.New("the HTTPGet.Path in probe not equal")
	ErrorDifferentHTTPGetPort = errors.New("the HTTPGet.Port in probe not equal")
	ErrorDifferentHTTPGetScheme = errors.New("the HTTPGet.Scheme in probe not equal")
	ErrorDifferentFailureThreshold = errors.New("the FailureThreshold in probe not equal")
	ErrorDifferentInitialDelaySeconds = errors.New("the InitialDelaySeconds in probe not equal")
	ErrorDifferentPeriodSeconds = errors.New("the PeriodSeconds in probe not equal")
	ErrorDifferentSuccessThreshold = errors.New("the SuccessThreshold in probe not equal")
	ErrorDifferentTimeoutSeconds = errors.New("the TimeoutSeconds in probe not equal")
)
