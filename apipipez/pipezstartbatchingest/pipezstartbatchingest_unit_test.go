package pipezstartbatchingest_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apipipez/pipezstartbatchingest"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbnamez"
	"github.com/stretchr/testify/require"
)

func TestUnitPipezStartBatchIngestApi(t *testing.T) {
	type testScenario struct {
		name     string
		dataset  string
		ingestId string
		table    string
		runId    string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()

		api := &pipezstartbatchingest.Api{
			Req: &pipezstartbatchingest.Req{
				Params: &pipezstartbatchingest.Params{
					Dataset: lakedsz.MustParse(s.dataset),
				},
				Payload: &pipezstartbatchingest.Payload{
					IngestId: s.ingestId,
					Table:    laketbnamez.MustParse(s.table),
				},
			},
			Resp: &pipezstartbatchingest.Resp{
				Payload: &pipezstartbatchingest.RespPayload{
					Id: s.runId,
				},
			},
		}

		apiCopy := jsonz.MustCopy(api, &pipezstartbatchingest.Api{})
		require.Equal(t, api.Req, apiCopy.Req)
		require.Equal(t, api.Resp, apiCopy.Resp)

		spec := api.ApiSpec()
		require.Equal(t, "pipez-start-relevance_fashion-relevance-fashion-batch-ingest", spec.Id)
		require.Equal(t, "POST", spec.Method)
		require.Equal(t, "/api/pipe/dataset/{dataset}/start/relevance_fashion/relevance-fashion-batch-ingest", spec.Path)
		require.NotNil(t, spec.Desc)
		desc := spec.Desc()
		require.Equal(t, "Start Batch Ingestion Pipeline", desc.Name)
		require.Equal(t, "Starts a batch ingestion pipeline on a dataset with the given params and returns its id.", desc.Summary)
		require.Equal(t, "dataset", desc.Labels["token_types"])

		refs := api.ApiDataRefs()
		require.NotNil(t, refs)
		require.NotNil(t, refs.PathParams)
		require.NotNil(t, refs.ReqBody)
		require.NotNil(t, refs.RespBody)

		pathParams := refs.PathParams.Object(true)
		require.Contains(t, pathParams, "dataset")
		require.Equal(t, "Dataset identifier.", pathParams["dataset"].Desc.Summary)

		reqBody := refs.ReqBody.Object(true)
		require.Contains(t, reqBody, "ingest_id")
		require.Contains(t, reqBody, "table")

		respBody := refs.RespBody.Object(true)
		require.Contains(t, respBody, "id")
	}

	t.Run("Valid batch ingest pipeline start request and response", func(t *testing.T) {
		check(t, testScenario{
			name:     "Valid batch ingest pipeline start request and response",
			dataset:  "my-dataset",
			ingestId: "ing_123",
			table:    "catalog",
			runId:    "pipe_run_999",
		})
	})
}
