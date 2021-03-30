package diff

import (
	"context"
)

type ctxKey string

var (
	diffsStoragePtrCtxKey ctxKey = "diffsStoragePtrCtxKey"
)

func With(ctx context.Context, diffs *DiffsStorage) context.Context {
	return context.WithValue(ctx, diffsStoragePtrCtxKey, diffs)
}

func FromContext(ctx context.Context) *DiffsStorage {
	diffs, ok := ctx.Value(diffsStoragePtrCtxKey).(*DiffsStorage)
	if !ok {
		return nil
	}

	return diffs
}
