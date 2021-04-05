package diff

import (
	"context"
)

type ctxKey string

var (
	diffsStoragePtrCtxKey ctxKey = "diffsStoragePtrCtxKey"
	diffsBatchPtrCtxKey   ctxKey = "diffsBatchPtrCtxKey"
	//diffsChannelPtrCtxKey ctxKey = "diffsChannelPtrCtxKey"
)

func WithDiffStorage(ctx context.Context, diffs *DiffsStorage) context.Context {
	return context.WithValue(ctx, diffsStoragePtrCtxKey, diffs)
}

func StorageFromContext(ctx context.Context) *DiffsStorage {
	diffsStorage, ok := ctx.Value(diffsStoragePtrCtxKey).(*DiffsStorage)
	if !ok {
		return nil
	}

	return diffsStorage
}

func WithDiffBatch(ctx context.Context, batch *DiffsBatch) context.Context {
	return context.WithValue(ctx, diffsBatchPtrCtxKey, batch)
}

func BatchFromContext(ctx context.Context) *DiffsBatch {
	diffsBatch, ok := ctx.Value(diffsBatchPtrCtxKey).(*DiffsBatch)
	if !ok {
		return nil
	}

	return diffsBatch
}

//func WithDiffChannel(ctx context.Context, channel *ChanForDiff) context.Context {
//	return context.WithValue(ctx, diffsChannelPtrCtxKey, channel)
//}
//
//func ChanFromContext(ctx context.Context) *ChanForDiff {
//	channel, ok := ctx.Value(diffsChannelPtrCtxKey).(*ChanForDiff)
//	if !ok {
//		return nil
//	}
//
//	return channel
//}
