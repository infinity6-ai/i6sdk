package relevancefashiondashimpactedmetricsbyshowcasecategory_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashimpactedmetricsbyshowcasecategory"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
	"github.com/stretchr/testify/require"
)

func TestUnitRelevanceFashionDashImpactedMetricsByShowcaseCategory(t *testing.T) {
	cat := "shoes"
	fn := &relevancefashiondashimpactedmetricsbyshowcasecategory.Func{
		Req: &relevancefashiondashimpactedmetricsbyshowcasecategory.Req{
			Input: &relevancefashiondashimpactedmetricsbyshowcasecategory.ReqInput{
				Dataset: lakedsz.MustParse("my-dataset"),
				Params: &relevancefashiondashimpactedmetricsbyshowcasecategory.Params{
					MinYearMonth: "2026-07",
					MaxYearMonth: "2026-11",
				},
			},
			Partitions: &relevancefashiondashimpactedmetricsbyshowcasecategory.ReqPartitions{
				StoreId: "store-123",
			},
		},
		Resp: &relevancefashiondashimpactedmetricsbyshowcasecategory.Resp{
			Output: &relevancefashiondashimpactedmetricsbyshowcasecategory.RespOutput{
				Results: relevancefashiondashimpactedmetricsbyshowcasecategory.Result{
					Rows: []*relevancefashiondashimpactedmetricsbyshowcasecategory.Row{
						{
							ShowcaseCategory: &cat,
							Clicks:           25,
							Revenue:          300.0,
						},
					},
				},
			},
		},
	}

	fnCopy := jsonz.MustCopy(fn, &relevancefashiondashimpactedmetricsbyshowcasecategory.Func{})

	require.Equal(t, fn.Req.Input.Dataset, fnCopy.Req.Input.Dataset)
	require.Equal(t, "my-dataset", fnCopy.Req.Input.Dataset.String())
	require.Equal(t, "store-123", fnCopy.Req.Partitions.StoreId)
	require.Len(t, fnCopy.Resp.Output.Results.Rows, 1)
	require.Equal(t, "shoes", *fnCopy.Resp.Output.Results.Rows[0].ShowcaseCategory)
}

func TestUnitRelevanceFashionDashImpactedMetricsByShowcaseCategorySpec(t *testing.T) {
	fn := &relevancefashiondashimpactedmetricsbyshowcasecategory.Func{}
	spec := fn.SdkFuncSpec()
	require.Equal(t, "product-relevance-fashion-dash-impacted-metrics-by-showcase-category", spec.Id)
	require.Equal(t, "v1", spec.Version)
	require.Equal(t, sdkfuncz.TypClientQuery, spec.Typ)
	require.NotNil(t, spec.Desc)
}

func TestUnitRelevanceFashionDashImpactedMetricsByShowcaseCategoryDataRefs(t *testing.T) {
	fn := &relevancefashiondashimpactedmetricsbyshowcasecategory.Func{}
	refs := fn.SdkFuncDataRefs()
	require.NotNil(t, refs)
	require.NotNil(t, refs.Input)
	require.NotNil(t, refs.Partitions)
	require.NotNil(t, refs.Output)

	require.NotNil(t, refs.Input.Object)
	inputProps := refs.Input.Object(false)
	require.Contains(t, inputProps, "dataset")
	require.Contains(t, inputProps, "params")

	require.NotNil(t, refs.Partitions.Object)
	partitionProps := refs.Partitions.Object(false)
	require.Contains(t, partitionProps, "store_id")

	require.NotNil(t, refs.Output.Object)
	outputProps := refs.Output.Object(false)
	require.Contains(t, outputProps, "results")
}
