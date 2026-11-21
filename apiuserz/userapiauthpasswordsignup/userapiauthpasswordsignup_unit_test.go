package userapiauthpasswordsignup_test

import (
	"testing"

	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthpasswordsignup"
	"github.com/stretchr/testify/require"
)

func TestUnitUserAuthPasswordSignUpApi(t *testing.T) {
	api := &userapiauthpasswordsignup.Api{
		Req: userapiauthpasswordsignup.Req{
			Email:    "bob@infinity6.ai",
			Password: "securePassword123",
		},
		Resp: userapiauthpasswordsignup.Resp{
			Status: "verification_email_sent",
			Email:  "bob@infinity6.ai",
		},
	}

	spec := api.ApiSpec()
	require.Equal(t, "user-auth-password-signup", spec.Id)
	require.Equal(t, "POST", spec.Method)
	require.Equal(t, "/api/userz/auth/signup/password", spec.Path)
	require.Equal(t, "userz", spec.Desc().Labels["servname"])

	refs := api.ApiDataRefs()
	require.NotNil(t, refs)
	require.NotNil(t, refs.ReqBody)
	require.NotNil(t, refs.RespBody)
}
