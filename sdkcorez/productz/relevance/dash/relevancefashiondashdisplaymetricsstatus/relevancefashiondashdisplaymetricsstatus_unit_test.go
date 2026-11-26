package relevancefashiondashdisplaymetricsstatus_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashdisplaymetricsstatus"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
	"github.com/stretchr/testify/require"
)

func TestUnitRelevanceFashionDashDisplayMetricsStatus(t *testing.T) {
	fn := &relevancefashiondashdisplaymetricsstatus.Func{
		Req: &relevancefashiondashdisplaymetricsstatus.Req{
			Input: &relevancefashiondashdisplaymetricsstatus.ReqInput{
				Dataset: lakedsz.MustParse("my-dataset"),
				Params: &relevancefashiondashdisplaymetricsstatus.Params{
					MinYearMonth: "2026-07",
					MaxYearMonth: "2026-11",
				},
			},
			Partitions: &relevancefashiondashdisplaymetricsstatus.ReqPartitions{
				StoreId: "store-123",
			},
		},
		Resp: &relevancefashiondashdisplaymetricsstatus.Resp{
			Output: &relevancefashiondashdisplaymetricsstatus.RespOutput{
				Results: relevancefashiondashdisplaymetricsstatus.Result{
					Rows: []*relevancefashiondashdisplaymetricsstatus.Row{
						{
							SrcDisplayId: "display-1",
							Clicks:       100,
							Purchases:    10,
							Revenue:      500.0,
						},
					},
				},
			},
		},
	}

	fnCopy := jsonz.MustCopy(fn, &relevancefashiondashdisplaymetricsstatus.Func{})

	require.Equal(t, fn.Req.Input.Dataset, fnCopy.Req.Input.Dataset)
	require.Equal(t, "my-dataset", fnCopy.Req.Input.Dataset.String())
	require.Equal(t, "store-123", fnCopy.Req.Partitions.StoreId)
	require.Len(t, fnCopy.Resp.Output.Results.Rows, 1)
	require.Equal(t, "display-1", fnCopy.Resp.Output.Results.Rows[0].SrcDisplayId)
}

func TestUnitRelevanceFashionDashDisplayMetricsStatusSpec(t *testing.T) {
	fn := &relevancefashiondashdisplaymetricsstatus.Func{}
	spec := fn.SdkFuncSpec()
	require.Equal(t, "product-relevance-fashion-dash-display-metrics-status", spec.Id)
	require.Equal(t, "v1", spec.Version)
	require.Equal(t, sdkfuncz.TypClientQuery, spec.Typ)
	require.NotNil(t, spec.Desc)
}

func TestUnitRelevanceFashionDashDisplayMetricsStatusDataRefs(t *testing.T) {
	fn := &relevancefashiondashdisplaymetricsstatus.Func{}
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
