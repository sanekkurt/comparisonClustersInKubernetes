package utils

import (
	"context"
	"errors"
	"fmt"
	"k8s-cluster-comparator/internal/logging"
	"testing"
)

func TestAreStringListsEqual(t *testing.T) {
	list1 := []string{"1", "2", "3", "4"}
	list2 := []string{"1", "2", "3", "4", "5"}

	err := AreStringListsEqual(context.Background(), list1, list2)
	if !errors.Is(err, ErrorDifferentNumberValues) {
		t.Errorf("Error expected: '%s. But it was returned: %s", ErrorDifferentNumberValues, err)
	}

	list1 = append(list1, "8")
	err = AreStringListsEqual(context.Background(), list1, list2)
	if !errors.Is(err, ErrorDifferentValues) {
		t.Errorf("Error expected: '%s. But it was returned: %s", ErrorDifferentValues, err)
	}
}

func TestStringsListToMap(t *testing.T) {

	var (
		ctx = context.Background()
	)

	err := logging.ConfigureForTests()
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
	}

	slice := []string{"1", "2", "3", "4"}

	resultMap := make(map[string]struct{})
	resultMap["1"] = struct{}{}
	resultMap["2"] = struct{}{}
	resultMap["3"] = struct{}{}
	resultMap["4"] = struct{}{}

	mapForChecks, err := StringsListToMap(ctx, slice, true)
	if len(resultMap) != len(mapForChecks) {

		t.Errorf("the resulting map does not match the required map. Different length. %d vs %d", len(resultMap), len(mapForChecks))

	} else {
		for key, _ := range resultMap {
			if _, ok := mapForChecks[key]; ok {
				delete(mapForChecks, key)
				delete(resultMap, key)
			} else {
				t.Errorf("no value found in the final map: %s", key)
			}
		}
	}

	slice = append(slice, "4")
	_, err = StringsListToMap(ctx, slice, true)
	if err != nil {
		if !errors.Is(err, ErrorDuplicateElement) {
			t.Errorf("Error expected: '%s. But it was returned: %s", ErrorDuplicateElement, err)
		}
	} else {
		t.Errorf("Error expected: '%s. But the error was not returned", ErrorDuplicateElement)
	}

}
