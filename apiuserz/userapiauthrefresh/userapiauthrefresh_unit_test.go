package userapiauthrefresh_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthrefresh"
	"github.com/stretchr/testify/require"
)

func TestUnitUserAuthRefreshApi(t *testing.T) {
	api := &userapiauthrefresh.Api{
		Req: &userapiauthrefresh.Req{
			Payload: &userapiauthrefresh.ReqPayload{
				RefreshToken: "ref_123",
			},
		},
		Resp: &userapiauthrefresh.Resp{
			Payload: &userapiauthrefresh.RespPayload{
				UserAccessToken:       "jwt_token",
				UserRefreshToken:      "new_ref_456",
				UserId:                "u_123",
				ExpiresAt:             123456789,
				RefreshTokenExpiresAt: 123499999,
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapiauthrefresh.Api{})
	require.Equal(t, api.Req.Payload.RefreshToken, apiCopy.Req.Payload.RefreshToken)

	spec := api.ApiSpec()
	require.Equal(t, "user-auth-refresh", spec.Id)
	require.Equal(t, "Human User Session Refresh", spec.Desc().Name)
}
