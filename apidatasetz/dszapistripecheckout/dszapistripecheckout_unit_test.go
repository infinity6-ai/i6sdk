package dszapistripecheckout_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapistripecheckout"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszStripeCheckoutApi(t *testing.T) {
	api := &dszapistripecheckout.Api{
		Req: &dszapistripecheckout.Req{
			Payload: &dszapistripecheckout.ReqPayload{
				DatasetName: "myds",
				Plan:        "plan_123",
				SuccessUrl:  "https://example.com/success",
				CancelUrl:   "https://example.com/cancel",
			},
		},
		Resp: &dszapistripecheckout.Resp{
			Payload: &dszapistripecheckout.RespPayload{
				Id:      "cs_123",
				Url:     "https://checkout.stripe.com/pay/cs_123",
				Dataset: lakedsz.MustParse("myds"),
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapistripecheckout.Api{})
	require.Equal(t, api.Req.Payload.Plan, apiCopy.Req.Payload.Plan)
	require.Equal(t, api.Resp.Payload.Url, apiCopy.Resp.Payload.Url)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-stripe-checkout", spec.Id)
}
