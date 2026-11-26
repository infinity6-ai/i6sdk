package dszapicustomerportalcreate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapicustomerportalcreate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszCreateCustomerPortalApi(t *testing.T) {
	api := &dszapicustomerportalcreate.Api{
		Req: &dszapicustomerportalcreate.Req{
			Params: &dszapicustomerportalcreate.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
			Payload: &dszapicustomerportalcreate.ReqPayload{
				ReturnUrl: "https://example.com/return",
			},
		},
		Resp: &dszapicustomerportalcreate.Resp{
			Payload: &dszapicustomerportalcreate.RespPayload{
				Url: "https://billing.stripe.com/session/123",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapicustomerportalcreate.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)
	require.Equal(t, api.Resp.Payload.Url, apiCopy.Resp.Payload.Url)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-create-customer-portal", spec.Id)
}
