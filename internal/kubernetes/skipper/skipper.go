package skipper

import (
	"errors"
	"strings"

	"k8s-cluster-comparator/internal/kubernetes/types"
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
func ParseSkipConfig(skipSpec string) (SkipEntitiesList, error) {
	if !strings.Contains(skipSpec, ";") {
		return nil, errors.New("does not contain valid data in the 'skip' variable. Between entities put ';' please")
	}

	var (
		tempSlice []string

		temp         = strings.Split(skipSpec, ";")
		skipEntities = make(map[types.ObjectKind]SkipComponentNames)

		tempMap = SkipComponentNames{}
	)

	for _, value := range temp {
		if !strings.Contains(value, ":") {
			return nil, errors.New("does not contain valid data in the 'skip' variable. The enumeration of the names of entities start after ':' please or don't finish the line ';'")
		}

		tempSlice = strings.Split(value, ":")

		if strings.Contains(tempSlice[1], ",") {
			for _, val := range strings.Split(tempSlice[1], ",") {
				tempMap[types.ObjectName(val)] = struct{}{}
			}

			skipEntities[types.ObjectKind(tempSlice[0])] = make(map[types.ObjectName]struct{})

			for key, value := range tempMap {
				skipEntities[types.ObjectKind(tempSlice[0])][key] = value
				delete(tempMap, key)
			}
		} else {
			skipEntities[types.ObjectKind(tempSlice[0])] = make(map[types.ObjectName]struct{})
			skipEntities[types.ObjectKind(tempSlice[0])][types.ObjectName(tempSlice[0])] = struct{}{}
		}
	}

	return skipEntities, nil
}
