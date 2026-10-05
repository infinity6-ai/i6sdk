package dszapilistdatasets_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapilistdatasets"
)

func TestUnitDszListDatasetsApi(t *testing.T) {
	api := &dszapilistdatasets.DszListDatasetsApi{
		Req: &dszapilistdatasets.DszListDatasetsReq{
			User: "user1",
		},
		Resp: &dszapilistdatasets.DszListDatasetsResp{
			Datasets: []*dszapilistdatasets.Dataset{
				{
					Id:   "myds",
					Name: "My Dataset",
					Solutions: []*dszapilistdatasets.Solution{
						{
							Id:   "relevance",
							Name: "Relevance Solution",
						},
					},
					Grants: []string{"read"},
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapilistdatasets.DszListDatasetsApi{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	parsed := jsparser.ParseApi(api)
	require.NotNil(t, parsed)
	require.Equal(t, "dsz-list-datasets", parsed.Spec.Id)
}
