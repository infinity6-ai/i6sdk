package userapitokencheck_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapitokencheck"
	"github.com/stretchr/testify/require"
)

func TestUnitUserTokenCheckApi(t *testing.T) {
	api := &userapitokencheck.Api{
		Req: &userapitokencheck.Req{
			Payload: &userapitokencheck.ReqPayload{
				Token: "token_123",
			},
		},
		Resp: &userapitokencheck.Resp{
			Payload: &userapitokencheck.RespPayload{
				Valid:     true,
				TokenKind: "user_jwt",
				Summary:   "Token valid",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapitokencheck.Api{})
	require.Equal(t, api.Req.Payload.Token, apiCopy.Req.Payload.Token)
	require.True(t, apiCopy.Resp.Payload.Valid)

	spec := api.ApiSpec()
	require.Equal(t, "user-token-check", spec.Id)
	require.Equal(t, "Universal Token & Signature Diagnostics", spec.Desc().Name)
}
