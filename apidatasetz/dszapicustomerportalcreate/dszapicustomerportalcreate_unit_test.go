package dszapicustomerportalcreate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapicustomerportalcreate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszCreateCustomerPortalApi(t *testing.T) {
	t.Run("round trip", func(t *testing.T) {
		api := &dszapicustomerportalcreate.Api{
			Req: &dszapicustomerportalcreate.Req{
				Dataset:   lakedsz.MustParse("myds"),
				ReturnUrl: "https://app.infinity6.ai/dash",
				Locale:    "pt-BR",
			},
			Resp: &dszapicustomerportalcreate.Resp{
				Url: "https://billing.stripe.com/p/session?secret=xyz",
			},
		}
		apiCopy := jsonz.MustCopy(api, &dszapicustomerportalcreate.Api{})
		require.Equal(t, api.Req, apiCopy.Req)
		require.Equal(t, api.Resp, apiCopy.Resp)
	})

	t.Run("body schema has return_url, locale and url", func(t *testing.T) {
		api := &dszapicustomerportalcreate.Api{
			Req: &dszapicustomerportalcreate.Req{
				Dataset:   lakedsz.MustParse("myds"),
				ReturnUrl: "https://app.infinity6.ai/dash",
				Locale:    "pt-BR",
			},
			Resp: &dszapicustomerportalcreate.Resp{
				Url: "https://billing.stripe.com/p/session?secret=xyz",
			},
		}
		reqBody := string(jsonz.MustFormat(api.ApiDataRefs().ReqBody).Bytes())
		require.JSONEq(t, `{"return_url": "https://app.infinity6.ai/dash", "locale": "pt-BR"}`, reqBody)

		respBody := string(jsonz.MustFormat(api.ApiDataRefs().RespBody).Bytes())
		require.JSONEq(t, `{"url": "https://billing.stripe.com/p/session?secret=xyz"}`, respBody)
	})

	t.Run("spec", func(t *testing.T) {
		spec := (&dszapicustomerportalcreate.Api{}).ApiSpec()
		require.Equal(t, "dsz-create-customer-portal", spec.Id)
		require.Equal(t, "POST", spec.Method)
		require.Equal(t, "/api/ds/dataset/{dataset}/portal", spec.Path)
	})
}
