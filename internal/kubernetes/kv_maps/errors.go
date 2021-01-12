package kv_maps

import (
	"errors"
)

var (
	ErrorKVMapNoSuchKey = errors.New("required key does not exist in map")
)
