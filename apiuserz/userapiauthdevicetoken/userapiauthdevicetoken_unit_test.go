package userapiauthdevicetoken_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthdevicetoken"
	"github.com/stretchr/testify/require"
)

func TestUnitUserAuthDeviceTokenApi(t *testing.T) {
	api := &userapiauthdevicetoken.Api{
		Req: userapiauthdevicetoken.Req{
			DeviceCode: "dev-code-secret-123",
		},
		Resp: userapiauthdevicetoken.Resp{
			UserAccessToken:  "jwt-token-abc",
			UserRefreshToken: "rt-session-def",
			UserId:           "u-1234",
			ExpiresAt:        1792933262,
			Status:           "authorized",
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapiauthdevicetoken.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	spec := api.ApiSpec()
	require.Equal(t, "user-auth-device-token", spec.Id)
}
