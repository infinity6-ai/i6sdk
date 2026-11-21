package relevancefashiondashfunnel_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashfunnel"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
	"github.com/stretchr/testify/require"
)

func TestUnitRelevanceFashionDashFunnel(t *testing.T) {
	fn := &relevancefashiondashfunnel.Func{
		Req: relevancefashiondashfunnel.Req{
			Dataset: lakedsz.MustParse("my-dataset"),
			DashPartitions: relevancefashiondashfunnel.Partitions{
				StoreId: "store-123",
			},
			Params: relevancefashiondashfunnel.Params{
				MinYearMonth: "2024-01",
				MaxYearMonth: "2024-12",
			},
		},
		Resp: relevancefashiondashfunnel.Resp{
			Results: relevancefashiondashfunnel.Result{
				Rows: []*relevancefashiondashfunnel.Row{
					{
						EventType: "click",
						FromRec:   "true",
						QtEvents:  "150",
					},
					{
						EventType: "purchase",
						FromRec:   "false",
						QtEvents:  "25",
					},
				},
			},
		},
	}

	fnCopy := jsonz.MustCopy(fn, &relevancefashiondashfunnel.Func{})

	require.Equal(t, fn.Req, fnCopy.Req)
	require.Equal(t, fn.Resp, fnCopy.Resp)
	require.Equal(t, "my-dataset", fnCopy.Req.Dataset.String())
	require.Equal(t, "store-123", fnCopy.Req.DashPartitions.StoreId)
	require.Equal(t, "2024-01", fnCopy.Req.Params.MinYearMonth)
	require.Equal(t, "2024-12", fnCopy.Req.Params.MaxYearMonth)
	require.Len(t, fnCopy.Resp.Results.Rows, 2)
	require.Equal(t, "click", fnCopy.Resp.Results.Rows[0].EventType)
	require.Equal(t, "true", fnCopy.Resp.Results.Rows[0].FromRec)
	require.Equal(t, "150", fnCopy.Resp.Results.Rows[0].QtEvents)
	require.Equal(t, "purchase", fnCopy.Resp.Results.Rows[1].EventType)
	require.Equal(t, "false", fnCopy.Resp.Results.Rows[1].FromRec)
	require.Equal(t, "25", fnCopy.Resp.Results.Rows[1].QtEvents)
}

func TestUnitRelevanceFashionDashFunnelSpec(t *testing.T) {
	fn := &relevancefashiondashfunnel.Func{}
	spec := fn.SdkFuncSpec()
	require.Equal(t, "product-relevance-fashion-dash-funnel", spec.Id)
	require.Equal(t, "v1", spec.Version)
	require.Equal(t, sdkfuncz.TypClientQuery, spec.Typ)
	require.NotNil(t, spec.Desc)
}

func TestUnitRelevanceFashionDashFunnelDataRefs(t *testing.T) {
	fn := &relevancefashiondashfunnel.Func{}
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
