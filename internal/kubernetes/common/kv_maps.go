package common

import (
	"bytes"
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
			log.With(zap.String("key", k)).Warn("key does not exist in map2")
			return false
		}

		if strings.Compare(val1, map2[k]) != 0 {
			log.With(zap.String("key", k), zap.String("value1", val1), zap.String("value2", map2[k])).Warn("key values do not match")
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
			log.With(zap.String("extraKeys", strings.Join(keys, ", "))).Warn("Extra keys found in 2nd map")

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
			log.With(zap.String("key", k)).Warn("key does not exist in map2")
			return false
		}

		if bytes.Compare(val1, map2[k]) != 0 {
			log.With(zap.String("key", k)).Warn("key values do not match")
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
			log.With(zap.String("extraKeys", strings.Join(keys, ", "))).Warn("Extra keys found in 2nd map")

			return false
		}
	}

	return true
}
