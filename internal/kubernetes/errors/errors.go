package errors

import (
	"errors"
)

var (
	ErrorTooManyObjectsToCompare = errors.New("too many object given for comparison")

	//ErrorContainerNamesTemplate  = errors.New("container names in template are not equal")
	//ErrorContainerImagesTemplate = errors.New("container images in pod template are not equal")

	ErrorPodsCount = errors.New("the pods count are different")
)
