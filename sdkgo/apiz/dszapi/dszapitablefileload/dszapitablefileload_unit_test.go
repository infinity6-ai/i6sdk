package dszapitablefileload_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparserapi"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapitablefileload"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/resultz/relevancefashionresultz"
)

func TestUnitDszTableFileLoadApi(t *testing.T) {
	ctx := t.Context()
	api := &dszapitablefileload.Api[*relevancefashionresultz.CatalogSimple]{
		Req: dszapitablefileload.Req[*relevancefashionresultz.CatalogSimple]{
			Dataset: lakedsz.MustParse("my-dataset"),
		},
		Resp: dszapitablefileload.Resp[*relevancefashionresultz.CatalogSimple]{
			Data: []*relevancefashionresultz.CatalogSimple{
				{
					StoreId: "store-1",
					SkuId:   "sku-1",
					SkuDs:   "sku-description",
				},
			},
		},
	}

	// Verify standard JSON copy / two-way binding
	apiCopy := jsonz.MustCopy(api, &dszapitablefileload.Api[*relevancefashionresultz.CatalogSimple]{})
	require.Equal(t, api.Req.Dataset.Get(), apiCopy.Req.Dataset.Get())
	require.Len(t, apiCopy.Resp.Data, 1)
	require.Equal(t, api.Resp.Data[0].StoreId, apiCopy.Resp.Data[0].StoreId)
	require.Equal(t, api.Resp.Data[0].SkuId, apiCopy.Resp.Data[0].SkuId)
	require.Equal(t, api.Resp.Data[0].SkuDs, apiCopy.Resp.Data[0].SkuDs)

	// Verify API specifications
	spec := api.ApiSpec()
	require.Equal(t, "dsz-table-file-load-catalog-simple", spec.Id)
	require.Equal(t, "/api/ds/dataset/{dataset}/table-file-load/catalog-simple", spec.Path)

	// Verify Schema parsing and mapping
	parsed := jsparserapi.ParseApi(ctx, api)
	require.NotNil(t, parsed)
	require.Equal(t, "dsz-table-file-load-catalog-simple", parsed.Spec.Id)
}
