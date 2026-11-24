package sdkgo_test

import (
	"testing"

	"github.com/infinity6-ai/i6sdk/sdkgo/sdkgo"
	"github.com/stretchr/testify/require"
)

func TestUnitClient(t *testing.T) {
	c := sdkgo.New()
	require.NotNil(t, c)
	require.NotEmpty(t, c.Version())
	require.Equal(t, "default", c.Config().App)
	require.NotEmpty(t, c.Config().ServUrls["decsuite"])
	require.NotEmpty(t, c.Apis())
	require.NotEmpty(t, c.SdkFuncs())

	api := c.Api("user-auth-device-code")
	require.NotNil(t, api)
	require.Equal(t, "user-auth-device-code", api.ApiSpec().Id)

	fn := c.SdkFunc("product-relevance-fashion-fbt")
	require.NotNil(t, fn)
	require.Equal(t, "product-relevance-fashion-fbt", fn.SdkFuncSpec().Id)
}

func TestUnitClientCustomConfig(t *testing.T) {
	c := sdkgo.New(sdkgo.Config{
		App: "custom-app",
	})
	require.Equal(t, "custom-app", c.Config().App)
}

func TestUnitClientEnvConfig(t *testing.T) {
	t.Setenv("I6_SDK_DECSUITE_URL", "https://custom-decsuite.test")
	t.Setenv("I6_SDK_TFINFER_URL", "https://custom-tfinfer.test")
	t.Setenv("I6_SDK_USERZ_URL", "https://custom-userz.test")

	c := sdkgo.New()
	require.Equal(t, "https://custom-decsuite.test", c.Config().ServUrls["decsuite"])
	require.Equal(t, "https://custom-tfinfer.test", c.Config().ServUrls["tfinfer"])
	require.Equal(t, "https://custom-userz.test", c.Config().ServUrls["userz"])
}
