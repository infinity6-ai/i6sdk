package relevancefashiondashattachedrevenue_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/productz/relevance/dash/relevancefashiondashattachedrevenue"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
)

func TestUnitRelevanceFashionDashAttachedRevenue(t *testing.T) {
	fn := &relevancefashiondashattachedrevenue.Func{
		Req: relevancefashiondashattachedrevenue.Req{
			Dataset: lakedsz.MustParse("my-dataset"),
			DashPartitions: relevancefashiondashattachedrevenue.Partitions{
				StoreId: "store-123",
			},
			Params: relevancefashiondashattachedrevenue.Params{
				MinYearMonth: "2024-01",
				MaxYearMonth: "2024-12",
			},
		},
		Resp: relevancefashiondashattachedrevenue.Resp{
			Results: relevancefashiondashattachedrevenue.Result{
				Rows: []*relevancefashiondashattachedrevenue.Row{
					{
						AttachedRevenue:   1500.5,
						AttachedPurchases: 12,
						AttachedSkus:      30,
						AvgTicket:         125.04,
					},
					{
						AttachedRevenue:   10,
						AttachedPurchases: 1,
						AttachedSkus:      2,
						AvgTicket:         10,
					},
				},
			},
		},
	}

	fnCopy := jsonz.MustCopy(fn, &relevancefashiondashattachedrevenue.Func{})

	require.Equal(t, fn.Req, fnCopy.Req)
	require.Equal(t, fn.Resp, fnCopy.Resp)
	require.Equal(t, "my-dataset", fnCopy.Req.Dataset.String())
	require.Equal(t, "store-123", fnCopy.Req.DashPartitions.StoreId)
	require.Equal(t, "2024-01", fnCopy.Req.Params.MinYearMonth)
	require.Equal(t, "2024-12", fnCopy.Req.Params.MaxYearMonth)
	require.Len(t, fnCopy.Resp.Results.Rows, 2)
	require.Equal(t, 1500.5, fnCopy.Resp.Results.Rows[0].AttachedRevenue)
	require.Equal(t, 12.0, fnCopy.Resp.Results.Rows[0].AttachedPurchases)
	require.Equal(t, 30.0, fnCopy.Resp.Results.Rows[0].AttachedSkus)
	require.Equal(t, 125.04, fnCopy.Resp.Results.Rows[0].AvgTicket)
	require.Equal(t, 10.0, fnCopy.Resp.Results.Rows[1].AttachedRevenue)
	require.Equal(t, 2.0, fnCopy.Resp.Results.Rows[1].AttachedSkus)
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
