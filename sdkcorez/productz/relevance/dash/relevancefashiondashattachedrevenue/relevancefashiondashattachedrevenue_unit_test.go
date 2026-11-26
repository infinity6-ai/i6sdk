package relevancefashiondashattachedrevenue_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashattachedrevenue"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
	"github.com/stretchr/testify/require"
)

func TestUnitRelevanceFashionDashAttachedRevenue(t *testing.T) {
	fn := &relevancefashiondashattachedrevenue.Func{
		Req: &relevancefashiondashattachedrevenue.Req{
			Input: &relevancefashiondashattachedrevenue.ReqInput{
				Dataset: lakedsz.MustParse("my-dataset"),
				Params: &relevancefashiondashattachedrevenue.Params{
					MinYearMonth: "2026-07",
					MaxYearMonth: "2026-11",
				},
			},
			Partitions: &relevancefashiondashattachedrevenue.ReqPartitions{
				StoreId: "store-123",
			},
		},
		Resp: &relevancefashiondashattachedrevenue.Resp{
			Output: &relevancefashiondashattachedrevenue.RespOutput{
				Results: relevancefashiondashattachedrevenue.Result{
					Rows: []*relevancefashiondashattachedrevenue.Row{
						{
							AttachedRevenue:   1000.5,
							AttachedPurchases: 50,
							AttachedSkus:      20,
							AvgTicket:         20.01,
						},
					},
				},
			},
		},
	}

	fnCopy := jsonz.MustCopy(fn, &relevancefashiondashattachedrevenue.Func{})

	require.Equal(t, fn.Req.Input.Dataset, fnCopy.Req.Input.Dataset)
	require.Equal(t, "my-dataset", fnCopy.Req.Input.Dataset.String())
	require.Equal(t, "2026-07", fnCopy.Req.Input.Params.MinYearMonth)
	require.Equal(t, "2026-11", fnCopy.Req.Input.Params.MaxYearMonth)
	require.Equal(t, "store-123", fnCopy.Req.Partitions.StoreId)
	require.Len(t, fnCopy.Resp.Output.Results.Rows, 1)
	require.Equal(t, 1000.5, fnCopy.Resp.Output.Results.Rows[0].AttachedRevenue)
}

func TestUnitRelevanceFashionDashAttachedRevenueSpec(t *testing.T) {
	fn := &relevancefashiondashattachedrevenue.Func{}
	spec := fn.SdkFuncSpec()
	require.Equal(t, "product-relevance-fashion-dash-attached-revenue", spec.Id)
	require.Equal(t, "v1", spec.Version)
	require.Equal(t, sdkfuncz.TypClientQuery, spec.Typ)
	require.NotNil(t, spec.Desc)
}

func TestUnitRelevanceFashionDashAttachedRevenueDataRefs(t *testing.T) {
	fn := &relevancefashiondashattachedrevenue.Func{}
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
