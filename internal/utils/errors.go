package utils

import "errors"

var (
	ErrorDifferentNumberValues = errors.New("different number of values in string lists")
	ErrorDifferentValues       = errors.New("different values in string lists")
	ErrorDuplicateElement      = errors.New("duplicate element in the list detected")
)
