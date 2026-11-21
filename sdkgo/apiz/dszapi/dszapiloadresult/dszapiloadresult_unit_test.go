package dszapiloadresult_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapiloadresult"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
	"github.com/infinity6-ai/i6sdk/sdkcorez/resultz/relevancefashionresultz"
)

func TestUnitDszTableFileLoadApi(t *testing.T) {
	store := storeid.MustParse("store-1")

	api := &dszapiloadresult.Api[*relevancefashionresultz.CatalogSimple]{
		ReqDataset: lakedsz.MustParse("my-dataset"),
		ReqRespData: &relevancefashionresultz.CatalogSimple{
			Req: relevancefashionresultz.CatalogSimpleReq{
				StoreId: store,
			},
			Resp: []*relevancefashionresultz.CatalogSimpleData{
				{
					StoreId: store,
					SkuId:   "sku-1",
					SkuDs:   "sku-description",
				},
			},
		},
	}

	// Verify standard JSON copy / two-way binding
	apiCopy := jsonz.MustCopy(api, &dszapiloadresult.Api[*relevancefashionresultz.CatalogSimple]{})
	require.Equal(t, api.ReqDataset.Get(), apiCopy.ReqDataset.Get())
	require.Len(t, apiCopy.ReqRespData.Resp, 1)
	require.Equal(t, api.ReqRespData.Resp[0].StoreId, apiCopy.ReqRespData.Resp[0].StoreId)
	require.Equal(t, api.ReqRespData.Resp[0].SkuId, apiCopy.ReqRespData.Resp[0].SkuId)
	require.Equal(t, api.ReqRespData.Resp[0].SkuDs, apiCopy.ReqRespData.Resp[0].SkuDs)

	// Verify query fields are copied correctly
	require.Equal(t, store, apiCopy.ReqRespData.Req.StoreId)

	// Verify API specifications
	spec := api.ApiSpec()
	require.Equal(t, "dsz-result-relevance-fashion-catalog", spec.Id)
	require.Equal(t, "/api/ds/dataset/{dataset}/results/relevance-fashion-catalog", spec.Path)

	// Verify ApiDataRefs does not overwrite existing query
	refs := api.ApiDataRefs()
	require.NotNil(t, refs)
	require.Equal(t, store, api.ReqRespData.Req.StoreId)

	// Verify ApiDataRefs initializes query when nil
	emptyApi := &dszapiloadresult.Api[*relevancefashionresultz.CatalogSimple]{}
	emptyRefs := emptyApi.ApiDataRefs()
	require.NotNil(t, emptyRefs)
}
