package userapiauthpasswordsignin_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthpasswordsignin"
	"github.com/stretchr/testify/require"
)

func TestUnitUserAuthPasswordSigninApi(t *testing.T) {
	api := &userapiauthpasswordsignin.Api{
		Req: &userapiauthpasswordsignin.Req{
			Payload: &userapiauthpasswordsignin.ReqPayload{
				Email:         "user@test.com",
				Password:      "secretpassword",
				CodeChallenge: "chal_123",
				App:           "cli",
			},
		},
		Resp: &userapiauthpasswordsignin.Resp{
			Payload: &userapiauthpasswordsignin.RespPayload{
				Code: "auth_code_123",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapiauthpasswordsignin.Api{})
	require.Equal(t, api.Req.Payload.Email, apiCopy.Req.Payload.Email)

	spec := api.ApiSpec()
	require.Equal(t, "user-auth-password-signin", spec.Id)
	require.Equal(t, "Email/Password Sign-In", spec.Desc().Name)
}
