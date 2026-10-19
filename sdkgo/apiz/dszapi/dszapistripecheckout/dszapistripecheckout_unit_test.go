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
			Plan:                            "pro_monthly",
			User:                            "user-1",
			DatasetName:                     "my-dataset",
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
	ctx := t.Context()
	parsed := jsparserapi.ParseApi(ctx, &dszapistripecheckout.Api{})
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
		"plan", "user", "dataset_name", "customer_email",
		"success_url", "cancel_url", "subscription_data_trial_period_days", "cnpj_hub",
		"responsible_name_hub", "responsible_email_hub", "id_solution_hub", "id_campanha_hub",
		"terms_version", "contract_version", "terms_accepted_at",
	} {
		require.Contains(t, req, name)
	}
	require.NotContains(t, req, "mode")
	require.NotContains(t, req, "locale")
	require.NotContains(t, req, "allow_promotion_codes")
	require.NotContains(t, req, "line_items")
	require.Equal(t, "string", req["plan"].Typ)
	require.Equal(t, "number", req["subscription_data_trial_period_days"].Typ)

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

	api.Req = dszapistripecheckout.Req{
		Plan:          "pro_monthly",
		CustomerEmail: "a@b.com",
	}
	body, err := refs.ReqBody.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, `{
		"plan": "pro_monthly",
		"user": "",
		"dataset_name": "",
		"customer_email": "a@b.com",
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
		"terms_accepted_at": ""
	}`, string(body))

	// Reading a body fills Req
	other := &dszapistripecheckout.Api{}
	err = other.ApiDataRefs().ReqBody.UnmarshalJSON([]byte(`{"plan":"p","customer_email":"c@d.com","subscription_data_trial_period_days":7}`))
	require.NoError(t, err)
	require.Equal(t, "c@d.com", other.Req.CustomerEmail)
	require.Equal(t, 7, other.Req.SubscriptionDataTrialPeriodDays)
	require.Equal(t, "p", other.Req.Plan)

	respJSON := []byte(`{"id":"cs_1","url":"https://checkout/cs_1"}`)
	require.NoError(t, refs.RespBody.UnmarshalJSON(respJSON))
	require.Equal(t, "cs_1", api.Resp.Id)
	require.Equal(t, "https://checkout/cs_1", api.Resp.Url)
	out, err := refs.RespBody.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, string(respJSON), string(out))
}
