package ingestzapieventevent_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapieventevent"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz/relevancefashiondataz"
	"github.com/stretchr/testify/require"
)

func TestUnitIngestzEventEventApi(t *testing.T) {
	api := &ingestzapieventevent.Api{
		Req: &ingestzapieventevent.Req{
			Params: &ingestzapieventevent.ReqParams{
				Dataset: lakedsz.MustParse("my-dataset"),
			},
			Payload: &relevancefashiondataz.Event{
				ID:        "ev_1",
				StoreId:   "store_1",
				SessionId: "sess_1",
				Click: &relevancefashiondataz.Click{
					SkuId: "sku_1",
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &ingestzapieventevent.Api{})
	require.Equal(t, api.Req.Params.Dataset.Get(), apiCopy.Req.Params.Dataset.Get())
	require.NotNil(t, apiCopy.Req.Payload)
	require.Equal(t, "ev_1", apiCopy.Req.Payload.ID)
	require.Equal(t, "store_1", apiCopy.Req.Payload.StoreId)
	require.Equal(t, "sess_1", apiCopy.Req.Payload.SessionId)
	require.NotNil(t, apiCopy.Req.Payload.Click)
	require.Equal(t, "sku_1", apiCopy.Req.Payload.Click.SkuId)

	spec := api.ApiSpec()
	require.Equal(t, "ingestz-event-relevance-fashion-event", spec.Id)
	require.Equal(t, "/api/ingest/dataset/{dataset}/event/relevance-fashion-event/text.txt", spec.Path)
}
