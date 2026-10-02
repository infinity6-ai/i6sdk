package dszapiliststores_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapiliststores"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

func TestUnitDszListStoresApi(t *testing.T) {
	api := &dszapiliststores.DszListStoresApi{
		Req: &dszapiliststores.DszListStoresReq{
			Dataset: lakedsz.MustParse("myds"),
		},
		Resp: &dszapiliststores.DszListStoresResp{
			Dataset: lakedsz.MustParse("myds"),
			Stores: []dszapiliststores.DszStore{
				{
					Id:   "store1",
					Name: "Store One",
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapiliststores.DszListStoresApi{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	parsed := jsparser.ParseApi(api)
	require.NotNil(t, parsed)
	require.Equal(t, "dsz-list-stores", parsed.Spec.Id)
}
