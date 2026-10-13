package dszapistripecheckout_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparserapi"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapistripecheckout"
)

func TestUnitStripeCheckout(t *testing.T) {
	api := &dszapistripecheckout.Api{
		Req: dszapistripecheckout.Req{
			LineItems: []*dszapistripecheckout.LineItem{
				{Price: "price_1", Quantity: 1},
				{Price: "price_2", Quantity: 3},
			},
			CustomerEmail:                   "a@b.com",
			SuccessUrl:                      "https://x/ok",
			CancelUrl:                       "https://x/cancel",
			SubscriptionDataTrialPeriodDays: 14,
			CnpjHub:                         "12345678000199",
			ResponsibleNameHub:              "Ana",
			ResponsibleEmailHub:             "ana@b.com",
			IdSolutionHub:                   "sol-1",
			IdCampanhaHub:                   "camp-1",
			TermsVersion:                    "v1",
			ContractVersion:                 "v2",
			TermsAcceptedAt:                 "2026-01-01T00:00:00Z",
		},
		Resp: dszapistripecheckout.Resp{Id: "cs_1", Url: "https://checkout/cs_1"},
	}

	apiCopy := jsonz.MustCopy(api, &dszapistripecheckout.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)
}

func TestUnitStripeCheckoutParseApi(t *testing.T) {
	parsed := jsparserapi.ParseApi(&dszapistripecheckout.Api{})
	require.NotNil(t, parsed)
	require.Equal(t, "dsz-stripe-checkout", parsed.Spec.Id)
	require.Equal(t, "apiDszStripeCheckout", parsed.Spec.CamelId)
	require.Equal(t, "POST", parsed.Spec.Method)
	require.Equal(t, "/api/stripe/checkout", parsed.Spec.Path)

	tdfMap := make(map[string]*jstypez.Typedef)
	for _, tdf := range parsed.Typedefs.All() {
		tdfMap[tdf.Head.Name] = tdf
	}

	props := func(name string) map[string]jstypez.Tuple {
		tdf, ok := tdfMap[name]
		require.True(t, ok, name)
		ret := make(map[string]jstypez.Tuple)
		for _, p := range tdf.Properties.All() {
			ret[p.Name] = p
		}
		return ret
	}

	req := props("apiDszStripeCheckoutReqPayload")
	for _, name := range []string{
		"mode", "payment_method_types", "line_items", "customer_email", "locale", "expires_at",
		"success_url", "cancel_url", "subscription_data_trial_period_days", "cnpj_hub",
		"responsible_name_hub", "responsible_email_hub", "id_solution_hub", "id_campanha_hub",
		"terms_version", "contract_version", "terms_accepted_at", "allow_promotion_codes",
	} {
		require.Contains(t, req, name)
	}
	require.Equal(t, "string", req["mode"].Typ)
	require.Equal(t, "number", req["subscription_data_trial_period_days"].Typ)
	require.Equal(t, "boolean", req["allow_promotion_codes"].Typ)

	item := props("apiDszStripeCheckoutReqPayloadLineItemsElement")
	require.Equal(t, "string", item["price"].Typ)
	require.Equal(t, "number", item["quantity"].Typ)

	resp := props("apiDszStripeCheckoutRespPayload")
	require.Equal(t, "string", resp["id"].Typ)
	require.Equal(t, "string", resp["url"].Typ)
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

	// The fixed values are always sent, whatever the request holds
	api.Req = dszapistripecheckout.Req{
		LineItems:     []*dszapistripecheckout.LineItem{{Price: "price_1", Quantity: 2}},
		CustomerEmail: "a@b.com",
	}
	body, err := refs.ReqBody.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, `{
		"mode": "subscription",
		"payment_method_types": ["card"],
		"line_items": [{"price": "price_1", "quantity": 2}],
		"customer_email": "a@b.com",
		"locale": "pt-BR",
		"expires_at": "3600",
		"success_url": "",
		"cancel_url": "",
		"subscription_data_trial_period_days": 0,
		"cnpj_hub": "",
		"responsible_name_hub": "",
		"responsible_email_hub": "",
		"id_solution_hub": "",
		"id_campanha_hub": "",
		"terms_version": "",
		"contract_version": "",
		"terms_accepted_at": "",
		"allow_promotion_codes": true
	}`, string(body))

	// Reading a body fills Req and ignores the fixed values
	other := &dszapistripecheckout.Api{}
	err = other.ApiDataRefs().ReqBody.UnmarshalJSON([]byte(`{"mode":"payment","locale":"en","line_items":[{"price":"p","quantity":1}],"customer_email":"c@d.com","subscription_data_trial_period_days":7}`))
	require.NoError(t, err)
	require.Equal(t, "c@d.com", other.Req.CustomerEmail)
	require.Equal(t, 7, other.Req.SubscriptionDataTrialPeriodDays)
	require.Len(t, other.Req.LineItems, 1)
	body, err = other.ApiDataRefs().ReqBody.MarshalJSON()
	require.NoError(t, err)
	require.Contains(t, string(body), `"mode":"subscription"`)
	require.Contains(t, string(body), `"locale":"pt-BR"`)

	respJSON := []byte(`{"id":"cs_1","url":"https://checkout/cs_1"}`)
	require.NoError(t, refs.RespBody.UnmarshalJSON(respJSON))
	require.Equal(t, "cs_1", api.Resp.Id)
	require.Equal(t, "https://checkout/cs_1", api.Resp.Url)
	out, err := refs.RespBody.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, string(respJSON), string(out))
}
