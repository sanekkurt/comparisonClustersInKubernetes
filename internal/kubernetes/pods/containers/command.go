package containers

import (
	"context"
	"errors"
	"fmt"
)

// CompareCommandsOrArgsInContainer compares commands or args in containers
func CompareMassStringsInContainers(ctx context.Context, mass1, mass2 []string) error {
	if len(mass1) != len(mass2) {
		return errors.New(fmt.Sprintf("different number of values in containers. count values in container 1 - %d value, count values in container 2 - %d value", len(mass1), len(mass2)))
	}

	for index, value := range mass1 {
		if value != mass2[index] {
			return errors.New(fmt.Sprintf("value in container 1  - %s, value in container 2 - %s", value, mass2[index]))
		}
	}
	return nil
}
