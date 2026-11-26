package userapiauthexchange_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthexchange"
	"github.com/stretchr/testify/require"
)

func TestUnitUserAuthExchangeApi(t *testing.T) {
	api := &userapiauthexchange.Api{
		Req: &userapiauthexchange.Req{
			Payload: &userapiauthexchange.ReqPayload{
				Code:         "code_123",
				CodeVerifier: "verifier_456",
			},
		},
		Resp: &userapiauthexchange.Resp{
			Payload: &userapiauthexchange.RespPayload{
				UserAccessToken:       "jwt_token",
				UserRefreshToken:      "ref_token",
				UserId:                "u_123",
				ExpiresAt:             123456789,
				RefreshTokenExpiresAt: 123499999,
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapiauthexchange.Api{})
	require.Equal(t, api.Req.Payload.Code, apiCopy.Req.Payload.Code)

	spec := api.ApiSpec()
	require.Equal(t, "user-auth-exchange", spec.Id)
	require.Equal(t, "OAuth 2.0 PKCE Code Exchange", spec.Desc().Name)
}
