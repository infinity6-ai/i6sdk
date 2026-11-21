package jswriter

import (
	"context"
	"strings"

	"github.com/infinity6-ai/gox/commonz/logz"
	"go.code.infinity6.ai/sdkgo/genjsz/jsapi"
	"go.code.infinity6.ai/sdkgo/genjsz/jssdkfunc"
	"go.code.infinity6.ai/sdkgo/genjsz/jsschema"
	"github.com/infinity6-ai/gox/schemaz/jstypez"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

type Options struct {
	Apis     []*jsapi.Api
	SdkFuncs []*jssdkfunc.SdkFunc
}

type jsgenerator struct {
	Apis     []*jsapi.Api
	SdkFuncs []*jssdkfunc.SdkFunc
	Typedefs jstypez.Typedefs
}

func (js *jsgenerator) Prepare(ctx context.Context, opts Options) {
	js.Apis = opts.Apis
	js.SdkFuncs = opts.SdkFuncs
	for _, element := range js.Apis {
		logger.Info(ctx, "parsing js api", map[string]any{"spec": element.Spec.Id})
		js.Typedefs.Merge(element.Typedefs)
	}
	js.Typedefs.Merge(jsapi.ExtractTypedefs(js.Apis))

	for _, element := range js.SdkFuncs {
		logger.Info(ctx, "parsing js sdk func", map[string]any{"spec": element.Spec.Id})
		js.Typedefs.Merge(element.Typedefs)
	}
	js.Typedefs.Merge(jssdkfunc.ExtractTypedefs(js.SdkFuncs))
}

func Write(ctx context.Context, opts Options) string {
	var js jsgenerator
	js.Prepare(ctx, opts)
	var code strings.Builder
	code.WriteString(" //\n // typedefs\n //\n")
	code.WriteString(jsschema.WriteModuleTypedefs(ctx, js.Typedefs))
	code.WriteString(" //\n // apis\n //\n")
	code.WriteString(jsapi.Write(ctx, js.Apis))
	code.WriteString(" //\n // sdkfuncs\n //\n")
	code.WriteString(jssdkfunc.Write(ctx, js.SdkFuncs))
	return code.String()
}
