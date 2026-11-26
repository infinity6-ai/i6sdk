package userapiauthgoogle_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthgoogle"
	"github.com/stretchr/testify/require"
)

func TestUnitUserAuthGoogleApi(t *testing.T) {
	api := &userapiauthgoogle.Api{
		Req: &userapiauthgoogle.Req{
			Payload: &userapiauthgoogle.ReqPayload{
				IdToken:       "google_id_token",
				CodeChallenge: "chal_123",
				App:           "web",
			},
		},
		Resp: &userapiauthgoogle.Resp{
			Payload: &userapiauthgoogle.RespPayload{
				Code: "auth_code_123",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapiauthgoogle.Api{})
	require.Equal(t, api.Req.Payload.IdToken, apiCopy.Req.Payload.IdToken)

	spec := api.ApiSpec()
	require.Equal(t, "user-auth-google", spec.Id)
	require.Equal(t, "Google One-Tap / ID Token Sign-In", spec.Desc().Name)
}
