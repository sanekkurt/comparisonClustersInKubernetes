package context

import (
	"context"

	"k8s.io/client-go/kubernetes"
)

type ctxKey string

var (
	clientSetCtxKey ctxKey = "clientSetCtxKey"

	namespaceCtxKey ctxKey = "namespaceCtxKey"
)

func WithClientSet(ctx context.Context, cset kubernetes.Interface) context.Context {
	return context.WithValue(ctx, clientSetCtxKey, cset)
}

func ClientSetFromContext(ctx context.Context) kubernetes.Interface {
	cfg, ok := ctx.Value(clientSetCtxKey).(kubernetes.Interface)
	if !ok {
		return nil
	}

	return cfg
}

func WithNamespace(ctx context.Context, namespace string) context.Context {
	return context.WithValue(ctx, namespaceCtxKey, namespace)
}

func NamespaceFromContext(ctx context.Context) string {
	namespace, ok := ctx.Value(namespaceCtxKey).(string)
	if !ok {
		return ""
	}

	return namespace
}
