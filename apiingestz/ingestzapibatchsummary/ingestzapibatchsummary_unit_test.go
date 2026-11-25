package ingestzapibatchsummary_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapibatchsummary"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbnamez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
	"github.com/stretchr/testify/require"
)

func TestUnitIngestzBatchSummaryApi(t *testing.T) {
	api := &ingestzapibatchsummary.Api{
		Req: &ingestzapibatchsummary.Req{
			Params: &ingestzapibatchsummary.ReqParams{
				Dataset: lakedsz.MustParse("my-dataset"),
			},
			Query: &ingestzapibatchsummary.ReqQuery{
				Table:   laketbnamez.MustParse("my-table"),
				StoreId: storeid.MustParse("store-123"),
			},
		},
		Resp: &ingestzapibatchsummary.Resp{
			Payload: &ingestzapibatchsummary.RespPayload{
				StoreId:    storeid.MustParse("store-123"),
				RowCount:   5000,
				IngestedAt: "2026-10-16T12:00:00Z",
				ColumnSummaries: []ingestzapibatchsummary.ColumnSummary{
					{
						ColumnName:    "id",
						DistinctCount: 5000,
						NullCount:     0,
					},
					{
						ColumnName:    "category",
						DistinctCount: 12,
						NullCount:     15,
					},
				},
			},
		},
	}

	api2 := jsonz.MustCopy(api, &ingestzapibatchsummary.Api{})

	require.Equal(t, api.Req.Params.Dataset, api2.Req.Params.Dataset)
	require.Equal(t, api.Req.Query.Table, api2.Req.Query.Table)
	require.Equal(t, api.Req.Query.StoreId, api2.Req.Query.StoreId)

	require.NotNil(t, api2.Resp)
	require.Equal(t, api.Resp.Payload.StoreId, api2.Resp.Payload.StoreId)
	require.Equal(t, api.Resp.Payload.RowCount, api2.Resp.Payload.RowCount)
	require.Equal(t, api.Resp.Payload.IngestedAt, api2.Resp.Payload.IngestedAt)

	require.Equal(t, len(api.Resp.Payload.ColumnSummaries), len(api2.Resp.Payload.ColumnSummaries))
	require.Equal(t, api.Resp.Payload.ColumnSummaries[0].ColumnName, api2.Resp.Payload.ColumnSummaries[0].ColumnName)
	require.Equal(t, api.Resp.Payload.ColumnSummaries[0].DistinctCount, api2.Resp.Payload.ColumnSummaries[0].DistinctCount)
	require.Equal(t, api.Resp.Payload.ColumnSummaries[0].NullCount, api2.Resp.Payload.ColumnSummaries[0].NullCount)

	spec := api.ApiSpec()
	require.Equal(t, "ingestz-last-batch-summary", spec.Id)
	require.Equal(t, "Get Last Successful Ingestion Data Summary", spec.Desc().Name)
}
