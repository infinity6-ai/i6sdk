package dszapipaymentautorenewupdate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapipaymentautorenewupdate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/stretchr/testify/require"
)

func TestUnitDszUpdatePaymentAutoRenewApi(t *testing.T) {
	api := &dszapipaymentautorenewupdate.Api{
		Req: &dszapipaymentautorenewupdate.Req{
			Params: &dszapipaymentautorenewupdate.ReqParams{
				Dataset:  lakedsz.MustParse("myds"),
				Solution: solution.MustParse("relevance"),
			},
			Payload: &dszapipaymentautorenewupdate.ReqPayload{
				AutoRenew: true,
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapipaymentautorenewupdate.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-update-payment-auto-renew", spec.Id)
}
