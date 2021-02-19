package skipper

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

// SkipComponentNames is a map of blank structs of k8s objects that must be skipped from the comparison
type SkipComponentNames map[types.ObjectName]struct{}

// SkipEntitiesList represents map[objectKind][listOfObjectNamesToSkipDuringCompare]
type SkipEntitiesList struct {
	FullResourceNamesSkip map[types.ObjectKind]SkipComponentNames
	NameBasedSkip         map[types.ObjectName]struct{}
}

func (skipCfg SkipEntitiesList) IsSkippedEntity(kind string, name string) bool {
	if _, bToSkip := skipCfg.NameBasedSkip[types.ObjectName(name)]; bToSkip {
		return true
	}

	if skipNames, ok := skipCfg.FullResourceNamesSkip[types.ObjectKind(kind)]; ok {
		if _, bToSkip := skipNames[types.ObjectName(name)]; bToSkip {
			return true
		}
	}

	return false
}

func ParseFullResourceNameSkipConfig(ctx context.Context, skipCfg map[types.ObjectKind][]types.ObjectName) (map[types.ObjectKind]SkipComponentNames, error) {
	var (
		log = logging.FromContext(ctx)

		objKinds = make(map[types.ObjectKind]struct{})

		fullResourceNames = make(map[types.ObjectKind]SkipComponentNames)
	)

	for k, names := range skipCfg {
		k := types.ObjectKind(strings.ToLower(string(k)))
		objNames := make(SkipComponentNames)

		if _, ok := objKinds[k]; ok {
			return nil, fmt.Errorf("kind '%s' specified multiple times", string(k))
		}

		for _, n := range names {
			n := types.ObjectName(strings.ToLower(string(n))) //nolint:govet

			if _, ok := objNames[n]; ok {
				return nil, fmt.Errorf("resource '%s/%s' specified multiple times", string(k), string(n))
			}

			objNames[n] = struct{}{}

			log.With(zap.String("kind", string(k)), zap.String("name", string(n))).Infof("resource '%s/%s' added to skip list", string(k), string(n))
		}

		fullResourceNames[k] = objNames
	}

	return fullResourceNames, nil
}

func ParseNameBasedSkipConfig(ctx context.Context, skipCfg []types.ObjectName) (SkipComponentNames, error) {
	var (
		log = logging.FromContext(ctx)

		list = make(SkipComponentNames)
	)

	for _, n := range skipCfg {
		n := types.ObjectName(strings.ToLower(string(n))) //nolint:govet

		if _, ok := list[n]; ok {
			return nil, fmt.Errorf("name '%s' specified multiple times", string(n))
		}
		list[n] = struct{}{}

		log.With(zap.String("name", string(n))).Infof("name '%s' added to skip list", string(n))
	}

	return list, nil
}
