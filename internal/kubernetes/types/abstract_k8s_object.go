package types

import (
	"context"
	"fmt"
	"reflect"

	"k8s-cluster-comparator/internal/logging"
)

// CompareAbstractObjectMetadata compares metadata of two abstract k8s objects
func CompareAbstractObjectMetadata(ctx context.Context, obj1, obj2 AbstractObjectMetadata) (bool, error) {
	log := logging.FromContext(ctx)

	if !reflect.DeepEqual(obj1.Type, obj2.Type) {
		return true, fmt.Errorf("object types are different: %s/%s and %s/%s, most likely this is an error", obj1.Type.APIVersion, obj1.Type.Kind, obj2.Type.APIVersion, obj2.Type.Kind)
	}

	if !reflect.DeepEqual(obj1.Meta.Labels, obj2.Meta.Labels) {
		log.Infof("object labels are different")
		return true, nil
	}

	if !reflect.DeepEqual(obj1.Meta.Annotations, obj2.Meta.Annotations) {
		log.Infof("object annotations are different")
		return true, nil
	}

	return false, nil
}
