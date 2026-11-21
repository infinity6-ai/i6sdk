package dszapisubscriptionplanupdate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapisubscriptionplanupdate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
)

func TestUnitDszSubscriptionPlanUpdateApi(t *testing.T) {
	t.Run("round trip", func(t *testing.T) {
		api := &dszapisubscriptionplanupdate.Api{
			Req: &dszapisubscriptionplanupdate.Req{
				Dataset:  lakedsz.MustParse("myds"),
				Solution: solution.Relevance().Solution,
				Plan:     "relevance_growth_base_monthly",
			},
		}
		apiCopy := jsonz.MustCopy(api, &dszapisubscriptionplanupdate.Api{})
		require.Equal(t, api.Req, apiCopy.Req)
	})

	t.Run("body schema has plan", func(t *testing.T) {
		api := &dszapisubscriptionplanupdate.Api{
			Req: &dszapisubscriptionplanupdate.Req{
				Dataset:  lakedsz.MustParse("myds"),
				Solution: solution.Relevance().Solution,
				Plan:     "relevance_growth_base_monthly",
			},
		}
		body := string(jsonz.MustFormat(api.ApiDataRefs().ReqBody).Bytes())
		require.JSONEq(t, `{"plan": "relevance_growth_base_monthly"}`, body)
	})

	t.Run("spec", func(t *testing.T) {
		spec := (&dszapisubscriptionplanupdate.Api{}).ApiSpec()
		require.Equal(t, "dsz-update-subscription-plan", spec.Id)
		require.Equal(t, "POST", spec.Method)
		require.Equal(t, "/api/ds/dataset/{dataset}/solution/{solution}/subscription/plan", spec.Path)
	})
}
