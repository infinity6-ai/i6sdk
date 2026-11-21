package relevancefashiondashlistshowcasecategories_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashlistshowcasecategories"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
)

func TestUnitRelevanceFashionDashListShowcaseCategories(t *testing.T) {
	fn := &relevancefashiondashlistshowcasecategories.Func{
		Req: relevancefashiondashlistshowcasecategories.Req{
			Dataset: lakedsz.MustParse("my-dataset"),
			DashPartitions: relevancefashiondashlistshowcasecategories.Partitions{
				StoreId: "store-123",
			},
		},
		Resp: relevancefashiondashlistshowcasecategories.Resp{
			Results: relevancefashiondashlistshowcasecategories.Result{
				Rows: []*relevancefashiondashlistshowcasecategories.Row{
					{
						ShowcaseCategory: strPtr("shirt"),
					},
					{
						ShowcaseCategory: nil,
					},
				},
			},
		},
	}

	fnCopy := jsonz.MustCopy(fn, &relevancefashiondashlistshowcasecategories.Func{})

	require.Equal(t, fn.Req, fnCopy.Req)
	require.Equal(t, fn.Resp, fnCopy.Resp)
	require.Equal(t, "my-dataset", fnCopy.Req.Dataset.String())
	require.Equal(t, "store-123", fnCopy.Req.DashPartitions.StoreId)
	require.Len(t, fnCopy.Resp.Results.Rows, 2)
	require.Equal(t, "shirt", *fnCopy.Resp.Results.Rows[0].ShowcaseCategory)
	require.Nil(t, fnCopy.Resp.Results.Rows[1].ShowcaseCategory)
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
	inputProps := refs.Input.Object(true)
	require.Contains(t, inputProps, "dataset")
	require.NotContains(t, inputProps, "params")

	require.NotNil(t, refs.Partitions.Object)
	partitionProps := refs.Partitions.Object(true)
	require.Contains(t, partitionProps, "store_id")

	require.NotNil(t, refs.Output.Object)
	outputProps := refs.Output.Object(true)
	require.Contains(t, outputProps, "rows")
}


func strPtr(s string) *string {
	return &s
}
