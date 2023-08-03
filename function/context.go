package function

import (
	"context"
	"time"
)

type FunctionCtx struct {
	ctx context.Context
}

type keyType string

const (
	blockNumberKey keyType = "block_number"
	blockHashKey   keyType = "block_hash"
)

func NewFunctionCtx(parentCtx context.Context, blockNumber int64, blockHash string) *FunctionCtx {
	ctx := context.WithValue(parentCtx, blockNumberKey, blockNumber)
	ctx = context.WithValue(ctx, blockHashKey, blockHash)
	return &FunctionCtx{ctx: ctx}
}

func (f *FunctionCtx) BlockNumber() int64 {
	return f.ctx.Value(blockNumberKey).(int64)
}

func (f *FunctionCtx) BlockHash() string {
	return f.ctx.Value(blockHashKey).(string)
}

// Implement the context.Context interface

func (f *FunctionCtx) Deadline() (deadline time.Time, ok bool) {
	return f.ctx.Deadline()
}

func (f *FunctionCtx) Done() <-chan struct{} {
	return f.ctx.Done()
}

func (f *FunctionCtx) Err() error {
	return f.ctx.Err()
}

func (f *FunctionCtx) Value(key interface{}) interface{} {
	return f.ctx.Value(key)
}
