package tfzapimodelmetadata_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apitfz/tfzapimodelmetadata"
	"github.com/stretchr/testify/require"
)

func TestUnitModelMetadataApi(t *testing.T) {
	api := &tfzapimodelmetadata.Api{
		Req: &tfzapimodelmetadata.Req{
			Params: &tfzapimodelmetadata.ReqParams{
				ModelName: "embedding_model",
				Version:   "v1.0.0",
			},
		},
		Resp: &tfzapimodelmetadata.Resp{
			Payload: &tfzapimodelmetadata.RespPayload{
				Inputs: []tfzapimodelmetadata.TensorDetail{
					{
						Name:  "input_1",
						Index: 0,
						Shape: []int{1, 256},
						DType: "float32",
					},
				},
				Outputs: []tfzapimodelmetadata.TensorDetail{
					{
						Name:  "output_1",
						Index: 0,
						Shape: []int{1, 128},
						DType: "float32",
					},
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &tfzapimodelmetadata.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)
	require.Equal(t, "embedding_model", apiCopy.Req.Params.ModelName)
	require.Equal(t, "v1.0.0", apiCopy.Req.Params.Version)
	require.Len(t, apiCopy.Resp.Payload.Inputs, 1)
	require.Equal(t, "input_1", apiCopy.Resp.Payload.Inputs[0].Name)
	require.Equal(t, []int{1, 256}, apiCopy.Resp.Payload.Inputs[0].Shape)

	spec := api.ApiSpec()
	require.Equal(t, "tfz-model-metadata", spec.Id)
	require.Equal(t, "Get Model Metadata", spec.Desc().Name)
}
