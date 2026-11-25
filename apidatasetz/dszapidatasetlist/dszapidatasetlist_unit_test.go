package dszapidatasetlist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapidatasetlist"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/stretchr/testify/require"
)

func TestUnitDszListDatasetsApi(t *testing.T) {
	type testScenario struct {
		name   string
		userId string
		dsId   string
		dsName string
		grants grants.Grant
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()

		api := &dszapidatasetlist.Api{
			Req: &dszapidatasetlist.Req{
				Params: &dszapidatasetlist.Params{
					UserId: s.userId,
				},
			},
			Resp: &dszapidatasetlist.Resp{
				Payload: &dszapidatasetlist.Payload{
					Datasets: []*dszapidatasetlist.Dataset{
						{
							Id:     s.dsId,
							Name:   s.dsName,
							Grants: s.grants,
						},
					},
				},
			},
		}

		apiCopy := jsonz.MustCopy(api, &dszapidatasetlist.Api{})
		require.Equal(t, api.Req, apiCopy.Req)
		require.Equal(t, api.Resp, apiCopy.Resp)

		spec := api.ApiSpec()
		require.Equal(t, "dsz-list-datasets", spec.Id)

		refs := api.ApiDataRefs()
		require.NotNil(t, refs)
		require.NotNil(t, refs.PathParams)
		require.NotNil(t, refs.RespBody)
		require.Nil(t, refs.QueryParams)
		require.Nil(t, refs.ReqHeaders)
		require.Nil(t, refs.ReqBody)
		require.Nil(t, refs.RespHeaders)

		pathParams := refs.PathParams.Object(true)
		require.Contains(t, pathParams, "userId")
		require.Equal(t, "User identifier (id).", pathParams["userId"].Desc.Summary)

		respBody := refs.RespBody.Object(true)
		require.Contains(t, respBody, "datasets")
		require.Equal(t, "List of datasets available to the user.", respBody["datasets"].Desc.Summary)
	}

	t.Run("Valid dataset list API serialization and generated schema", func(t *testing.T) {
		check(t, testScenario{
			name:   "Valid dataset list API serialization and generated schema",
			userId: "user1",
			dsId:   "myds",
			dsName: "My Dataset",
			grants: grants.Read,
		})
	})
}
