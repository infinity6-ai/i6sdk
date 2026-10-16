package dszapipendingpaymentlist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparserapi"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapipendingpaymentlist"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

func TestUnitPendingPaymentList(t *testing.T) {
	api := &dszapipendingpaymentlist.Api{
		Req: dszapipendingpaymentlist.Req{User: "user-1"},
		Resp: dszapipendingpaymentlist.Resp{
			PendingPayments: []*dszapipendingpaymentlist.PendingPayment{
				{Dataset: lakedsz.MustParse("ds1")},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapipendingpaymentlist.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)
}

func TestUnitPendingPaymentListParseApi(t *testing.T) {
	ctx := t.Context()
	parsed := jsparserapi.ParseApi(ctx, &dszapipendingpaymentlist.Api{})
	require.NotNil(t, parsed)
	require.Equal(t, "dsz-list-pending-payments", parsed.Spec.Id)
	require.Equal(t, "apiDszListPendingPayments", parsed.Spec.CamelId)
	require.Equal(t, "GET", parsed.Spec.Method)
	require.Equal(t, "/api/ds/pending-payments", parsed.Spec.Path)

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

	query := props("apiDszListPendingPaymentsQuery")
	require.Equal(t, "string", query["user"].Typ)

	resp := props("apiDszListPendingPaymentsRespPayload")
	require.Contains(t, resp, "pending_payments")

	item := props("apiDszListPendingPaymentsRespPayloadPendingPaymentsElement")
	require.Equal(t, "string", item["dataset"].Typ, "dataset")
}

func TestUnitPendingPaymentListApiDataRefs(t *testing.T) {
	api := &dszapipendingpaymentlist.Api{}
	refs := api.ApiDataRefs()

	require.Nil(t, refs.PathParams)
	require.NotNil(t, refs.QueryParams)
	require.Nil(t, refs.ReqHeaders)
	require.Nil(t, refs.ReqBody)
	require.Nil(t, refs.RespHeaders)
	require.NotNil(t, refs.RespBody)

	respJSON := []byte(`{"pending_payments":[{"dataset":"ds1"}]}`)
	require.NoError(t, refs.RespBody.UnmarshalJSON(respJSON))
	require.Len(t, api.Resp.PendingPayments, 1)
	require.Equal(t, "ds1", api.Resp.PendingPayments[0].Dataset.String())
	out, err := refs.RespBody.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, string(respJSON), string(out))
}
