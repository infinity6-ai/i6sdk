package inferzrelevancefashionfbt_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/inferz/inferzrelevancefashionfbt"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
)

func TestUnitInferzRelevanceFashionFbt(t *testing.T) {
	fn := &inferzrelevancefashionfbt.Func{
		Req: inferzrelevancefashionfbt.Req{
			Dataset: lakedsz.MustParse("my-dataset"),
			Skus: []*inferzrelevancefashionfbt.Sku{
				{
					SkuId: "sku-1",
				},
				{
					SkuId: "sku-2",
				},
			},
			ModelPartitions: inferzrelevancefashionfbt.Partitions{
				StoreId: "store-123",
			},
		},
		Resp: inferzrelevancefashionfbt.Resp{
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
	}

	fnCopy := jsonz.MustCopy(fn, &inferzrelevancefashionfbt.Func{})

	require.Equal(t, fn.Req, fnCopy.Req)
	require.Equal(t, fn.Resp, fnCopy.Resp)
	require.Equal(t, "my-dataset", fnCopy.Req.Dataset.String())
	require.Len(t, fnCopy.Req.Skus, 2)
	require.Equal(t, "sku-1", fnCopy.Req.Skus[0].SkuId)
	require.Equal(t, "sku-2", fnCopy.Req.Skus[1].SkuId)
	require.Equal(t, "store-123", fnCopy.Req.ModelPartitions.StoreId)
	require.Len(t, fnCopy.Resp.Results, 2)
	require.Equal(t, "sku-3", fnCopy.Resp.Results[0].Skus[0].SkuId)
	require.Equal(t, 0.95, fnCopy.Resp.Results[0].Rating)
	require.Equal(t, "sku-4", fnCopy.Resp.Results[1].Skus[0].SkuId)
	require.Equal(t, 0.85, fnCopy.Resp.Results[1].Rating)
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
	inputProps := refs.Input.Object(true)
	require.Contains(t, inputProps, "dataset")
	require.Contains(t, inputProps, "skus")

	require.NotNil(t, refs.Partitions.Object)
	partitionProps := refs.Partitions.Object(true)
	require.Contains(t, partitionProps, "store_id")

	require.NotNil(t, refs.Output.Object)
	outputProps := refs.Output.Object(true)
	require.Contains(t, outputProps, "results")
}
