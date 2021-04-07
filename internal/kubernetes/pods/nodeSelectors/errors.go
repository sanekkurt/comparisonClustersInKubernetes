package nodeSelectors

import "errors"

var (
	ErrorDiffersNodeSelectorsInTemplates = errors.New("different NodeSelector in Pod templates")                   //nolint
	ErrorNodeSelectorsDoesNotExist       = errors.New("nodeSelector does not exist in other cluster pod template") //nolint
)
