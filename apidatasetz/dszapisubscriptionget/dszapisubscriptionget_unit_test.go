package dszapisubscriptionget_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapisubscriptionget"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/stretchr/testify/require"
)

func TestUnitDszGetSubscriptionApi(t *testing.T) {
	api := &dszapisubscriptionget.Api{
		Req: &dszapisubscriptionget.Req{
			Params: &dszapisubscriptionget.ReqParams{
				Dataset:  lakedsz.MustParse("myds"),
				Solution: solution.MustParse("relevance"),
			},
		},
		Resp: &dszapisubscriptionget.Resp{
			Payload: &dszapisubscriptionget.RespPayload{
				PlanKey: "relevance_growth_base_monthly",
				Status:  "active",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapisubscriptionget.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)
	require.Equal(t, api.Resp.Payload.Status, apiCopy.Resp.Payload.Status)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-get-subscription", spec.Id)
}
