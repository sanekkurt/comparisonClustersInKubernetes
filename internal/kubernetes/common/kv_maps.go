package common

import (
	"bytes"
	"context"
	"strings"

	"go.uber.org/zap"

	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

// AreKVMapsEqual is a general function to compare two key-value maps
func AreKVMapsEqual(ctx context.Context, map1, map2 types.KVMap, skipKeys map[string]struct{}, dumpValues bool) bool {
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

			delete(map1, k)
			delete(map2, k)

			continue
		}

		if strings.Compare(val1, map2[k]) != 0 {
			log := log.With(zap.String("key", k))

			if dumpValues {
				log = log.With(zap.String("value1", val1), zap.String("value2", map2[k]))
			}

			log.Warn("key values do not match")

			delete(map1, k)
			delete(map2, k)

			continue
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

func AreKVBytesMapsEqual(ctx context.Context, map1, map2 map[string][]byte, skipKeys map[string]struct{}, dumpValues bool) bool {
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

			delete(map1, k)
			delete(map2, k)

			continue
		}

		if bytes.Compare(val1, map2[k]) != 0 {
			log := log.With(zap.String("key", k))

			if dumpValues {
				log = log.With(zap.String("value1", string(val1)), zap.String("value2", string(map2[k])))
			}

			log.Warn("key values do not match")

			delete(map1, k)
			delete(map2, k)

			continue
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
		}
	}

	return true
}
