package relevancefashiondashimpactedmetricsbysku_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/productz/relevance/dash/relevancefashiondashimpactedmetricsbysku"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
)

func TestUnitRelevanceFashionDashImpactedMetricsBySku(t *testing.T) {
	fn := &relevancefashiondashimpactedmetricsbysku.Func{
		Req: relevancefashiondashimpactedmetricsbysku.Req{
			Dataset: lakedsz.MustParse("my-dataset"),
			DashPartitions: relevancefashiondashimpactedmetricsbysku.Partitions{
				StoreId: "store-123",
			},
			Params: relevancefashiondashimpactedmetricsbysku.Params{
				MinYearMonth: "2026-07",
				MaxYearMonth: "2026-09",
			},
		},
		Resp: relevancefashiondashimpactedmetricsbysku.Resp{
			Results: relevancefashiondashimpactedmetricsbysku.Result{
				Rows: []*relevancefashiondashimpactedmetricsbysku.Row{
					{
						SkuId:            "0",
						ShowcaseCategory: strPtr("shirt"),
						Description:      "Seamless Ribbed Sports Bra - Brown",
						Recomendations:   205,
						Clicks:           5,
						CartAdds:         3,
						Purchases:        3,
						Revenue:          1368.01,
						Participation:    0.12,
						ConversionRate:   0.6,
						EngagementRate:   0.024,
					},
					{
						SkuId:            "1",
						ShowcaseCategory: nil,
						Description:      "Women's Seamless Ribbed Leggings - Brown",
						Recomendations:   226,
						Clicks:           6,
						CartAdds:         5,
						Purchases:        3,
						Revenue:          103.5,
						Participation:    0.08,
						ConversionRate:   0.5,
						EngagementRate:   0.026,
					},
				},
			},
		},
	}

	fnCopy := jsonz.MustCopy(fn, &relevancefashiondashimpactedmetricsbysku.Func{})

	require.Equal(t, fn.Req, fnCopy.Req)
	require.Equal(t, fn.Resp, fnCopy.Resp)
	require.Equal(t, "my-dataset", fnCopy.Req.Dataset.String())
	require.Equal(t, "store-123", fnCopy.Req.DashPartitions.StoreId)
	require.Equal(t, "2026-07", fnCopy.Req.Params.MinYearMonth)
	require.Equal(t, "2026-09", fnCopy.Req.Params.MaxYearMonth)
	require.Len(t, fnCopy.Resp.Results.Rows, 2)
	require.Equal(t, "shirt", *fnCopy.Resp.Results.Rows[0].ShowcaseCategory)
	require.Nil(t, fnCopy.Resp.Results.Rows[1].ShowcaseCategory)
	require.Equal(t, fn.Resp.Results.Rows[0].Revenue, fnCopy.Resp.Results.Rows[0].Revenue)
	require.Equal(t, fn.Resp.Results.Rows[0].EngagementRate, fnCopy.Resp.Results.Rows[0].EngagementRate)
}

func TestUnitRelevanceFashionDashImpactedMetricsBySkuSpec(t *testing.T) {
	fn := &relevancefashiondashimpactedmetricsbysku.Func{}
	spec := fn.SdkFuncSpec()
	require.Equal(t, "product-relevance-fashion-dash-impacted-metrics-by-sku", spec.Id)
	require.Equal(t, "v1", spec.Version)
	require.Equal(t, sdkfuncz.TypClientQuery, spec.Typ)
	require.NotNil(t, spec.Desc)
}

func TestUnitRelevanceFashionDashImpactedMetricsBySkuDataRefs(t *testing.T) {
	fn := &relevancefashiondashimpactedmetricsbysku.Func{}
	refs := fn.SdkFuncDataRefs()
	require.NotNil(t, refs)
	require.NotNil(t, refs.Input)
	require.NotNil(t, refs.Partitions)
	require.NotNil(t, refs.Output)

	require.NotNil(t, refs.Input.Object)
	inputProps := refs.Input.Object(true)
	require.Contains(t, inputProps, "dataset")
	require.Contains(t, inputProps, "params")

	require.NotNil(t, refs.Partitions.Object)
	partitionProps := refs.Partitions.Object(true)
	require.Contains(t, partitionProps, "store_id")

	require.NotNil(t, refs.Output.Object)
	outputProps := refs.Output.Object(true)
	require.Contains(t, outputProps, "rows")
}

func strPtr(s string) *string {
	return &s
}
