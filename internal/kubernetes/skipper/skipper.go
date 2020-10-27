package skipper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
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

	if !strings.Contains(skipSpec, ";") {
		return nil, fmt.Errorf("does not contain valid data in the 'skip' variable. Between entities put ';' please")
	}

	var (
		tempSlice []string

		temp         = strings.Split(skipSpec, ";")
		skipEntities = make(map[types.ObjectKind]SkipComponentNames)

		tempMap = SkipComponentNames{}

		kind    string
		entries string
	)

	for _, entries = range temp {
		if !strings.Contains(entries, ":") {
			return nil, errors.New("does not contain valid data in the 'skip' variable. The enumeration of the names of entities start after ':' please or don't finish the line ';'")
		}

		tempSlice = strings.Split(entries, ":")

		kind = tempSlice[0]

		if strings.Contains(tempSlice[1], ",") {
			for _, entryName := range strings.Split(tempSlice[1], ",") {
				tempMap[types.ObjectName(entryName)] = struct{}{}
			}

			skipEntities[types.ObjectKind(tempSlice[0])] = make(map[types.ObjectName]struct{})

			for entryName := range tempMap {
				log.Infof("%s '%s' added to skip list", kind, entryName)

				skipEntities[types.ObjectKind(tempSlice[0])][entryName] = struct{}{}
				delete(tempMap, entryName)
			}
		} else {
			log.Infof("%s '%s' added to skip list", kind, tempSlice[1])

			skipEntities[types.ObjectKind(tempSlice[0])] = make(map[types.ObjectName]struct{})
			skipEntities[types.ObjectKind(tempSlice[0])][types.ObjectName(tempSlice[1])] = struct{}{}
		}
	}

	return skipEntities, nil
}
