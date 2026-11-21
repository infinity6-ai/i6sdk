package dszapidatasetlist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapidatasetlist"
	"go.code.infinity6.ai/sdkgo/sdkgo/grants"
)

func TestUnitDszListDatasetsApi(t *testing.T) {
	api := &dszapidatasetlist.Api{
		Req: &dszapidatasetlist.DszListDatasetsReq{
			UserId: "user1",
		},
		Resp: &dszapidatasetlist.DszListDatasetsResp{
			Datasets: []*dszapidatasetlist.Dataset{
				{
					Id:     "myds",
					Name:   "My Dataset",
					Grants: grants.Read,
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapidatasetlist.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-list-datasets", spec.Id)
}
