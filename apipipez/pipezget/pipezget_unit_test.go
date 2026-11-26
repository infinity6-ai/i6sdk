package pipezget_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apipipez/pipezget"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solutiondomain"
	"github.com/infinity6-ai/i6sdk/sdkcorez/pipez"
	"github.com/stretchr/testify/require"
)

func TestUnitPipezGetApi(t *testing.T) {
	api := &pipezget.Api{
		Req: &pipezget.Req{
			Params: &pipezget.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
				Id:      "pipe-123",
			},
		},
		Resp: &pipezget.Resp{
			Payload: &pipezget.RespPayload{
				Pipeline: &pipez.Pipeline{
					Id:             "pipe-123",
					ConcurrentId:   "concurrent-123",
					SolutionDomain: solutiondomain.MustParseString("relevance_fashion"),
					Name:           "mypipe",
					Dataset:        "myds",
					States: []*pipez.PipeState{
						{
							Status:    "CREATED",
							Params:    map[string]any{"foo": "bar"},
							CreatedAt: "2026-10-13T12:00:00Z",
						},
					},
				},
			},
		},
	}

	api2 := jsonz.MustCopy(api, &pipezget.Api{})

	require.Equal(t, api.Req.Params.Dataset, api2.Req.Params.Dataset)
	require.Equal(t, api.Req.Params.Id, api2.Req.Params.Id)
	require.NotNil(t, api2.Resp.Payload.Pipeline)
	require.Equal(t, api.Resp.Payload.Pipeline.Id, api2.Resp.Payload.Pipeline.Id)
	require.Equal(t, api.Resp.Payload.Pipeline.ConcurrentId, api2.Resp.Payload.Pipeline.ConcurrentId)
	require.Equal(t, api.Resp.Payload.Pipeline.SolutionDomain, api2.Resp.Payload.Pipeline.SolutionDomain)
	require.Equal(t, api.Resp.Payload.Pipeline.Name, api2.Resp.Payload.Pipeline.Name)
	require.Equal(t, api.Resp.Payload.Pipeline.Dataset, api2.Resp.Payload.Pipeline.Dataset)
	require.Equal(t, len(api.Resp.Payload.Pipeline.States), len(api2.Resp.Payload.Pipeline.States))
	require.Equal(t, api.Resp.Payload.Pipeline.States[0].Status, api2.Resp.Payload.Pipeline.States[0].Status)

	spec := api.ApiSpec()
	require.Equal(t, "pipez-get", spec.Id)
	require.Equal(t, "Get Pipeline", spec.Desc().Name)
}
