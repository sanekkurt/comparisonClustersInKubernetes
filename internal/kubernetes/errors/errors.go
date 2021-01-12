package errors

import (
	"errors"
)

var (
	ErrorMatchlabelsNotEqual = errors.New("different matchLabels")

	//ErrorContainerNamesTemplate  = errors.New("container names in template are not equal")
	//ErrorContainerImagesTemplate = errors.New("container images in pod template are not equal")

	ErrorPodsCount = errors.New("the pods count are different")
)
