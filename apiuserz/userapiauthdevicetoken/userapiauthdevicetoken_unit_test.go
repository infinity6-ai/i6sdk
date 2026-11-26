package userapiauthdevicetoken_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthdevicetoken"
	"github.com/stretchr/testify/require"
)

func TestUnitUserAuthDeviceTokenApi(t *testing.T) {
	api := &userapiauthdevicetoken.Api{
		Req: &userapiauthdevicetoken.Req{
			Payload: &userapiauthdevicetoken.ReqPayload{
				DeviceCode: "dev_123",
			},
		},
		Resp: &userapiauthdevicetoken.Resp{
			Payload: &userapiauthdevicetoken.RespPayload{
				UserAccessToken:       "jwt_token",
				UserRefreshToken:      "ref_token",
				UserId:                "u_123",
				ExpiresAt:             123456789,
				RefreshTokenExpiresAt: 123499999,
				Status:                "authorized",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapiauthdevicetoken.Api{})
	require.Equal(t, api.Req.Payload.DeviceCode, apiCopy.Req.Payload.DeviceCode)
	require.Equal(t, api.Resp.Payload.UserAccessToken, apiCopy.Resp.Payload.UserAccessToken)

	spec := api.ApiSpec()
	require.Equal(t, "user-auth-device-token", spec.Id)
	require.Equal(t, "Device Token Polling & Exchange", spec.Desc().Name)
}
