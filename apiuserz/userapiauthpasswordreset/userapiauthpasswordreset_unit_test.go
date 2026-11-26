package userapiauthpasswordreset_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthpasswordreset"
	"github.com/stretchr/testify/require"
)

func TestUnitUserAuthPasswordResetApi(t *testing.T) {
	api := &userapiauthpasswordreset.Api{
		Req: &userapiauthpasswordreset.Req{
			Payload: &userapiauthpasswordreset.ReqPayload{
				Email: "user@test.com",
			},
		},
		Resp: &userapiauthpasswordreset.Resp{
			Payload: &userapiauthpasswordreset.RespPayload{
				Status: "password_reset_email_sent",
				Email:  "user@test.com",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapiauthpasswordreset.Api{})
	require.Equal(t, api.Req.Payload.Email, apiCopy.Req.Payload.Email)

	spec := api.ApiSpec()
	require.Equal(t, "user-auth-password-reset", spec.Id)
	require.Equal(t, "Password Reset Initiation", spec.Desc().Name)
}
