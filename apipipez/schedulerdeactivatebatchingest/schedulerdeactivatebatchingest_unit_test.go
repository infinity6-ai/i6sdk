package schedulerdeactivatebatchingest_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apipipez/schedulerdeactivatebatchingest"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbnamez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/pipez/relevancefashionpipez"
	"github.com/stretchr/testify/require"
)

func TestUnitSchedulerDeactivateBatchIngestApi(t *testing.T) {
	api := &schedulerdeactivatebatchingest.Api{
		Req: &schedulerdeactivatebatchingest.Req{
			Params: &schedulerdeactivatebatchingest.ReqParams{
				Dataset: lakedsz.MustParse("my-dataset"),
			},
			Payload: &relevancefashionpipez.BatchIngest{
				IngestId: "ing_123",
				Table:    laketbnamez.MustParse("catalog"),
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &schedulerdeactivatebatchingest.Api{})
	require.Equal(t, api.Req.Params.Dataset.Get(), apiCopy.Req.Params.Dataset.Get())
	require.NotNil(t, apiCopy.Req.Payload)
	require.Equal(t, "ing_123", apiCopy.Req.Payload.IngestId)

	spec := api.ApiSpec()
	require.Equal(t, "pipez-scheduler-deactivate-relevance-fashion-relevance-fashion-batch-ingest", spec.Id)
	require.Equal(t, "/api/pipe/dataset/{dataset}/scheduler/deactivate/relevance-fashion/relevance-fashion-batch-ingest", spec.Path)
}
