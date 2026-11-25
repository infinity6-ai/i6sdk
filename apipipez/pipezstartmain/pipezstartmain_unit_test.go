package pipezstartmain_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apipipez/pipezstartmain"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitPipezStartMainApi(t *testing.T) {
	type testScenario struct {
		name    string
		dataset string
		runId   string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()

		api := &pipezstartmain.Api{
			Req: &pipezstartmain.Req{
				Params: &pipezstartmain.Params{
					Dataset: lakedsz.MustParse(s.dataset),
				},
			},
			Resp: &pipezstartmain.Resp{
				Payload: &pipezstartmain.RespPayload{
					Id: s.runId,
				},
			},
		}

		apiCopy := jsonz.MustCopy(api, &pipezstartmain.Api{})
		require.Equal(t, api.Req, apiCopy.Req)
		require.Equal(t, api.Resp, apiCopy.Resp)

		spec := api.ApiSpec()
		require.Equal(t, "pipez-start-relevance-fashion-relevance-fashion-main", spec.Id)
		require.Equal(t, "POST", spec.Method)
		require.Equal(t, "/api/pipe/dataset/{dataset}/start/relevance_fashion/relevance-fashion-main", spec.Path)
		require.NotNil(t, spec.Desc)
		desc := spec.Desc()
		require.Equal(t, "Start Main Pipeline", desc.Name)
		require.Equal(t, "Starts the main relevance fashion pipeline on a dataset and returns its id.", desc.Summary)
		require.Equal(t, "dataset", desc.Labels["token_types"])

		refs := api.ApiDataRefs()
		require.NotNil(t, refs)
		require.NotNil(t, refs.PathParams)
		require.Nil(t, refs.ReqBody)
		require.NotNil(t, refs.RespBody)

		pathParams := refs.PathParams.Object(true)
		require.Contains(t, pathParams, "dataset")
		require.Equal(t, "Dataset identifier.", pathParams["dataset"].Desc.Summary)

		respBody := refs.RespBody.Object(true)
		require.Contains(t, respBody, "id")

		require.Equal(t, s.dataset, api.Dataset().String())
		pipe := api.Pipe()
		require.NotNil(t, pipe)
		require.Equal(t, "relevance-fashion-main", pipe.Pipename())
		require.Equal(t, "relevance_fashion", pipe.SolutionDomain().String())

		api.SetRespId("new_run_id")
		require.Equal(t, "new_run_id", api.Resp.Payload.Id)
	}

	t.Run("Valid main pipeline start request and response", func(t *testing.T) {
		check(t, testScenario{
			name:    "Valid main pipeline start request and response",
			dataset: "my-dataset",
			runId:   "run_main_001",
		})
	})
}
