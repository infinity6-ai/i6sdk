package dszapidatasetsolutions_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapidatasetsolutions"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
)

func TestUnitDszListActiveSolutionsApi(t *testing.T) {
	api := &dszapidatasetsolutions.Api{
		Req: &dszapidatasetsolutions.DszListActiveSolutionsReq{
			Dataset: lakedsz.MustParse("myds"),
		},
		Resp: &dszapidatasetsolutions.DszListActiveSolutionsResp{
			Dataset: lakedsz.MustParse("myds"),
			Solutions: []solution.Solution{
				solution.Relevance().Solution,
				solution.Plat().Solution,
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapidatasetsolutions.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-dataset-solutions", spec.Id)
}
