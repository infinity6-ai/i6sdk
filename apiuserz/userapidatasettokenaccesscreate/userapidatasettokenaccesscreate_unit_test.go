package userapidatasettokenaccesscreate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenaccesscreate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/stretchr/testify/require"
)

func TestUnitUserDatasetTokenAccessCreateApi(t *testing.T) {
	api := &userapidatasettokenaccesscreate.Api{
		Req: &userapidatasettokenaccesscreate.Req{
			Params: &userapidatasettokenaccesscreate.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
			Headers: &userapidatasettokenaccesscreate.ReqHeaders{
				Authorization: "Bearer token",
			},
			Payload: &userapidatasettokenaccesscreate.ReqPayload{
				Grants: grants.Read,
			},
		},
		Resp: &userapidatasettokenaccesscreate.Resp{
			Payload: &userapidatasettokenaccesscreate.RespPayload{
				AccessToken: "token_123",
				Dataset:     lakedsz.MustParse("myds"),
				Grants:      grants.Read,
				ExpiresAt:   123456789,
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapidatasettokenaccesscreate.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "user-dataset-token-access-create", spec.Id)
	require.Equal(t, "Create Dataset Access Token", spec.Desc().Name)
}
