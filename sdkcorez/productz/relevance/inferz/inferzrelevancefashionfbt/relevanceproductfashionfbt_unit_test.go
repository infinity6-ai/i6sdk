package inferzrelevancefashionfbt_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/inferz/inferzrelevancefashionfbt"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
	"github.com/stretchr/testify/require"
)

func TestUnitInferzRelevanceFashionFbt(t *testing.T) {
	fn := &inferzrelevancefashionfbt.Func{
		Req: &inferzrelevancefashionfbt.Req{
			Input: &inferzrelevancefashionfbt.ReqInput{
				Dataset: lakedsz.MustParse("my-dataset"),
				Skus: []*inferzrelevancefashionfbt.Sku{
					{
						SkuId: "sku-1",
					},
					{
						SkuId: "sku-2",
					},
				},
			},
			Partitions: &inferzrelevancefashionfbt.ReqPartitions{
				StoreId: "store-123",
			},
		},
		Resp: &inferzrelevancefashionfbt.Resp{
			Output: &inferzrelevancefashionfbt.RespOutput{
				Results: []*inferzrelevancefashionfbt.Result{
					{
						Skus:   []*inferzrelevancefashionfbt.Sku{{SkuId: "sku-3"}},
						Rating: 0.95,
					},
					{
						Skus:   []*inferzrelevancefashionfbt.Sku{{SkuId: "sku-4"}},
						Rating: 0.85,
					},
				},
			},
		},
	}

	fnCopy := jsonz.MustCopy(fn, &inferzrelevancefashionfbt.Func{})

	require.Equal(t, fn.Req.Input.Dataset, fnCopy.Req.Input.Dataset)
	require.Equal(t, "my-dataset", fnCopy.Req.Input.Dataset.String())
	require.Len(t, fnCopy.Req.Input.Skus, 2)
	require.Equal(t, "sku-1", fnCopy.Req.Input.Skus[0].SkuId)
	require.Equal(t, "sku-2", fnCopy.Req.Input.Skus[1].SkuId)
	require.Equal(t, "store-123", fnCopy.Req.Partitions.StoreId)
	require.Len(t, fnCopy.Resp.Output.Results, 2)
	require.Equal(t, "sku-3", fnCopy.Resp.Output.Results[0].Skus[0].SkuId)
	require.Equal(t, 0.95, fnCopy.Resp.Output.Results[0].Rating)
	require.Equal(t, "sku-4", fnCopy.Resp.Output.Results[1].Skus[0].SkuId)
	require.Equal(t, 0.85, fnCopy.Resp.Output.Results[1].Rating)
}

func TestUnitInferzRelevanceFashionFbtSpec(t *testing.T) {
	fn := &inferzrelevancefashionfbt.Func{}
	spec := fn.SdkFuncSpec()
	require.Equal(t, "product-relevance-fashion-fbt", spec.Id)
	require.Equal(t, "v1", spec.Version)
	require.Equal(t, sdkfuncz.TypClientInfer, spec.Typ)
	require.NotNil(t, spec.Desc)
}

func TestUnitInferzRelevanceFashionFbtDataRefs(t *testing.T) {
	fn := &inferzrelevancefashionfbt.Func{}
	refs := fn.SdkFuncDataRefs()
	require.NotNil(t, refs)
	require.NotNil(t, refs.Input)
	require.NotNil(t, refs.Partitions)
	require.NotNil(t, refs.Output)

	require.NotNil(t, refs.Input.Object)
	inputProps := refs.Input.Object(false)
	require.Contains(t, inputProps, "dataset")
	require.Contains(t, inputProps, "skus")

	require.NotNil(t, refs.Partitions.Object)
	partitionProps := refs.Partitions.Object(false)
	require.Contains(t, partitionProps, "store_id")

	require.NotNil(t, refs.Output.Object)
	outputProps := refs.Output.Object(false)
	require.Contains(t, outputProps, "results")
}
