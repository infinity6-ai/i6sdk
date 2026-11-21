package dszapistripecheckout_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapistripecheckout"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

func TestUnitStripeCheckout(t *testing.T) {
	api := &dszapistripecheckout.Api{
		Req: dszapistripecheckout.Req{
			DatasetName:   "my-dataset",
			Plan:          "pro_monthly",
			SuccessUrl:    "https://x/ok",
			CancelUrl:     "https://x/cancel",
			CampaignState: "c3RhdGU",
			PaymentPhone:  "+5511999999999",
			PaymentEmail:  "pay@x.com",
			Document:      "12345678900",
		},
		Resp: dszapistripecheckout.Resp{Id: "cs_1", Url: "https://checkout/cs_1", Dataset: lakedsz.MustParse("dsabcdefghijklmnopq")},
	}

	apiCopy := jsonz.MustCopy(api, &dszapistripecheckout.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)
}

func TestUnitStripeCheckoutApiSpec(t *testing.T) {
	spec := (&dszapistripecheckout.Api{}).ApiSpec()
	require.Equal(t, "dsz-stripe-checkout", spec.Id)
	require.Equal(t, "POST", spec.Method)
	require.Equal(t, "/api/stripe/checkout", spec.Path)
}

func TestUnitStripeCheckoutApiDataRefs(t *testing.T) {
	api := &dszapistripecheckout.Api{}
	refs := api.ApiDataRefs()

	require.Nil(t, refs.PathParams)
	require.Nil(t, refs.QueryParams)
	require.Nil(t, refs.ReqHeaders)
	require.NotNil(t, refs.ReqBody)
	require.Nil(t, refs.RespHeaders)
	require.NotNil(t, refs.RespBody)

	api.Req = dszapistripecheckout.Req{
		Plan:          "pro_monthly",
		SuccessUrl:    "https://x/ok",
		CampaignState: "c3RhdGU",
	}
	body, err := refs.ReqBody.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, `{
		"dataset_name": "",
		"plan": "pro_monthly",
		"success_url": "https://x/ok",
		"cancel_url": "",
		"campaign_state": "c3RhdGU",
		"payment_phone": "",
		"payment_email": "",
		"document": "",
		"user_exchange_callback_url": ""
	}`, string(body))

	// Reading a body fills Req
	other := &dszapistripecheckout.Api{}
	err = other.ApiDataRefs().ReqBody.UnmarshalJSON([]byte(`{"campaign_state":"abc","plan":"p","cancel_url":"https://x/cancel"}`))
	require.NoError(t, err)
	require.Equal(t, "abc", other.Req.CampaignState)
	require.Equal(t, "p", other.Req.Plan)
	require.Equal(t, "https://x/cancel", other.Req.CancelUrl)

	err = other.ApiDataRefs().ReqBody.UnmarshalJSON([]byte(`{"payment_phone":"1","payment_email":"a@b.c","document":"2"}`))
	require.NoError(t, err)
	require.Equal(t, "1", other.Req.PaymentPhone)
	require.Equal(t, "a@b.c", other.Req.PaymentEmail)
	require.Equal(t, "2", other.Req.Document)

	respJSON := []byte(`{"id":"cs_1","url":"https://checkout/cs_1","dataset":"dsabcdefghijklmnopq"}`)
	require.NoError(t, refs.RespBody.UnmarshalJSON(respJSON))
	require.Equal(t, "cs_1", api.Resp.Id)
	require.Equal(t, "https://checkout/cs_1", api.Resp.Url)
	require.Equal(t, "dsabcdefghijklmnopq", api.Resp.Dataset.String())
	out, err := refs.RespBody.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, string(respJSON), string(out))
}
