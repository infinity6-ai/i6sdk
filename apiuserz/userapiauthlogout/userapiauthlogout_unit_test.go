package userapiauthlogout_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthlogout"
	"github.com/stretchr/testify/require"
)

func TestUnitUserAuthLogoutApi(t *testing.T) {
	api := &userapiauthlogout.Api{
		Req: &userapiauthlogout.Req{
			Payload: &userapiauthlogout.ReqPayload{
				RefreshToken: "ref_123",
			},
		},
		Resp: &userapiauthlogout.Resp{
			Payload: &userapiauthlogout.RespPayload{
				Success: true,
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapiauthlogout.Api{})
	require.Equal(t, api.Req.Payload.RefreshToken, apiCopy.Req.Payload.RefreshToken)

	spec := api.ApiSpec()
	require.Equal(t, "user-auth-logout", spec.Id)
	require.Equal(t, "User Logout", spec.Desc().Name)
}
