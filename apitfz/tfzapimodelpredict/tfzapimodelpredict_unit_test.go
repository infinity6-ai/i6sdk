package tfzapimodelpredict_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"github.com/infinity6-ai/i6sdk/apitfz/tfzapimodelpredict"
)

func TestUnitModelPredictApi(t *testing.T) {
	api := &tfzapimodelpredict.Api{
		Req: &tfzapimodelpredict.Req{
			Dataset:   "myds",
			ModelName: "embedding_model",
			Version:   "v1.0.0",
			ModelPartitions: map[string]string{
				"store_id": "store1",
			},
			// Inputs: []any{
			// 	[]any{"a", "b"},
			// },
			Inputs: `[["a","b"]]`,
		},
		Resp: &tfzapimodelpredict.Resp{
			// Outputs: map[string]any{
			// 	"output": []any{"c", "d"},
			// },
			Outputs: `{"output":{"c","d"}}`,
		},
	}

	apiCopy := jsonz.MustCopy(api, &tfzapimodelpredict.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)
	require.Equal(t, "myds", apiCopy.Req.Dataset)
	require.Equal(t, "embedding_model", apiCopy.Req.ModelName)
	require.Equal(t, "v1.0.0", apiCopy.Req.Version)
	require.Equal(t, "store1", apiCopy.Req.ModelPartitions["store_id"])
	require.Equal(t, `[["a","b"]]`, apiCopy.Req.Inputs)
	require.Equal(t, `{"output":{"c","d"}}`, apiCopy.Resp.Outputs)

	spec := api.ApiSpec()
	require.Equal(t, "tfz-model-predict", spec.Id)
}
