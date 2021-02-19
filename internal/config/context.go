package config

import (
	"context"
)

var (
	configCtxKey struct{}
)

func With(ctx context.Context, cfg *AppConfig) context.Context {
	return context.WithValue(ctx, configCtxKey, cfg)
}

func FromContext(ctx context.Context) *AppConfig {
	cfg, ok := ctx.Value(configCtxKey).(*AppConfig)
	if !ok {
		return nil
	}

	return cfg
}
