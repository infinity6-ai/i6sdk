package dszapistoredisplaylist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapistoredisplaylist"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
	"github.com/stretchr/testify/require"
)

func TestUnitDszStoreDisplayListApi(t *testing.T) {
	api := &dszapistoredisplaylist.Api{
		Req: &dszapistoredisplaylist.Req{
			Params: &dszapistoredisplaylist.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
			Query: &dszapistoredisplaylist.ReqQuery{
				StoreId: storeid.MustParse("store-1"),
			},
		},
		Resp: &dszapistoredisplaylist.Resp{
			Payload: &dszapistoredisplaylist.RespPayload{
				Dataset: lakedsz.MustParse("myds"),
				StoreId: storeid.MustParse("store-1"),
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapistoredisplaylist.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-store-display-list", spec.Id)
}
