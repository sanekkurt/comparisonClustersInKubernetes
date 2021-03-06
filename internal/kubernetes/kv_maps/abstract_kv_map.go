package kv_maps

import (
	"strings"

	"k8s-cluster-comparator/internal/kubernetes/types"
)

// CompareKVMap is a general function to compare two key-value maps
//func CompareKVMap(map1, map2 map[string]string) bool {
//	return reflect.DeepEqual(map1, map2)
//}

// AreKVMapsEqual is a general function to compare two key-value maps
func AreKVMapsEqual(map1, map2 types.KVMap, skipKeys map[string]struct{}) bool {
	for k := range map1 {
		if skipKeys != nil {
			if _, ok := skipKeys[k]; ok {
				log.Debugf("skip '%s' key from comparison due to skip rule", k)

				delete(map1, k)
				delete(map2, k)

				continue
			}
		}

		if _, ok := map2[k]; !ok {
			log.Debugf("key '%s' does not exist in map2")
			return false
		}

		delete(map1, k)
		delete(map2, k)
	}

	if len(map2) > 0 {
		var keys = make([]string, 0, len(map2))

		for k := range map2 {
			if skipKeys != nil {
				if _, ok := skipKeys[k]; ok {
					log.Debugf("skip '%s' extra key from comparison due to skip rule", k)
					continue
				}
			}

			keys = append(keys, k)
		}

		if len(keys) > 0 {
			log.Debugf("the number of keys is not equal in maps. map2 contains following keys that does not exist in the map1: %s", strings.Join(keys, ","))

			return false
		}
	}

	return true
}
