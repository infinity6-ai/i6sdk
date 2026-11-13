package userapisigncheck_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparserapi"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/userapi/userapisigncheck"
)

func TestUnitUserSignCheckApi(t *testing.T) {
	ctx := t.Context()
	api := &userapisigncheck.Api{
		Req: userapisigncheck.Req{
			SignatureHeader: "sig-header",
			Signature:       "sig-payload",
			Method:          "POST",
			Path:            "/api/relevance/dataset/ds1/predict",
			Query:           "limit=10",
			Body:            `{"item":123}`,
		},
		Resp: userapisigncheck.Resp{
			Valid:     true,
			ErrorCode: "",
			Summary:   "Request signature is valid and mathematically verified",
			Diagnostics: userapisigncheck.Diagnostics{
				ServerCanonicalString: "canonical",
				Breakdown: &userapisigncheck.Breakdown{
					Method:         "POST",
					Path:           "/api/relevance/dataset/ds1/predict",
					CanonicalQuery: "limit=10",
					BodySha256:     "sha",
					TokenSha256:    "tok",
					Exp:            1792934000,
					Nonce:          "00000000000000000000000000000000",
					UnsignedBody:   false,
				},
				Envelope: &userapisigncheck.Envelope{
					TokenType:            2,
					Exp:                  1792934000,
					TimeRemainingSeconds: 86350,
					IsExpired:            false,
					NonceHex:             "00000000000000000000000000000000",
					IsTimeWindowed:       true,
				},
				EmbeddedToken: &userapisigncheck.EmbeddedToken{
					Dataset:         "ds64l7qauv4l6512kf0c",
					ClientPubKeyHex: "d75a980182b10ab7d54bfed3c964073a0ee172f3daa62325af021a68f707511a",
					Grants:          14,
					Solutions:       3,
					ExpiresAt:       1792936000000,
				},
				CryptoCheck: &userapisigncheck.CryptoCheck{
					VerifiedAgainstPubKey: "d75a980182b10ab7...",
					Match:                 true,
				},
				Hint: "hint",
			},
		},
	}

	spec := api.ApiSpec()
	require.Equal(t, "user-sign-check", spec.Id)
	require.Equal(t, "POST", spec.Method)
	require.Equal(t, "/api/userz/sign/check", spec.Path)
	require.Equal(t, "userz", spec.Desc().Labels["servname"])

	refs := api.ApiDataRefs()
	require.NotNil(t, refs)
	require.NotNil(t, refs.ReqHeaders)
	require.NotNil(t, refs.ReqBody)
	require.NotNil(t, refs.RespBody)

	apiCopy := jsonz.MustCopy(api, &userapisigncheck.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	parsed := jsparserapi.ParseApi(ctx, api)
	require.NotNil(t, parsed)
	require.Equal(t, "user-sign-check", parsed.Spec.Id)
}
