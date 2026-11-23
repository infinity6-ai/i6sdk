package relevancefashiondashdisplaymetricsstatus_test

import (
	"encoding/json"
	"testing"

	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashdisplaymetricsstatus"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
	"github.com/stretchr/testify/require"
)

func decodeInput(t *testing.T, js string) *relevancefashiondashdisplaymetricsstatus.Func {
	t.Helper()
	fn := &relevancefashiondashdisplaymetricsstatus.Func{}
	require.NoError(t, json.Unmarshal([]byte(js), fn.SdkFuncDataRefs().Input))
	return fn
}

func TestUnitRelevanceFashionDashDisplayMetricsStatusInput(t *testing.T) {
	fn := decodeInput(t, `{"dataset":"xpto2","params":{"min_year_month":"2026-07","max_year_month":"2026-11"}}`)
	var expectedDataset lakedsz.Dataset
	require.NoError(t, json.Unmarshal([]byte(`"xpto2"`), &expectedDataset))
	require.Equal(t, expectedDataset, fn.Req.Dataset)
	require.Equal(t, "2026-07", fn.Req.Params.MinYearMonth)
	require.Equal(t, "2026-11", fn.Req.Params.MaxYearMonth)
}

func TestUnitRelevanceFashionDashDisplayMetricsStatusOutput(t *testing.T) {
	fn := &relevancefashiondashdisplaymetricsstatus.Func{}
	js := `{"rows":[{"src_display_id":"NO_DISPLAY","recraws":1241,"recshows":587,"clicks":14659,"cart_adds":4421,"purchases":1526,"revenue":709480.8,"unique_sessions":442310,"last_display_activity":"2026-11-09T19:08:32Z","ctr":24.9727427597956,"conversion_rate":0.104099870386793,"uplift":0.0,"uplift_color":"YELLOW"},{"src_display_id":"312ee60c-b9d1-4895-878b-5220b8fd84bc","recraws":628,"recshows":192,"clicks":11,"cart_adds":7,"purchases":4,"revenue":11943.6,"unique_sessions":184,"last_display_activity":"2026-11-06T21:24:55Z","ctr":0.0572916666666667,"conversion_rate":0.363636363636364,"uplift":2.49314905278208,"uplift_color":"GREEN"},{"src_display_id":"9a07a931-9e4d-46bf-93db-710708024f50","recraws":1168,"recshows":399,"clicks":48,"cart_adds":48,"purchases":18,"revenue":23377.9,"unique_sessions":864,"last_display_activity":"2026-11-09T19:07:59Z","ctr":0.120300751879699,"conversion_rate":0.375,"uplift":2.60230996068152,"uplift_color":"GREEN"}]}`
	require.NoError(t, json.Unmarshal([]byte(js), fn.SdkFuncDataRefs().Output))

	rows := fn.Resp.Results.Rows
	require.Len(t, rows, 3)

	require.Equal(t, relevancefashiondashdisplaymetricsstatus.DisplayNoDisplay, rows[0].SrcDisplayId)
	require.Equal(t, 1241.0, rows[0].RecRaws)
	require.Equal(t, 587.0, rows[0].RecShows)
	require.Equal(t, 14659.0, rows[0].Clicks)
	require.Equal(t, 4421.0, rows[0].CartAdds)
	require.Equal(t, 1526.0, rows[0].Purchases)
	require.InDelta(t, 709480.8, rows[0].Revenue, 1e-6)
	require.Equal(t, 442310.0, rows[0].UniqueSessions)
	require.Equal(t, "2026-11-09T19:08:32Z", rows[0].LastDisplayActivity)
	require.InDelta(t, 24.9727427597956, rows[0].Ctr, 1e-9)
	require.InDelta(t, 0.104099870386793, rows[0].ConversionRate, 1e-9)
	require.InDelta(t, 0.0, rows[0].Uplift, 1e-9)
	require.Equal(t, relevancefashiondashdisplaymetricsstatus.UpliftColorYellow, rows[0].UpliftColor)

	require.Equal(t, "312ee60c-b9d1-4895-878b-5220b8fd84bc", rows[1].SrcDisplayId)
	require.Equal(t, relevancefashiondashdisplaymetricsstatus.UpliftColorGreen, rows[1].UpliftColor)
	require.InDelta(t, 2.49314905278208, rows[1].Uplift, 1e-9)

	require.Equal(t, "9a07a931-9e4d-46bf-93db-710708024f50", rows[2].SrcDisplayId)
	require.Equal(t, relevancefashiondashdisplaymetricsstatus.UpliftColorGreen, rows[2].UpliftColor)
	require.InDelta(t, 2.60230996068152, rows[2].Uplift, 1e-9)
}

func TestUnitRelevanceFashionDashDisplayMetricsStatusSpec(t *testing.T) {
	spec := (&relevancefashiondashdisplaymetricsstatus.Func{}).SdkFuncSpec()
	require.Equal(t, "product-relevance-fashion-dash-display-metrics-status", spec.Id)
	require.Equal(t, "v1", spec.Version)
	require.Equal(t, sdkfuncz.TypClientQuery, spec.Typ)
	require.NotNil(t, spec.Desc)
}

func TestUnitRelevanceFashionDashDisplayMetricsStatusDataRefs(t *testing.T) {
	refs := (&relevancefashiondashdisplaymetricsstatus.Func{}).SdkFuncDataRefs()
	require.NotNil(t, refs)
	require.NotNil(t, refs.Input)
	require.NotNil(t, refs.Input.Object)
	inputProps := refs.Input.Object(true)
	require.Contains(t, inputProps, "dataset")
	require.Contains(t, inputProps, "params")

	require.NotNil(t, refs.Partitions)
	require.NotNil(t, refs.Partitions.Object)
	partitionProps := refs.Partitions.Object(true)
	require.Contains(t, partitionProps, "store_id")

	require.NotNil(t, refs.Output)
	require.NotNil(t, refs.Output.Object)
	outputProps := refs.Output.Object(true)
	require.Contains(t, outputProps, "rows")
}
