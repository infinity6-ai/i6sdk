package forgejsbundle_test

import (
	"testing"

	"github.com/infinity6-ai/gox/routez/apispecz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/forgez/forgejs/forgejsapi"
	"go.code.infinity6.ai/sdkgo/forgez/forgejs/forgejsbundle"
	"go.code.infinity6.ai/sdkgo/forgez/forgejs/forgejssdkfunc"
	"go.code.infinity6.ai/sdkgo/mold/moldsdkfunc"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/inferz/inferzrelevancefashionfbt"
)

func TestUnitWrite(t *testing.T) {
	ctx := t.Context()

	t.Run("empty options writes typedefs and classes", func(t *testing.T) {
		got := forgejsbundle.Write(ctx, forgejsbundle.Options{})
		require.Contains(t, got, "// typedefs")
		require.Contains(t, got, "@typedef {import(\"../services/services.js\").Services} Services")
	})

	t.Run("with sdkfunc and api", func(t *testing.T) {
		api := forgejsapi.NewSdk().ParseApi(ctx, routezsamplefraction.NewService().Api())
		fn := forgejssdkfunc.NewSdk().ParseSdkFunc(ctx, &inferzrelevancefashionfbt.Func{})
		got := forgejsbundle.Write(ctx, forgejsbundle.Options{
			Apis:     []*apispecz.Api{api},
			SdkFuncs: []*moldsdkfunc.SdkFunc{fn},
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
