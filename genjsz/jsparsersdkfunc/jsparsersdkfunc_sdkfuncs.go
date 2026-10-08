package jsparsersdkfunc

import (
	"context"

	"go.code.infinity6.ai/platform/logz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

type SdkFuncs struct {
	SdkFuncs []*SdkFunc       `json:"sdk_funcs"`
	Typedefs jstypez.Typedefs `json:"typedefs"`
}

func (s *SdkFuncs) Funcs() []*SdkFunc {
	return s.SdkFuncs
}

type Funcs = SdkFuncs

var ParseFuncs = ParseSdkFuncs

func ParseSdkFuncs(ctx context.Context, funcs []sdkfuncz.Func) *SdkFuncs {
	ret := SdkFuncs{
		SdkFuncs: make([]*SdkFunc, len(funcs)),
		Typedefs: jstypez.NewTypedefs(nil),
	}
	for idx, fn := range funcs {
		logger.Info(ctx, "parsing js sdk func", map[string]any{"func": fn.SdkFuncSpec().Id})
		parsed := ParseSdkFunc(fn)
		ret.Typedefs.Merge(parsed.Typedefs)
		ret.SdkFuncs[idx] = parsed
	}
	return &ret
}
