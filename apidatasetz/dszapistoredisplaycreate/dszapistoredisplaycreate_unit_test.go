package dszapistoredisplaycreate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapistoredisplaycreate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
	"github.com/stretchr/testify/require"
)

func TestUnitDszStoreDisplayCreateApi(t *testing.T) {
	api := &dszapistoredisplaycreate.Api{
		Req: &dszapistoredisplaycreate.Req{
			Params: &dszapistoredisplaycreate.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
			Payload: &dszapistoredisplaycreate.ReqPayload{
				StoreId: storeid.MustParse("store-1"),
			},
		},
		Resp: &dszapistoredisplaycreate.Resp{
			Payload: &dszapistoredisplaycreate.RespPayload{},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapistoredisplaycreate.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-store-display-create", spec.Id)
}
