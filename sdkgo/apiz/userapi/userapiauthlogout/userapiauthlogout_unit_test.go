package userapiauthlogout_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/userapi/userapiauthlogout"
)

func TestUnitUserAuthLogoutApi(t *testing.T) {
	api := &userapiauthlogout.Api{
		Req: userapiauthlogout.Req{
			RefreshToken: "jti-123.secret-abc",
		},
		Resp: userapiauthlogout.Resp{
			Success: true,
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapiauthlogout.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	spec := api.ApiSpec()
	require.Equal(t, "user-auth-logout", spec.Id)
}
