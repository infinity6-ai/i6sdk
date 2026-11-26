package dszapipaymentlist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapipaymentlist"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszListPaymentsApi(t *testing.T) {
	api := &dszapipaymentlist.Api{
		Req: &dszapipaymentlist.Req{
			Params: &dszapipaymentlist.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
		},
		Resp: &dszapipaymentlist.Resp{
			Payload: &dszapipaymentlist.RespPayload{
				Payments: []dszapipaymentlist.Payment{
					{Id: "in_123"},
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapipaymentlist.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-list-payments", spec.Id)
}
