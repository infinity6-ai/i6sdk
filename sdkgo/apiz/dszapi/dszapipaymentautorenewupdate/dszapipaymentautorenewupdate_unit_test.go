package dszapipaymentautorenewupdate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapipaymentautorenewupdate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
)

func TestUnitDszPaymentAutoRenewUpdateApi(t *testing.T) {
	t.Run("round trip", func(t *testing.T) {
		api := &dszapipaymentautorenewupdate.Api{
			Req: &dszapipaymentautorenewupdate.Req{
				Dataset:   lakedsz.MustParse("myds"),
				Solution:  solution.Relevance().Solution,
				AutoRenew: true,
			},
		}
		apiCopy := jsonz.MustCopy(api, &dszapipaymentautorenewupdate.Api{})
		require.Equal(t, api.Req, apiCopy.Req)
	})

	t.Run("body schema has auto_renew", func(t *testing.T) {
		api := &dszapipaymentautorenewupdate.Api{
			Req: &dszapipaymentautorenewupdate.Req{
				Dataset:   lakedsz.MustParse("myds"),
				Solution:  solution.Relevance().Solution,
				AutoRenew: false,
			},
		}
		body := string(jsonz.MustFormat(api.ApiDataRefs().ReqBody).Bytes())
		require.JSONEq(t, `{"auto_renew": false}`, body)
	})

	t.Run("spec", func(t *testing.T) {
		spec := (&dszapipaymentautorenewupdate.Api{}).ApiSpec()
		require.Equal(t, "dsz-update-payment-auto-renew", spec.Id)
		require.Equal(t, "POST", spec.Method)
		require.Equal(t, "/api/ds/dataset/{dataset}/solution/{solution}/payment/autorenew", spec.Path)
	})
}
