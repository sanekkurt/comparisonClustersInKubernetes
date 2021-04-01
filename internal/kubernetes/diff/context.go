package diff

import (
	"context"
)

type ctxKey string

var (
	diffsStoragePtrCtxKey ctxKey = "diffsStoragePtrCtxKey"
	diffsBatchPtrCtxKey   ctxKey = "diffsBatchPtrCtxKey"
)

func WithDiffStorage(ctx context.Context, diffs *DiffsStorage) context.Context {
	return context.WithValue(ctx, diffsStoragePtrCtxKey, diffs)
}

func DiffStorageFromContext(ctx context.Context) *DiffsStorage {
	diffsStorage, ok := ctx.Value(diffsStoragePtrCtxKey).(*DiffsStorage)
	if !ok {
		return nil
	}

	return diffsStorage
}

func WithDiffBatch(ctx context.Context, batch *DiffsBatch) context.Context {
	return context.WithValue(ctx, diffsBatchPtrCtxKey, batch)
}

func DiffBatchFromContext(ctx context.Context) *DiffsBatch {
	diffsBatch, ok := ctx.Value(diffsBatchPtrCtxKey).(*DiffsBatch)
	if !ok {
		return nil
	}

	return diffsBatch
}
