package config

import (
	"context"

	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/kubernetes/types"
)

var (
	kubeConnsCtxKey = struct{}{}
)

func With(ctx context.Context, kubeConns *types.KubeConnections) context.Context {
	return context.WithValue(ctx, kubeConnsCtxKey, kubeConns)
}

func FromContext(ctx context.Context) (kubernetes.Interface, kubernetes.Interface, string) {
	if kubeConns, ok := ctx.Value(kubeConnsCtxKey).(*types.KubeConnections); ok {
		return kubeConns.C1, kubeConns.C2, kubeConns.Namespace
	}
	return nil, nil, ""
}
