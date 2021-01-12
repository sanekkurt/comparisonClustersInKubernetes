package pods

import (
	"errors"
)

var (
	ErrorDiffersTemplatesNumber = errors.New("different number of containers in Pod templates") //nolint
)
