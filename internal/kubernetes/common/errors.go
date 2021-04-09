package common

import "errors"

var (
	ErrorKeyDoesNotExistInMap2   = errors.New("key does not exist in map2")
	ErrorExtraKeysFoundInMap2    = errors.New("extra keys found in 2nd map")
	ErrorKeyValueDoNotMatchInMap = errors.New("key values do not match")
)
