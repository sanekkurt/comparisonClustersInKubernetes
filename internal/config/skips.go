package config

import (
	"context"
	"fmt"

	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"
)

type WorkMode struct{}

func (wm WorkMode) UnmarshalYAML(cb func(interface{}) error) error {
	var (
		mode string

		allowedModes = map[string]struct{}{
			"EverythingButNotExceptions": {},
			"NothingButGivenList":        {},
		}
	)

	err := cb(&mode)
	if err != nil {
		return err
	}

	if _, ok := allowedModes[mode]; !ok {
		return fmt.Errorf("unknown work mode '%s'", mode)
	}

	return nil
}

type SkipConfiguration struct {
	FullResourceNames map[types.ObjectKind][]types.ObjectName `yaml:"fullResourceNames"`
	NameBasedSkips    []types.ObjectName                      `yaml:"nameBasedSkips"`
}

// ParseSkipConfig parses information about entities to skip from a environment
func ParseSkipConfig(ctx context.Context, skipCfg SkipConfiguration) (skipper.SkipEntitiesList, error) {
	var (
		err error

		skipEntityList = skipper.SkipEntitiesList{}
	)

	fullResourceNamesSkips, err := skipper.ParseFullResourceNameSkipConfig(ctx, skipCfg.FullResourceNames)
	if err != nil {
		return skipEntityList, fmt.Errorf("invalid full resource name based skip config: %w", err)
	}

	skipEntityList.FullResourceNamesSkip = fullResourceNamesSkips

	nameBasedSkips, err := skipper.ParseNameBasedSkipConfig(ctx, skipCfg.NameBasedSkips)
	if err != nil {
		return skipEntityList, fmt.Errorf("invalid name based skip config: %w", err)
	}

	skipEntityList.NameBasedSkip = nameBasedSkips

	return skipEntityList, nil
}
