package jsparser

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"go.code.infinity6.ai/platform/logz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

type Apis struct {
	Apis     []*Api           `json:"apis"`
	Typedefs jstypez.Typedefs `json:"typedefs"`
}

func ParseApis(ctx context.Context, apis []apiz.Api) *Apis {
	ret := Apis{
		Apis:     make([]*Api, len(apis)),
		Typedefs: jstypez.NewTypedefs(nil),
	}
	for idx, api := range apis {
		logger.Info(ctx, "parsring js sdk api", map[string]any{"api": api.ApiSpec().Id})
		parsed := ParseApi(api)
		ret.Typedefs.Merge(parsed.Typedefs)
		ret.Apis[idx] = parsed
	}
	return &ret
}
