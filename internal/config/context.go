package config

import (
	"context"
)

type configCtxKeyT string

var (
	configCtxKey configCtxKeyT = "configCtxKey"
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
