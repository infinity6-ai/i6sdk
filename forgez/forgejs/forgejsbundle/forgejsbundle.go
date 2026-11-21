package forgejsbundle

import (
	"context"
	"strings"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/gox/routez/apispecz"
	"github.com/infinity6-ai/gox/schemaz/typeastz"
	"go.code.infinity6.ai/sdkgo/forgez/forgejs/forgejsapi"
	"go.code.infinity6.ai/sdkgo/forgez/forgejs/forgejssdkfunc"
	"go.code.infinity6.ai/sdkgo/forgez/forgejs/forgejsschema"
	"go.code.infinity6.ai/sdkgo/mold/moldsdkfunc"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

type Options struct {
	Apis     []*apispecz.Api
	SdkFuncs []*moldsdkfunc.SdkFunc
}

type jsgenerator struct {
	Apis     []*apispecz.Api
	SdkFuncs []*moldsdkfunc.SdkFunc
	Typedefs typeastz.Typedefs
}

func (js *jsgenerator) Prepare(ctx context.Context, opts Options) {
	js.Apis = opts.Apis
	js.SdkFuncs = opts.SdkFuncs
	for _, element := range js.Apis {
		logger.Info(ctx, "parsing js api", map[string]any{"spec": element.Spec.Id})
		js.Typedefs.Merge(element.Typedefs)
	}
	js.Typedefs.Merge(apispecz.ExtractTypedefs(js.Apis))

	for _, element := range js.SdkFuncs {
		logger.Info(ctx, "parsing js sdk func", map[string]any{"spec": element.Spec.Id})
		js.Typedefs.Merge(element.Typedefs)
	}
	js.Typedefs.Merge(forgejssdkfunc.ExtractTypedefs(js.SdkFuncs))
}

func Write(ctx context.Context, opts Options) string {
	var js jsgenerator
	js.Prepare(ctx, opts)
	var code strings.Builder
	code.WriteString(" //\n // typedefs\n //\n")
	code.WriteString(forgejsschema.WriteModuleTypedefs(ctx, js.Typedefs))
	code.WriteString(" //\n // apis\n //\n")
	code.WriteString(forgejsapi.Write(ctx, js.Apis))
	code.WriteString(" //\n // sdkfuncs\n //\n")
	code.WriteString(forgejssdkfunc.Write(ctx, js.SdkFuncs))
	return code.String()
}
