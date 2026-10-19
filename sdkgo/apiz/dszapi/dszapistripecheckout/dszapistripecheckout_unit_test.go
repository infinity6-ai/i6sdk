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
			User:        "user-1",
			DatasetName: "my-dataset",
			Plan:        "pro_monthly",
			SuccessUrl:  "https://x/ok",
			CancelUrl:   "https://x/cancel",
			State:       "c3RhdGU",
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
	for _, name := range []string{"user", "dataset_name", "plan", "success_url", "cancel_url", "state"} {
		require.Contains(t, req, name)
	}
	for _, name := range []string{
		"mode", "locale", "allow_promotion_codes", "line_items", "customer_email",
		"subscription_data_trial_period_days", "cnpj_hub", "responsible_name_hub", "responsible_email_hub",
		"id_solution_hub", "id_campanha_hub", "terms_version", "contract_version", "terms_accepted_at",
	} {
		require.NotContains(t, req, name)
	}
	require.Equal(t, "string", req["state"].Typ)

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
		Plan:       "pro_monthly",
		SuccessUrl: "https://x/ok",
		State:      "c3RhdGU",
	}
	body, err := refs.ReqBody.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, `{
		"user": "",
		"dataset_name": "",
		"plan": "pro_monthly",
		"success_url": "https://x/ok",
		"cancel_url": "",
		"state": "c3RhdGU"
	}`, string(body))

	// Reading a body fills Req
	other := &dszapistripecheckout.Api{}
	err = other.ApiDataRefs().ReqBody.UnmarshalJSON([]byte(`{"state":"abc","plan":"p","cancel_url":"https://x/cancel"}`))
	require.NoError(t, err)
	require.Equal(t, "abc", other.Req.State)
	require.Equal(t, "p", other.Req.Plan)
	require.Equal(t, "https://x/cancel", other.Req.CancelUrl)

	respJSON := []byte(`{"id":"cs_1","url":"https://checkout/cs_1"}`)
	require.NoError(t, refs.RespBody.UnmarshalJSON(respJSON))
	require.Equal(t, "cs_1", api.Resp.Id)
	require.Equal(t, "https://checkout/cs_1", api.Resp.Url)
	out, err := refs.RespBody.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, string(respJSON), string(out))
}
