package tfzapimodelpredict_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apitfz/tfzapimodelpredict"
	"github.com/stretchr/testify/require"
)

func TestUnitModelPredictApi(t *testing.T) {
	api := &tfzapimodelpredict.Api{
		Req: &tfzapimodelpredict.Req{
			Params: &tfzapimodelpredict.ReqParams{
				Dataset:   "myds",
				ModelName: "embedding_model",
				Version:   "v1.0.0",
			},
			Payload: &tfzapimodelpredict.ReqPayload{
				ModelPartitions: map[string]string{
					"store_id": "store1",
				},
				Inputs: `[["a","b"]]`,
			},
		},
		Resp: &tfzapimodelpredict.Resp{
			Payload: &tfzapimodelpredict.RespPayload{
				Outputs: `{"output":{"c","d"}}`,
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &tfzapimodelpredict.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)
	require.Equal(t, "myds", apiCopy.Req.Params.Dataset)
	require.Equal(t, "embedding_model", apiCopy.Req.Params.ModelName)
	require.Equal(t, "v1.0.0", apiCopy.Req.Params.Version)
	require.Equal(t, "store1", apiCopy.Req.Payload.ModelPartitions["store_id"])
	require.Equal(t, `[["a","b"]]`, apiCopy.Req.Payload.Inputs)
	require.Equal(t, `{"output":{"c","d"}}`, apiCopy.Resp.Payload.Outputs)

	spec := api.ApiSpec()
	require.Equal(t, "tfz-model-predict", spec.Id)
	require.Equal(t, "Model Predict", spec.Desc().Name)
}
