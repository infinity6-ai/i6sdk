package ingestzapistreamcatalog_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapistreamcatalog"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz/relevancefashiondataz"
	"github.com/stretchr/testify/require"
)

func TestUnitIngestzStreamCatalogApi(t *testing.T) {
	api := &ingestzapistreamcatalog.Api{
		Req: &ingestzapistreamcatalog.Req{
			Params: &ingestzapistreamcatalog.ReqParams{
				Dataset: lakedsz.MustParse("my-dataset"),
			},
			Payload: &ingestzapistreamcatalog.ReqPayload{
				Data: []*relevancefashiondataz.Catalog{
					{
						ID:        "obj1",
						ProductId: "prod_1",
						StoreId:   "store_1",
						SkuId:     "sku_1",
					},
				},
			},
		},
		Resp: &ingestzapistreamcatalog.Resp{
			Payload: &ingestzapistreamcatalog.RespPayload{
				Count: 1,
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &ingestzapistreamcatalog.Api{})
	require.Equal(t, api.Req.Params.Dataset.Get(), apiCopy.Req.Params.Dataset.Get())
	require.Len(t, apiCopy.Req.Payload.Data, 1)
	require.Equal(t, "obj1", apiCopy.Req.Payload.Data[0].ID)
	require.Equal(t, 1, apiCopy.Resp.Payload.Count)

	spec := api.ApiSpec()
	require.Equal(t, "ingestz-stream-relevance-fashion-catalog", spec.Id)
	require.Equal(t, "/api/ingest/dataset/{dataset}/stream/relevance-fashion-catalog", spec.Path)
}
