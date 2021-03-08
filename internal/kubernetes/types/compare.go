package types

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KubeObjectsDifference struct {
	ObjectType metav1.TypeMeta
	ObjectMeta metav1.ObjectMeta

	Error error

	Critical bool
}

type KubeResourceComparator interface {
	Compare(context.Context, string) ([]KubeObjectsDifference, error)
	//
	//fieldSelectorProvider (context.Context) string
	//labelSelectorProvider (context.Context) string
	//
	//collect(ctx context.Context) (interface{}, error)
}
