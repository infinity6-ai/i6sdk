package jswriter_test

import (
	"testing"

	"github.com/infinity6-ai/gox/routez/apispecz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsapi"
	"go.code.infinity6.ai/sdkgo/genjsz/jssdkfunc"
	"go.code.infinity6.ai/sdkgo/genjsz/jswriter"
	"go.code.infinity6.ai/sdkgo/sdkgo/productz/relevance/inferz/inferzrelevancefashionfbt"
)

func TestUnitWrite(t *testing.T) {
	ctx := t.Context()

	t.Run("empty options writes typedefs and classes", func(t *testing.T) {
		got := jswriter.Write(ctx, jswriter.Options{})
		require.Contains(t, got, "// typedefs")
		require.Contains(t, got, "@typedef {import(\"../services/services.js\").Services} Services")
	})

	t.Run("with sdkfunc and api", func(t *testing.T) {
		api := jsapi.NewSdk().ParseApi(ctx, routezsamplefraction.NewService().Api())
		fn := jssdkfunc.NewSdk().ParseSdkFunc(ctx, &inferzrelevancefashionfbt.Func{})
		got := jswriter.Write(ctx, jswriter.Options{
			Apis:     []*apispecz.Api{api},
			SdkFuncs: []*jssdkfunc.SdkFunc{fn},
		})
		require.Contains(t, got, "// typedefs")
		require.Contains(t, got, "// apis")
		require.Contains(t, got, "class Apis")
		require.Contains(t, got, "sampleFraction")
		require.Contains(t, got, "// sdkfuncs")
		require.Contains(t, got, "class SdkFuncs")
		require.Contains(t, got, "productRelevanceFashionFbtV1")
	})
}
