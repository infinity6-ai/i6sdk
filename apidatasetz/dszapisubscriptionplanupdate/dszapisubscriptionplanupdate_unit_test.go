package dszapisubscriptionplanupdate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapisubscriptionplanupdate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/stretchr/testify/require"
)

func TestUnitDszUpdateSubscriptionPlanApi(t *testing.T) {
	api := &dszapisubscriptionplanupdate.Api{
		Req: &dszapisubscriptionplanupdate.Req{
			Params: &dszapisubscriptionplanupdate.ReqParams{
				Dataset:  lakedsz.MustParse("myds"),
				Solution: solution.MustParse("relevance"),
			},
			Payload: &dszapisubscriptionplanupdate.ReqPayload{
				Plan: "relevance_growth_monthly_2027",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapisubscriptionplanupdate.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-update-subscription-plan", spec.Id)
}
