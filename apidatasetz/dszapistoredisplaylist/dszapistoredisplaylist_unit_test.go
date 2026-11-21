package dszapistoredisplaylist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapistoredisplaylist"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/dszstore"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storedisplayid"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
)

func TestUnitDszStoreDisplayListApi(t *testing.T) {
	api := &dszapistoredisplaylist.Api{
		Req: &dszapistoredisplaylist.DszStoreDisplayListReq{
			Dataset: lakedsz.MustParse("myds"),
			StoreId: storeid.MustParse("store1"),
		},
		Resp: &dszapistoredisplaylist.DszStoreDisplayListResp{
			Dataset: lakedsz.MustParse("myds"),
			StoreId: storeid.MustParse("store1"),
			Displays: []dszstore.DszStoreDisplay{
				{
					Id:   storedisplayid.MustParse("display1"),
					Name: "Display One",
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapistoredisplaylist.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-store-display-list", spec.Id)
}
