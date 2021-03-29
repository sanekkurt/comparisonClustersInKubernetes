package pods

import (
	"errors"
)

var (
	ErrorDiffersContainersNumberInTemplates    = errors.New("different number of containers in Pod templates")    //nolint
	ErrorDiffersNodeSelectorsNumberInTemplates = errors.New("different number of NodeSelectors in Pod templates") //nolint
	ErrorPodMissingNodeSelectors               = errors.New("one of the pods is missing NodeSelectors")           //nolint
	ErrorDiffersVolumesNumberInTemplates       = errors.New("different number of volumes in Pod templates")       //nolint
	ErrorPodMissingVolumes                     = errors.New("one of the pods is missing volumes")                 //nolint
)
