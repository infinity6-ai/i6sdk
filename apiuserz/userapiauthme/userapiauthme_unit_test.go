package userapiauthme_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthme"
	"github.com/stretchr/testify/require"
)

func TestUnitUserAuthMeApi(t *testing.T) {
	api := &userapiauthme.Api{
		Req: &userapiauthme.Req{
			Headers: &userapiauthme.ReqHeaders{
				Authorization: "Bearer token",
			},
		},
		Resp: &userapiauthme.Resp{
			Payload: &userapiauthme.RespPayload{
				UserId:        "usr_123",
				EmailRedacted: "u***@test.com",
				Providers:     []string{"google"},
				CreatedAt:     123456789,
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapiauthme.Api{})
	require.Equal(t, api.Req.Headers.Authorization, apiCopy.Req.Headers.Authorization)

	spec := api.ApiSpec()
	require.Equal(t, "user-auth-me", spec.Id)
	require.Equal(t, "Current User Profile", spec.Desc().Name)
}
