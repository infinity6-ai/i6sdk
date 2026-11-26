package apicollectionz_test

import (
	"context"
	"testing"

	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/inferz/inferzrelevancefashionfbt"
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

func TestUnitSdkFuncsMethodShadowing(t *testing.T) {
	called := false
	funcs := &apicollectionz.SdkFuncs{
		Invoker: func(ctx context.Context, fn sdkfuncz.Func) (int, error) {
			called = true
			return 200, nil
		},
	}

	code, err := funcs.ProductRelevanceFashionFbtV1(context.Background(), &inferzrelevancefashionfbt.Func{})
	require.NoError(t, err)
	require.Equal(t, 200, code)
	require.True(t, called)
}
