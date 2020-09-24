package kubernetes

import (
	"reflect"
)

// CompareKVMap is a general function to compare two key-value maps
func CompareKVMap(map1, map2 map[string]string) bool {
	return reflect.DeepEqual(map1, map2)
}
