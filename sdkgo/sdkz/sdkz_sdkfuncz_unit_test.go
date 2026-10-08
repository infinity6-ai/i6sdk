package sdkz_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkz"
)

func TestUnitSdkFuncs(t *testing.T) {
	funcs := sdkz.SdkFuncs()
	require.NotEmpty(t, funcs)

	f := sdkz.SdkFunc("product-relevance-fashion-fbt")
	require.NotNil(t, f)
	require.Equal(t, "product-relevance-fashion-fbt", f.InferSpec().Id)
	require.Equal(t, sdkfuncz.TypClientInfer, f.InferSpec().Typ)
	require.Equal(t, "fashionFbt", f.InferSpec().FuncName)

	// Test aliases
	infers := sdkz.Infers()
	require.Equal(t, len(funcs), len(infers))

	infer := sdkz.Infer("product-relevance-fashion-fbt")
	require.Equal(t, f, infer)
}
