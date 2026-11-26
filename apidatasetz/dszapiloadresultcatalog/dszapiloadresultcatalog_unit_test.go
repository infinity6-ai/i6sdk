package dszapiloadresultcatalog_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapiloadresultcatalog"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
	"github.com/infinity6-ai/i6sdk/sdkcorez/resultz/relevancefashionresultz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszResultRelevanceFashionCatalogApi(t *testing.T) {
	store := storeid.MustParse("store-1")

	api := &dszapiloadresultcatalog.Api{
		Req: &dszapiloadresultcatalog.Req{
			Params: &dszapiloadresultcatalog.ReqParams{
				Dataset: lakedsz.MustParse("my-dataset"),
			},
			Query: &dszapiloadresultcatalog.ReqQuery{
				StoreId: store,
			},
		},
		Resp: &dszapiloadresultcatalog.Resp{
			Payload: &dszapiloadresultcatalog.RespPayload{
				Data: &relevancefashionresultz.CatalogSimple{
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
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapiloadresultcatalog.Api{})
	require.Equal(t, api.Req.Params.Dataset.Get(), apiCopy.Req.Params.Dataset.Get())
	require.Equal(t, store, apiCopy.Req.Query.StoreId)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-result-relevance-fashion-catalog", spec.Id)
	require.Equal(t, "/api/ds/dataset/{dataset}/results/relevance-fashion-catalog", spec.Path)
}
