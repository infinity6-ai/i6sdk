package dszapistorelist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapistorelist"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/dszstore"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
)

func TestUnitDszListStoresApi(t *testing.T) {
	api := &dszapistorelist.Api{
		Req: &dszapistorelist.DszListStoresReq{
			Dataset: lakedsz.MustParse("myds"),
		},
		Resp: &dszapistorelist.DszListStoresResp{
			Dataset: lakedsz.MustParse("myds"),
			Stores: []dszstore.DszStore{
				{
					Id:   storeid.MustParse("store1"),
					Name: "Store One",
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapistorelist.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-list-stores", spec.Id)
}
