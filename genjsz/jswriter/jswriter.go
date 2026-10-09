package jswriter

import (
	"context"

	"go.code.infinity6.ai/sdkgo/genjsz/jsparserapi"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparsersdkfunc"
)

type Options struct {
	Apis     []*jsparserapi.Apis
	SdkFuncs []*jsparsersdkfunc.SdkFuncs
}

func Write(ctx context.Context, opts Options) {
	panic("implement it")
}
