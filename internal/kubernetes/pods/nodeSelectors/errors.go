package nodeSelectors

import "errors"

var (
	ErrorDiffersNodeSelectorsInTemplates = errors.New("different NodeSelector in Pod templates") //nolint
)
