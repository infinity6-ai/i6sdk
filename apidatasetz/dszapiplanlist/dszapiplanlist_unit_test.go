package dszapiplanlist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapiplanlist"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/stretchr/testify/require"
)

func TestUnitDszListPlansApi(t *testing.T) {
	api := &dszapiplanlist.Api{
		Req: &dszapiplanlist.Req{
			Params: &dszapiplanlist.ReqParams{
				Solution: solution.MustParse("relevance"),
			},
		},
		Resp: &dszapiplanlist.Resp{
			Payload: &dszapiplanlist.RespPayload{
				Plans: []dszapiplanlist.Plan{
					{Key: "relevance_growth_monthly_2027"},
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapiplanlist.Api{})
	require.Equal(t, api.Req.Params.Solution, apiCopy.Req.Params.Solution)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-list-plans", spec.Id)
}
