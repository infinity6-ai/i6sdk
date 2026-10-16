package dszapipaymentstatusget_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparserapi"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapipaymentstatusget"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

func TestUnitPaymentStatusGet(t *testing.T) {
	api := &dszapipaymentstatusget.Api{
		Req: dszapipaymentstatusget.Req{User: "user-1", Dataset: lakedsz.MustParse("ds1")},
		Resp: dszapipaymentstatusget.Resp{
			UserId:        "user-1",
			Dataset:       "ds1",
			DatasetName:   "My dataset",
			PaymentId:     "cs_1",
			PaymentType:   "stripe",
			PaymentStatus: "paid",
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapipaymentstatusget.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)
}

func TestUnitPaymentStatusGetParseApi(t *testing.T) {
	ctx := t.Context()
	parsed := jsparserapi.ParseApi(ctx, &dszapipaymentstatusget.Api{})
	require.NotNil(t, parsed)
	require.Equal(t, "dsz-get-payment-status", parsed.Spec.Id)
	require.Equal(t, "apiDszGetPaymentStatus", parsed.Spec.CamelId)
	require.Equal(t, "GET", parsed.Spec.Method)
	require.Equal(t, "/api/ds/payment-status", parsed.Spec.Path)

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

	query := props("apiDszGetPaymentStatusQuery")
	require.Equal(t, "string", query["user"].Typ)
	require.Equal(t, "string", query["dataset"].Typ)

	resp := props("apiDszGetPaymentStatusRespPayload")
	for _, name := range []string{"user_id", "dataset", "dataset_name", "payment_id", "payment_type", "payment_status", "solution"} {
		require.Equal(t, "string", resp[name].Typ, name)
	}
}

func TestUnitPaymentStatusGetApiDataRefs(t *testing.T) {
	api := &dszapipaymentstatusget.Api{}
	refs := api.ApiDataRefs()

	require.Nil(t, refs.PathParams)
	require.NotNil(t, refs.QueryParams)
	require.Nil(t, refs.ReqHeaders)
	require.Nil(t, refs.ReqBody)
	require.Nil(t, refs.RespHeaders)
	require.NotNil(t, refs.RespBody)

	respJSON := []byte(`{"user_id":"user-1","dataset":"ds1","dataset_name":"My dataset","solution":"relevance","payment_id":"cs_1","payment_type":"stripe","payment_status":"paid"}`)
	require.NoError(t, refs.RespBody.UnmarshalJSON(respJSON))
	require.Equal(t, "paid", api.Resp.PaymentStatus)
	require.Equal(t, "ds1", api.Resp.Dataset)
	out, err := refs.RespBody.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, string(respJSON), string(out))
}
