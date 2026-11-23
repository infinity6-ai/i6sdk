package forgepybundle_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/infinity6-ai/gox/routez/apispecz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/infinity6-ai/i6sdk/forgez/forgepy/forgepyapi"
	"github.com/infinity6-ai/i6sdk/forgez/forgepy/forgepybundle"
	"github.com/infinity6-ai/i6sdk/forgez/forgepy/forgepysdkfunc"
	"github.com/infinity6-ai/i6sdk/mold/moldsdkfunc"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/inferz/inferzrelevancefashionfbt"
	"github.com/stretchr/testify/require"
)

func TestUnitWrite(t *testing.T) {
	ctx := t.Context()

	t.Run("empty options writes header and resp container", func(t *testing.T) {
		got := forgepybundle.Write(ctx, forgepybundle.Options{})
		require.Contains(t, got, "from __future__ import annotations")
		require.Contains(t, got, "class Resp(Generic[T]):")
		require.Contains(t, got, "# typedefs")
	})

	t.Run("with sdkfunc and api", func(t *testing.T) {
		api := forgepyapi.NewSdk().ParseApi(ctx, routezsamplefraction.NewService().Api())
		fn := forgepysdkfunc.NewSdk().ParseSdkFunc(ctx, &inferzrelevancefashionfbt.Func{})
		got := forgepybundle.Write(ctx, forgepybundle.Options{
			Apis:     []*apispecz.Api{api},
			SdkFuncs: []*moldsdkfunc.SdkFunc{fn},
		})
		require.Contains(t, got, "# typedefs")
		require.Contains(t, got, "# apis")
		require.Contains(t, got, "class Apis:")
		require.Contains(t, got, "def sample_fraction(")
		require.Contains(t, got, "# sdkfuncs")
		require.Contains(t, got, "class SdkFuncs:")
		require.Contains(t, got, "def product_relevance_fashion_fbt_v1(")
	})
}

func TestUnitWriteFiles(t *testing.T) {
	ctx := t.Context()
	api := forgepyapi.NewSdk().ParseApi(ctx, routezsamplefraction.NewService().Api())
	fn := forgepysdkfunc.NewSdk().ParseSdkFunc(ctx, &inferzrelevancefashionfbt.Func{})

	files := forgepybundle.WriteFiles(ctx, forgepybundle.Options{
		Apis:     []*apispecz.Api{api},
		SdkFuncs: []*moldsdkfunc.SdkFunc{fn},
	})

	require.Contains(t, files, "types.py")
	require.Contains(t, files, "apis.py")
	require.Contains(t, files, "sdkfuncs.py")
	require.Contains(t, files, "__init__.py")

	require.Contains(t, files["types.py"], "class ServiceOptions(TypedDict, total=False):")
	require.Contains(t, files["apis.py"], "class Apis:")
	require.Contains(t, files["apis.py"], "def sample_fraction(")
	require.Contains(t, files["sdkfuncs.py"], "class SdkFuncs:")
	require.Contains(t, files["sdkfuncs.py"], "def product_relevance_fashion_fbt_v1(")
	require.Contains(t, files["__init__.py"], "\"Apis\"")
	require.Contains(t, files["__init__.py"], "\"SdkFuncs\"")
}

func TestUnitWriteDir(t *testing.T) {
	ctx := t.Context()
	tmpDir, err := os.MkdirTemp("", "forgepybundle_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	api := forgepyapi.NewSdk().ParseApi(ctx, routezsamplefraction.NewService().Api())
	fn := forgepysdkfunc.NewSdk().ParseSdkFunc(ctx, &inferzrelevancefashionfbt.Func{})

	err = forgepybundle.WriteDir(ctx, tmpDir, forgepybundle.Options{
		Apis:     []*apispecz.Api{api},
		SdkFuncs: []*moldsdkfunc.SdkFunc{fn},
	})
	require.NoError(t, err)

	for _, name := range []string{"types.py", "apis.py", "sdkfuncs.py", "__init__.py"} {
		path := filepath.Join(tmpDir, name)
		content, err := os.ReadFile(path)
		require.NoError(t, err)
		require.NotEmpty(t, content)
	}
}
