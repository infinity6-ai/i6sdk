package apicollectionz_test

import (
	"testing"

	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
	"github.com/stretchr/testify/require"
)

func TestUnitSdkFuncs(t *testing.T) {
	funcs := apicollectionz.Funcs()
	require.NotEmpty(t, funcs)

	f := apicollectionz.SdkFunc("product-relevance-fashion-fbt")
	require.NotNil(t, f)
	require.Equal(t, "product-relevance-fashion-fbt", f.SdkFuncSpec().Id)
	require.Equal(t, "v1", f.SdkFuncSpec().Version)
	require.Equal(t, sdkfuncz.TypClientInfer, f.SdkFuncSpec().Typ)

	// Test aliases
	sdkFuncs := apicollectionz.Funcs()
	require.Equal(t, len(funcs), len(sdkFuncs))

	sdkFunc := apicollectionz.SdkFunc("product-relevance-fashion-fbt")
	require.Equal(t, f, sdkFunc)
}
