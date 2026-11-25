package ingestzapieventcatalog_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapieventcatalog"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz/relevancefashiondataz"
	"github.com/stretchr/testify/require"
)

func TestUnitIngestzEventCatalogApi(t *testing.T) {
	api := &ingestzapieventcatalog.Api{
		Req: &ingestzapieventcatalog.Req{
			Params: &ingestzapieventcatalog.ReqParams{
				Dataset: lakedsz.MustParse("my-dataset"),
			},
			Payload: &relevancefashiondataz.Catalog{
				ID:        "obj1",
				ProductId: "prod_1",
				StoreId:   "store_1",
				SkuId:     "sku_1",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &ingestzapieventcatalog.Api{})
	require.Equal(t, api.Req.Params.Dataset.Get(), apiCopy.Req.Params.Dataset.Get())
	require.NotNil(t, apiCopy.Req.Payload)
	require.Equal(t, "obj1", apiCopy.Req.Payload.ID)
	require.Equal(t, "prod_1", apiCopy.Req.Payload.ProductId)

	spec := api.ApiSpec()
	require.Equal(t, "ingestz-event-relevance-fashion-catalog", spec.Id)
	require.Equal(t, "/api/ingest/dataset/{dataset}/event/relevance-fashion-catalog/text.txt", spec.Path)
}
