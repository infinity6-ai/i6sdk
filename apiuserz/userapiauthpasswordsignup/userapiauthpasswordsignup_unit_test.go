package userapiauthpasswordsignup_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthpasswordsignup"
	"github.com/stretchr/testify/require"
)

func TestUnitUserAuthPasswordSignupApi(t *testing.T) {
	api := &userapiauthpasswordsignup.Api{
		Req: &userapiauthpasswordsignup.Req{
			Payload: &userapiauthpasswordsignup.ReqPayload{
				Email:    "newuser@test.com",
				Password: "secretpassword",
			},
		},
		Resp: &userapiauthpasswordsignup.Resp{
			Payload: &userapiauthpasswordsignup.RespPayload{
				Status: "verification_email_sent",
				Email:  "newuser@test.com",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapiauthpasswordsignup.Api{})
	require.Equal(t, api.Req.Payload.Email, apiCopy.Req.Payload.Email)

	spec := api.ApiSpec()
	require.Equal(t, "user-auth-password-signup", spec.Id)
	require.Equal(t, "Email/Password Registration", spec.Desc().Name)
}
