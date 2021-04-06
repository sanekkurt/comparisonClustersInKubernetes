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
				return nil, fmt.Errorf("%w: %s", ErrorDuplicateElement, v)
			}

			log.Warnf("%s: %s", ErrorDuplicateElement, v)
		}

		out[v] = struct{}{}
	}

	return out, nil
}

// AreStringListsEqual compares two lists of strings
func AreStringListsEqual(ctx context.Context, l1, l2 []string) error {
	if len(l1) != len(l2) {
		return fmt.Errorf("%w: %d vs %d", ErrorDifferentNumberValues, len(l1), len(l2))
	}

	for index, value := range l1 {
		if value2 := l2[index]; value2 != value {
			return fmt.Errorf("%w at position #%d: %s vs %s", ErrorDifferentValues, index+1, value, l2[index])
		}
	}
	return nil
}
