package relevancefashiondashcoverage_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashcoverage"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
	"github.com/stretchr/testify/require"
)

func TestUnitRelevanceFashionDashCoverage(t *testing.T) {
	fn := &relevancefashiondashcoverage.Func{
		Req: &relevancefashiondashcoverage.Req{
			Input: &relevancefashiondashcoverage.ReqInput{
				Dataset: lakedsz.MustParse("my-dataset"),
				Params: &relevancefashiondashcoverage.Params{
					MinYearMonth: "2026-07",
					MaxYearMonth: "2026-11",
				},
			},
			Partitions: &relevancefashiondashcoverage.ReqPartitions{
				StoreId: "store-123",
			},
		},
		Resp: &relevancefashiondashcoverage.Resp{
			Output: &relevancefashiondashcoverage.RespOutput{
				Results: relevancefashiondashcoverage.Result{
					Rows: []*relevancefashiondashcoverage.Row{
						{
							RecomendationCoverage: 85.5,
							RecommendedItems:      500,
							RecommendableItems:    600,
						},
					},
				},
			},
		},
	}

	fnCopy := jsonz.MustCopy(fn, &relevancefashiondashcoverage.Func{})

	require.Equal(t, fn.Req.Input.Dataset, fnCopy.Req.Input.Dataset)
	require.Equal(t, "my-dataset", fnCopy.Req.Input.Dataset.String())
	require.Equal(t, "store-123", fnCopy.Req.Partitions.StoreId)
	require.Len(t, fnCopy.Resp.Output.Results.Rows, 1)
	require.Equal(t, 85.5, fnCopy.Resp.Output.Results.Rows[0].RecomendationCoverage)
}

func TestUnitRelevanceFashionDashCoverageSpec(t *testing.T) {
	fn := &relevancefashiondashcoverage.Func{}
	spec := fn.SdkFuncSpec()
	require.Equal(t, "product-relevance-fashion-dash-coverage", spec.Id)
	require.Equal(t, "v1", spec.Version)
	require.Equal(t, sdkfuncz.TypClientQuery, spec.Typ)
	require.NotNil(t, spec.Desc)
}

func TestUnitRelevanceFashionDashCoverageDataRefs(t *testing.T) {
	fn := &relevancefashiondashcoverage.Func{}
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
