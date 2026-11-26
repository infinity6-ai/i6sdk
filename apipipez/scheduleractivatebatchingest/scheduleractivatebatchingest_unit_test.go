package scheduleractivatebatchingest_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apipipez/scheduleractivatebatchingest"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbnamez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/pipez/relevancefashionpipez"
	"github.com/stretchr/testify/require"
)

func TestUnitSchedulerActivateBatchIngestApi(t *testing.T) {
	api := &scheduleractivatebatchingest.Api{
		Req: &scheduleractivatebatchingest.Req{
			Params: &scheduleractivatebatchingest.ReqParams{
				Dataset: lakedsz.MustParse("my-dataset"),
			},
			Payload: &relevancefashionpipez.BatchIngest{
				IngestId: "ing_123",
				Table:    laketbnamez.MustParse("catalog"),
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &scheduleractivatebatchingest.Api{})
	require.Equal(t, api.Req.Params.Dataset.Get(), apiCopy.Req.Params.Dataset.Get())
	require.NotNil(t, apiCopy.Req.Payload)
	require.Equal(t, "ing_123", apiCopy.Req.Payload.IngestId)

	spec := api.ApiSpec()
	require.Equal(t, "pipez-scheduler-activate-relevance-fashion-relevance-fashion-batch-ingest", spec.Id)
	require.Equal(t, "/api/pipe/dataset/{dataset}/scheduler/activate/relevance-fashion/relevance-fashion-batch-ingest", spec.Path)
}
