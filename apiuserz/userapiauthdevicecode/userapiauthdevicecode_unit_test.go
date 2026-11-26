package userapiauthdevicecode_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthdevicecode"
	"github.com/stretchr/testify/require"
)

func TestUnitUserAuthDeviceCodeApi(t *testing.T) {
	api := &userapiauthdevicecode.Api{
		Req: &userapiauthdevicecode.Req{
			Payload: &userapiauthdevicecode.ReqPayload{
				App: "cli",
			},
		},
		Resp: &userapiauthdevicecode.Resp{
			Payload: &userapiauthdevicecode.RespPayload{
				DeviceCode:              "dev_123",
				UserCode:                "USER-456",
				VerificationUri:         "https://auth.i6.ai/device",
				VerificationUriComplete: "https://auth.i6.ai/device?code=USER-456",
				ExpiresInSec:            600,
				IntervalSec:             5,
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapiauthdevicecode.Api{})
	require.Equal(t, api.Req.Payload.App, apiCopy.Req.Payload.App)
	require.Equal(t, api.Resp.Payload.DeviceCode, apiCopy.Resp.Payload.DeviceCode)

	spec := api.ApiSpec()
	require.Equal(t, "user-auth-device-code", spec.Id)
	require.Equal(t, "Device Authorization Code Initiation", spec.Desc().Name)
}
