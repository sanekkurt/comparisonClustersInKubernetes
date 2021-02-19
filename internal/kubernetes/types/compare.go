package types

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KubeObjectsDifference struct {
	ObjectType metav1.TypeMeta
	ObjectMeta metav1.ObjectMeta

	Critical bool
}

type KubeResourceComparator interface {
	Compare(context.Context, string) ([]KubeObjectsDifference, error)
}
