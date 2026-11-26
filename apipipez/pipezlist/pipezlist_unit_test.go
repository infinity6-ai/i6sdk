package pipezlist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apipipez/pipezlist"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/domain"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solutiondomain"
	"github.com/infinity6-ai/i6sdk/sdkcorez/pipez"
	"github.com/stretchr/testify/require"
)

func TestUnitPipezListApi(t *testing.T) {
	api := &pipezlist.Api{
		Req: &pipezlist.Req{
			Params: &pipezlist.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
			Query: &pipezlist.ReqQuery{
				SolutionDomain: solutiondomain.MustParse(solutiondomain.Value{
					Solution: solution.Relevance().Solution,
					Domain:   domain.MustParse("fashion"),
				}),
				Name:   "mypipe",
				Status: "CREATED",
				Cursor: "mycursor",
				Limit:  10,
			},
		},
		Resp: &pipezlist.Resp{
			Payload: &pipezlist.RespPayload{
				Pipelines: []*pipez.Pipeline{
					{
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
				Cursor: "next-cursor",
			},
		},
	}

	api2 := jsonz.MustCopy(api, &pipezlist.Api{})

	require.Equal(t, api.Req.Params.Dataset, api2.Req.Params.Dataset)
	require.Equal(t, api.Req.Query.SolutionDomain, api2.Req.Query.SolutionDomain)
	require.Equal(t, api.Req.Query.Name, api2.Req.Query.Name)
	require.Equal(t, api.Req.Query.Status, api2.Req.Query.Status)
	require.Equal(t, api.Req.Query.Cursor, api2.Req.Query.Cursor)
	require.Equal(t, api.Req.Query.Limit, api2.Req.Query.Limit)
	require.Equal(t, len(api.Resp.Payload.Pipelines), len(api2.Resp.Payload.Pipelines))
	require.Equal(t, api.Resp.Payload.Pipelines[0].Id, api2.Resp.Payload.Pipelines[0].Id)
	require.Equal(t, api.Resp.Payload.Pipelines[0].States[0].Status, api2.Resp.Payload.Pipelines[0].States[0].Status)
	require.Equal(t, api.Resp.Payload.Cursor, api2.Resp.Payload.Cursor)

	spec := api.ApiSpec()
	require.Equal(t, "pipez-list", spec.Id)
	require.Equal(t, "List Pipelines", spec.Desc().Name)
}
