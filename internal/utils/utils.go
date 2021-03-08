package utils

import (
	"context"
	"fmt"

	"k8s-cluster-comparator/internal/logging"
)

func StringsListToMap(ctx context.Context, in []string, errorOnDuplicates bool) (map[string]struct{}, error) {
	var (
		log = logging.FromContext(ctx)

		out = make(map[string]struct{})
	)

	for _, v := range in {
		if _, ok := out[v]; ok {
			if errorOnDuplicates {
				return nil, fmt.Errorf("duplicate element '%s' in the list detected", v)
			}

			log.Warnf("duplicate element '%s' in the list detected", v)
		}

		out[v] = struct{}{}
	}

	return out, nil
}
