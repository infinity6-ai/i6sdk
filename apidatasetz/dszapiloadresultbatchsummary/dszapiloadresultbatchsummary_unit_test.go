package dszapiloadresultbatchsummary_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapiloadresultbatchsummary"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbnamez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
	"github.com/infinity6-ai/i6sdk/sdkcorez/resultz/platformresultz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszResultLastBatchSummaryApi(t *testing.T) {
	store := storeid.MustParse("store-1")
	table := laketbnamez.MustParse("catalog")

	api := &dszapiloadresultbatchsummary.Api{
		Req: &dszapiloadresultbatchsummary.Req{
			Params: &dszapiloadresultbatchsummary.ReqParams{
				Dataset: lakedsz.MustParse("my-dataset"),
			},
			Query: &dszapiloadresultbatchsummary.ReqQuery{
				Table:   table,
				StoreId: store,
			},
		},
		Resp: &dszapiloadresultbatchsummary.Resp{
			Payload: &dszapiloadresultbatchsummary.RespPayload{
				Data: &platformresultz.BatchSummary{
					Req: &platformresultz.BatchSummaryReq{
						Table:   table,
						StoreId: store,
					},
					Resp: &platformresultz.BatchSummaryData{
						StoreId:  "store-1",
						RowCount: 100,
					},
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapiloadresultbatchsummary.Api{})
	require.Equal(t, api.Req.Params.Dataset.Get(), apiCopy.Req.Params.Dataset.Get())
	require.Equal(t, table, apiCopy.Req.Query.Table)
	require.Equal(t, store, apiCopy.Req.Query.StoreId)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-result-last-batch-summary", spec.Id)
	require.Equal(t, "/api/ds/dataset/{dataset}/results/last-batch-summary", spec.Path)
}
