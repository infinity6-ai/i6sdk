package ingestzapistreamevent_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapistreamevent"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz/relevancefashiondataz"
	"github.com/stretchr/testify/require"
)

func TestUnitIngestzStreamEventApi(t *testing.T) {
	api := &ingestzapistreamevent.Api{
		Req: &ingestzapistreamevent.Req{
			Params: &ingestzapistreamevent.ReqParams{
				Dataset: lakedsz.MustParse("my-dataset"),
			},
			Payload: &ingestzapistreamevent.ReqPayload{
				Data: []*relevancefashiondataz.Event{
					{
						ID:        "ev_1",
						StoreId:   "store_1",
						SessionId: "sess_1",
						Click: &relevancefashiondataz.Click{
							SkuId: "sku_1",
						},
					},
				},
			},
		},
		Resp: &ingestzapistreamevent.Resp{
			Payload: &ingestzapistreamevent.RespPayload{
				Count: 1,
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &ingestzapistreamevent.Api{})
	require.Equal(t, api.Req.Params.Dataset.Get(), apiCopy.Req.Params.Dataset.Get())
	require.Len(t, apiCopy.Req.Payload.Data, 1)
	require.Equal(t, "ev_1", apiCopy.Req.Payload.Data[0].ID)
	require.Equal(t, 1, apiCopy.Resp.Payload.Count)

	spec := api.ApiSpec()
	require.Equal(t, "ingestz-stream-relevance-fashion-event", spec.Id)
	require.Equal(t, "/api/ingest/dataset/{dataset}/stream/relevance-fashion-event", spec.Path)
}
