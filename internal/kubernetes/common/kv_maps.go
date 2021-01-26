package common

import (
	"context"
	"strings"

	"go.uber.org/zap"

	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

// CompareKVMap is a general function to compare two key-value maps
//func CompareKVMap(map1, map2 map[string]string) bool {
//	return reflect.DeepEqual(map1, map2)
//}

// AreKVMapsEqual is a general function to compare two key-value maps
func AreKVMapsEqual(ctx context.Context, map1, map2 types.KVMap, skipKeys map[string]struct{}) bool {
	log := logging.FromContext(ctx)

	for k, val1 := range map1 {
		if skipKeys != nil {
			if _, ok := skipKeys[k]; ok {
				log.Debugf("skip '%s' key from comparison due to skip rule", k)

				delete(map1, k)
				delete(map2, k)

				continue
			}
		}

		if _, ok := map2[k]; !ok {
			log.Warnf("key does not exist in map2", zap.String("key", k))
			return false
		}

		if val1 != map2[k] {
			log.Warnf("the value from map1 does not match the value from map2", zap.String("key", k), zap.String("value1", val1), zap.String("value2", map2[k]))
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
			log.Warnf("map2 has a number of extra keys which are not found in map1: %s", strings.Join(keys, ", "))

			return false
		}
	}

	return true
}

func AreKVBytesMapsEqual(ctx context.Context, map1, map2 map[string][]byte, skipKeys map[string]struct{}) bool {
	log := logging.FromContext(ctx)

	for k, val1 := range map1 {
		if skipKeys != nil {
			if _, ok := skipKeys[k]; ok {
				log.Debugf("skip '%s' key from comparison due to skip rule", k)

				delete(map1, k)
				delete(map2, k)

				continue
			}
		}

		if _, ok := map2[k]; !ok {
			log.Warnf("key does not exist in map2", zap.String("key", k))
			return false
		}

		if string(val1) != string(map2[k]) {
			log.Warnf("the data from map1 does not match the data from map2", zap.String("key", k))
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
			log.Warnf("map2 has a number of extra keys which are not found in map1: %s", strings.Join(keys, ", "))

			return false
		}
	}

	return true
}
