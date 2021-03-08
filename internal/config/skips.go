package config

import (
	"context"
	"fmt"

	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"
)

type ExcludeIncludeSpec struct {
	FullResourceNames map[types.ObjectKind][]types.ObjectName `yaml:"fullResourceNames"`
	NameBased         []types.ObjectName                      `yaml:"nameBased"`
}

// ParseSkipConfig parses information about entities to skip from a environment
func ParseSkipConfig(ctx context.Context, skipCfg ExcludeIncludeSpec) (skipper.SkipEntitiesList, error) {
	var (
		err error

		skipEntityList = skipper.SkipEntitiesList{}
	)

	fullResourceNamesSkips, err := skipper.ParseFullResourceNameSkipConfig(ctx, skipCfg.FullResourceNames)
	if err != nil {
		return skipEntityList, fmt.Errorf("invalid full resource name based skip config: %w", err)
	}

	skipEntityList.FullResourceNamesSkip = fullResourceNamesSkips

	nameBasedSkips, err := skipper.ParseNameBasedSkipConfig(ctx, skipCfg.NameBased)
	if err != nil {
		return skipEntityList, fmt.Errorf("invalid name based skip config: %w", err)
	}

	skipEntityList.NameBasedSkip = nameBasedSkips

	return skipEntityList, nil
}
