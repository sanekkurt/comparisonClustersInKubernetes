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

// AreStringListsEqual compares two lists of strings
func AreStringListsEqual(ctx context.Context, l1, l2 []string) (bool, string) {
	if len(l1) != len(l2) {
		return false, fmt.Sprintf("different number of values in string lists: %d vs %d", len(l1), len(l2))
	}

	for index, value := range l1 {
		if value2 := l2[index]; value2 != value {
			return false, fmt.Sprintf("different values in string lists at position #%d: %s vs %s", index+1, value, l2[index])
		}
	}
	return true, ""
}
