package sdkgo_test

import (
	"testing"

	"github.com/infinity6-ai/i6sdk/sdkgo"
	"github.com/stretchr/testify/require"
)

func TestUnitClient(t *testing.T) {
	c := sdkgo.New()
	require.NotNil(t, c)
	require.NotEmpty(t, c.Version())
	require.NotEmpty(t, c.Config().BaseUrl)
	require.NotEmpty(t, c.Apis())
	require.NotEmpty(t, c.SdkFuncs())

	api := c.Api("user-auth-device-code")
	require.NotNil(t, api)
	require.Equal(t, "user-auth-device-code", api.ApiSpec().Id)

	fn := c.SdkFunc("product-relevance-fashion-fbt")
	require.NotNil(t, fn)
	require.Equal(t, "product-relevance-fashion-fbt", fn.SdkFuncSpec().Id)
}
