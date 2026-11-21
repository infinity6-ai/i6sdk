package relevancefashiondashdisplaymetrics_test

import (
	"encoding/json"
	"testing"

	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashdisplaymetrics"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
	"github.com/stretchr/testify/require"
)

func decodeInput(t *testing.T, js string) *relevancefashiondashdisplaymetrics.Func {
	t.Helper()
	fn := &relevancefashiondashdisplaymetrics.Func{}
	require.NoError(t, json.Unmarshal([]byte(js), fn.SdkFuncDataRefs().Input))
	return fn
}

func TestUnitRelevanceFashionDashDisplayMetricsTriState(t *testing.T) {
	base := `"dataset":"xpto2","params":{"min_year_month":"2026-07","max_year_month":"2026-09"`

	fn := decodeInput(t, `{`+base+`}}`)
	require.Equal(t, "2026-07", fn.Req.Params.MinYearMonth)
	require.Nil(t, fn.Req.Params.FromRec)
	require.Nil(t, fn.Req.Params.DisplayId)

	fn = decodeInput(t, `{`+base+`,"from_rec":null,"display_id":null}}`)
	require.Nil(t, fn.Req.Params.FromRec)
	require.Nil(t, fn.Req.Params.DisplayId)

	fn = decodeInput(t, `{`+base+`,"from_rec":true,"display_id":"NO_DISPLAY"}}`)
	require.NotNil(t, fn.Req.Params.FromRec)
	require.True(t, *fn.Req.Params.FromRec)
	require.Equal(t, relevancefashiondashdisplaymetrics.DisplayNoDisplay, *fn.Req.Params.DisplayId)

	fn = decodeInput(t, `{`+base+`,"from_rec":false,"display_id":"d-1"}}`)
	require.NotNil(t, fn.Req.Params.FromRec)
	require.False(t, *fn.Req.Params.FromRec)
	require.Equal(t, "d-1", *fn.Req.Params.DisplayId)
}

func TestUnitRelevanceFashionDashDisplayMetricsOutput(t *testing.T) {
	fn := &relevancefashiondashdisplaymetrics.Func{}
	js := `{"rows":[{"src_display_id":null,"from_rec":0,"recraws":0,"recshows":0,"clicks":6368,"cart_adds":4258,"purchases":2362,"conversion_rate":0.370917085427136},{"src_display_id":null,"from_rec":1,"recraws":36300,"recshows":36300,"clicks":852,"cart_adds":484,"purchases":309,"conversion_rate":0.362676056338028}]}`
	require.NoError(t, json.Unmarshal([]byte(js), fn.SdkFuncDataRefs().Output))
	rows := fn.Resp.Results.Rows
	require.Len(t, rows, 2)
	require.Nil(t, rows[0].SrcDisplayId)
	require.Equal(t, 0, rows[0].FromRec)
	require.Equal(t, 1, rows[1].FromRec)
	require.Equal(t, 36300.0, rows[1].RecRaws)
	require.Equal(t, 852.0, rows[1].Clicks)
}

func TestUnitRelevanceFashionDashDisplayMetricsSpec(t *testing.T) {
	spec := (&relevancefashiondashdisplaymetrics.Func{}).SdkFuncSpec()
	require.Equal(t, "product-relevance-fashion-dash-display-metrics", spec.Id)
	require.Equal(t, "v1", spec.Version)
	require.Equal(t, sdkfuncz.TypClientQuery, spec.Typ)
	require.NotNil(t, spec.Desc)
}

func TestUnitRelevanceFashionDashDisplayMetricsDataRefs(t *testing.T) {
	refs := (&relevancefashiondashdisplaymetrics.Func{}).SdkFuncDataRefs()
	require.NotNil(t, refs)
	require.NotNil(t, refs.Input)
	require.NotNil(t, refs.Input.Object)
	inputProps := refs.Input.Object(true)
	require.Contains(t, inputProps, "dataset")
	require.Contains(t, inputProps, "params")

	require.NotNil(t, refs.Output)
	require.NotNil(t, refs.Output.Object)
	outputProps := refs.Output.Object(true)
	require.Contains(t, outputProps, "rows")
}
