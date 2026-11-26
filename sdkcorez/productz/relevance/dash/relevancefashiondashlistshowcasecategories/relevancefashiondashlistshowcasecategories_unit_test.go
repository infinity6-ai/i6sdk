package relevancefashiondashlistshowcasecategories_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashlistshowcasecategories"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
	"github.com/stretchr/testify/require"
)

func TestUnitRelevanceFashionDashListShowcaseCategories(t *testing.T) {
	cat := "dresses"
	fn := &relevancefashiondashlistshowcasecategories.Func{
		Req: &relevancefashiondashlistshowcasecategories.Req{
			Input: &relevancefashiondashlistshowcasecategories.ReqInput{
				Dataset: lakedsz.MustParse("my-dataset"),
			},
			Partitions: &relevancefashiondashlistshowcasecategories.ReqPartitions{
				StoreId: "store-123",
			},
		},
		Resp: &relevancefashiondashlistshowcasecategories.Resp{
			Output: &relevancefashiondashlistshowcasecategories.RespOutput{
				Results: relevancefashiondashlistshowcasecategories.Result{
					Rows: []*relevancefashiondashlistshowcasecategories.Row{
						{
							ShowcaseCategory: &cat,
						},
					},
				},
			},
		},
	}

	fnCopy := jsonz.MustCopy(fn, &relevancefashiondashlistshowcasecategories.Func{})

	require.Equal(t, fn.Req.Input.Dataset, fnCopy.Req.Input.Dataset)
	require.Equal(t, "my-dataset", fnCopy.Req.Input.Dataset.String())
	require.Equal(t, "store-123", fnCopy.Req.Partitions.StoreId)
	require.Len(t, fnCopy.Resp.Output.Results.Rows, 1)
	require.Equal(t, "dresses", *fnCopy.Resp.Output.Results.Rows[0].ShowcaseCategory)
}

func TestUnitRelevanceFashionDashListShowcaseCategoriesSpec(t *testing.T) {
	fn := &relevancefashiondashlistshowcasecategories.Func{}
	spec := fn.SdkFuncSpec()
	require.Equal(t, "product-relevance-fashion-dash-list-showcase-categories", spec.Id)
	require.Equal(t, "v1", spec.Version)
	require.Equal(t, sdkfuncz.TypClientQuery, spec.Typ)
	require.NotNil(t, spec.Desc)
}

func TestUnitRelevanceFashionDashListShowcaseCategoriesDataRefs(t *testing.T) {
	fn := &relevancefashiondashlistshowcasecategories.Func{}
	refs := fn.SdkFuncDataRefs()
	require.NotNil(t, refs)
	require.NotNil(t, refs.Input)
	require.NotNil(t, refs.Partitions)
	require.NotNil(t, refs.Output)

	require.NotNil(t, refs.Input.Object)
	inputProps := refs.Input.Object(false)
	require.Contains(t, inputProps, "dataset")

	require.NotNil(t, refs.Partitions.Object)
	partitionProps := refs.Partitions.Object(false)
	require.Contains(t, partitionProps, "store_id")

	require.NotNil(t, refs.Output.Object)
	outputProps := refs.Output.Object(false)
	require.Contains(t, outputProps, "results")
}
