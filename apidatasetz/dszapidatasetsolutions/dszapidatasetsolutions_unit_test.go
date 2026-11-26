package dszapidatasetsolutions_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapidatasetsolutions"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/stretchr/testify/require"
)

func TestUnitDszDatasetSolutionsApi(t *testing.T) {
	api := &dszapidatasetsolutions.Api{
		Req: &dszapidatasetsolutions.Req{
			Params: &dszapidatasetsolutions.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
		},
		Resp: &dszapidatasetsolutions.Resp{
			Payload: &dszapidatasetsolutions.RespPayload{
				Dataset: lakedsz.MustParse("myds"),
				Solutions: []solution.Solution{
					solution.MustParse("relevance"),
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapidatasetsolutions.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)
	require.Len(t, apiCopy.Resp.Payload.Solutions, 1)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-dataset-solutions", spec.Id)
}
