package jswriter

import (
	"context"
	"strings"

	"go.code.infinity6.ai/platform/logz"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparserapi"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparsersdkfunc"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/genjsz/jswriterapi"
	"go.code.infinity6.ai/sdkgo/genjsz/jswritersdkfunc"
	"go.code.infinity6.ai/sdkgo/genjsz/jswritertypedef"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

type Options struct {
	Apis     []*jsparserapi.Api
	SdkFuncs []*jsparsersdkfunc.SdkFunc
}

type jsgenerator struct {
	Apis     []*jsparserapi.Api
	SdkFuncs []*jsparsersdkfunc.SdkFunc
	Typedefs jstypez.Typedefs
}

func (js *jsgenerator) Prepare(ctx context.Context, opts Options) {
	js.Apis = opts.Apis
	js.SdkFuncs = opts.SdkFuncs
	for _, element := range js.Apis {
		logger.Info(ctx, "parsing js api", map[string]any{"spec": element.Spec.Id})
		js.Typedefs.Merge(element.Typedefs)
	}
	js.Typedefs.Merge(jswriterapi.ExtractTypedefs(js.Apis))

	for _, element := range js.SdkFuncs {
		logger.Info(ctx, "parsing js sdk func", map[string]any{"spec": element.Spec.Id})
		js.Typedefs.Merge(element.Typedefs)
	}
	js.Typedefs.Merge(jswritersdkfunc.ExtractTypedefs(js.SdkFuncs))
}

func Write(ctx context.Context, opts Options) string {
	var js jsgenerator
	js.Prepare(ctx, opts)
	var code strings.Builder
	code.WriteString(" //\n // typedefs\n //\n")
	code.WriteString(jswritertypedef.Write(ctx, js.Typedefs))
	code.WriteString(" //\n // apis\n //\n")
	code.WriteString(jswriterapi.WriteV2(ctx, js.Apis))
	code.WriteString(" //\n // sdkfuncs\n //\n")
	code.WriteString(jswritersdkfunc.Write(ctx, js.SdkFuncs))
	return code.String()
}
