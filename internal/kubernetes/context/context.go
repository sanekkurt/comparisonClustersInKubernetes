package context

import (
	"context"

	"k8s.io/client-go/kubernetes"
)

type clientSetCtxKeyT string

var (
	clientSetCtxKey clientSetCtxKeyT = "clientSetCtxKey"
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
