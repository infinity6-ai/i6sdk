package userapiauthexchange_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthexchange"
)

func TestUnitUserAuthExchangeApi(t *testing.T) {
	api := &userapiauthexchange.Api{
		Req: userapiauthexchange.Req{
			Code:         "code-123",
			CodeVerifier: "verifier-456",
			Ext:          "state-789",
		},
		Resp: userapiauthexchange.Resp{
			UserAccessToken:  "jwt-token-abc",
			UserRefreshToken: "rt-token-def",
			UserId:           "u-1234",
			ExpiresAt:        1792933262,
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapiauthexchange.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	spec := api.ApiSpec()
	require.Equal(t, "user-auth-exchange", spec.Id)
}
