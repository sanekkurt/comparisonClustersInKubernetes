package skipper

import (
	"context"
	"fmt"
	"strings"

	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
)

const (
	skipEntityTypesSep     = ";"
	skipEntityTypeNamesSep = ":"
	skipEntityNamesSep     = ","
)

// SkipComponentNames is a map of blank structs of k8s objects that must be skipped from the comparison
type SkipComponentNames map[types.ObjectName]struct{}

func (skipNames SkipComponentNames) IsSkippedEntity(name string) bool {
	_, bToSkip := skipNames[types.ObjectName(name)]
	return bToSkip
}

// SkipEntitiesList represents map[objectKind][listOfObjectNamesToSkipDuringCompare]
type SkipEntitiesList map[types.ObjectKind]SkipComponentNames

func (sel SkipEntitiesList) GetByKind(kind string) SkipComponentNames {
	l, ok := sel[types.ObjectKind(types.ObjectKindWrapper(kind))]
	if !ok {
		return nil
	}

	return l
}

// ParseSkipConfig parses information about entities to skip from a environment
func ParseSkipConfig(ctx context.Context, skipSpec string) (SkipEntitiesList, error) {
	log := logging.FromContext(ctx)

	var (
		skipEntitiesList = make(map[types.ObjectKind]SkipComponentNames)

		entityKind  string
		entityNames string
	)

	skipSpec = strings.ReplaceAll(skipSpec, " ", "")

	for _, skipEntityTypeNamesPair := range strings.Split(skipSpec, skipEntityTypesSep) {
		if !strings.Contains(skipEntityTypeNamesPair, skipEntityTypeNamesSep) {
			return nil, fmt.Errorf("does not contain valid data in the SKIP variable. The enumeration of the names of entities starts after '%s'", skipEntityTypeNamesSep)
		}

		skipEntityTypeNames := strings.Split(skipEntityTypeNamesPair, skipEntityTypeNamesSep)

		entityKind = skipEntityTypeNames[0]
		entityNames = skipEntityTypeNames[1]

		if skipEntitiesList[types.ObjectKind(entityKind)] == nil {
			skipEntitiesList[types.ObjectKind(entityKind)] = make(map[types.ObjectName]struct{})
		}

		for _, entityName := range strings.Split(entityNames, skipEntityNamesSep) {
			log.Infof("%s/%s added to skip list", entityKind, entityName)
			skipEntitiesList[types.ObjectKind(entityKind)][types.ObjectName(entityName)] = struct{}{}
		}
	}

	return skipEntitiesList, nil
}
