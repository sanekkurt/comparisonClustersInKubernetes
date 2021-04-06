package types

import (
	"context"

	"k8s-cluster-comparator/internal/kubernetes/diff"
)

type KubeResourceComparator interface {
	Compare(context.Context) (*diff.DiffsStorage, error)

	FieldSelectorProvider(context.Context) string
	LabelSelectorProvider(context.Context) string

	//collect(ctx context.Context) (interface{}, error)
}
