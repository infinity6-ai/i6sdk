package userapidatasettokenrefreshcreate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefreshcreate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/stretchr/testify/require"
)

func TestUnitUserDatasetTokenRefreshCreateApi(t *testing.T) {
	api := &userapidatasettokenrefreshcreate.Api{
		Req: &userapidatasettokenrefreshcreate.Req{
			Params: &userapidatasettokenrefreshcreate.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
			Headers: &userapidatasettokenrefreshcreate.ReqHeaders{
				Authorization: "Bearer token",
			},
			Payload: &userapidatasettokenrefreshcreate.ReqPayload{
				Grants: grants.Read,
			},
		},
		Resp: &userapidatasettokenrefreshcreate.Resp{
			Payload: &userapidatasettokenrefreshcreate.RespPayload{
				Jti:       "jti_123",
				Dataset:   lakedsz.MustParse("myds"),
				PrivKey:   "priv_key_b64",
				Grants:    grants.Read,
				ExpiresAt: 123456789,
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapidatasettokenrefreshcreate.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "user-dataset-token-refresh-create", spec.Id)
	require.Equal(t, "Provision Dataset Refresh Token", spec.Desc().Name)
}
