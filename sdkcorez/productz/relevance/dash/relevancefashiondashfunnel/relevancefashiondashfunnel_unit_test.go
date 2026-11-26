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
		Req: &relevancefashiondashfunnel.Req{
			Input: &relevancefashiondashfunnel.ReqInput{
				Dataset: lakedsz.MustParse("my-dataset"),
				Params: &relevancefashiondashfunnel.Params{
					MinYearMonth: "2026-07",
					MaxYearMonth: "2026-11",
				},
			},
			Partitions: &relevancefashiondashfunnel.ReqPartitions{
				StoreId: "store-123",
			},
		},
		Resp: &relevancefashiondashfunnel.Resp{
			Output: &relevancefashiondashfunnel.RespOutput{
				Results: relevancefashiondashfunnel.Result{
					Rows: []*relevancefashiondashfunnel.Row{
						{
							EventType: "rec_show",
							FromRec:   "true",
							QtEvents:  "1200",
						},
					},
				},
			},
		},
	}

	fnCopy := jsonz.MustCopy(fn, &relevancefashiondashfunnel.Func{})

	require.Equal(t, fn.Req.Input.Dataset, fnCopy.Req.Input.Dataset)
	require.Equal(t, "my-dataset", fnCopy.Req.Input.Dataset.String())
	require.Equal(t, "store-123", fnCopy.Req.Partitions.StoreId)
	require.Len(t, fnCopy.Resp.Output.Results.Rows, 1)
	require.Equal(t, "rec_show", fnCopy.Resp.Output.Results.Rows[0].EventType)
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
