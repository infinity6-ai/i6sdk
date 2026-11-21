package userapiauthme_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthme"
	"github.com/stretchr/testify/require"
)

func TestUnitUserAuthMeApi(t *testing.T) {
	api := &userapiauthme.Api{
		Req: userapiauthme.Req{
			AuthHeader: "Bearer eyJhbGci...",
		},
		Resp: userapiauthme.Resp{
			UserId:        "u-12345",
			EmailRedacted: "al***@infinity6.ai",
			Providers:     []string{"google"},
			CreatedAt:     1792900000,
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapiauthme.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	spec := api.ApiSpec()
	require.Equal(t, "user-auth-me", spec.Id)
}
