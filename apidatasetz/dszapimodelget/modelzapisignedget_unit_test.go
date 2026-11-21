package dszapimodelget_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapimodelget"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
)

func TestUnitModelzApiSignedGet(t *testing.T) {
	api := &dszapimodelget.Api{
		Req: dszapimodelget.Req{
			Dataset:      lakedsz.MustParse("myds"),
			ModelName:    "embedding_model",
			ModelVersion: "v1.0.0",
			ModelPartitions: map[string]string{
				"channel": "retail",
				"locale":  "en-US",
			},
		},
		Resp: dszapimodelget.Resp{
			Urls: []*dszapimodelget.ModelFile{
				{
					Name: "model.onnx",
					Url:  "https://example.com/models/embedding_model/model.onnx",
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapimodelget.Api{})

	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)
	require.Equal(t, "v1.0.0", apiCopy.Req.ModelVersion)
	require.Equal(t, "retail", apiCopy.Req.ModelPartitions["channel"])
	require.Len(t, apiCopy.Resp.Urls, 1)
}

func TestUnitModelzApiSignedGetApiSpec(t *testing.T) {
	api := &dszapimodelget.Api{}
	spec := api.ApiSpec()
	require.Equal(t, "dsz-model-get", spec.Id)
	require.Equal(t, "POST", spec.Method)
	require.Equal(t, "/api/model/dataset/{dataset}/signed-get/{model_name}", spec.Path)
}

func TestUnitModelzApiSignedApiDataRefs(t *testing.T) {
	api := &dszapimodelget.Api{}
	refs := api.ApiDataRefs()

	require.NotNil(t, refs.PathParams)
	require.Nil(t, refs.QueryParams)
	require.Nil(t, refs.ReqHeaders)
	require.NotNil(t, refs.ReqBody)
	require.Nil(t, refs.RespHeaders)
	require.NotNil(t, refs.RespBody)

	// Test ReqBody JSON unmarshaling
	reqBodyJSON := []byte(`{"model_version":"v2.1.0","model_partitions":{"arch":"arm64"}}`)
	err := refs.ReqBody.UnmarshalJSON(reqBodyJSON)
	require.NoError(t, err)
	require.Equal(t, "v2.1.0", api.Req.ModelVersion)
	require.Equal(t, "arm64", api.Req.ModelPartitions["arch"])

	// Test ReqBody JSON marshaling
	marshaledReq, err := refs.ReqBody.MarshalJSON()
	require.NoError(t, err)
	require.JSONEq(t, string(reqBodyJSON), string(marshaledReq))

	// Test RespBody JSON unmarshaling
	respBodyJSON := []byte(`{"urls":[{"name":"model.onnx","url":"https://example.com/model.onnx"}]}`)
	err = refs.RespBody.UnmarshalJSON(respBodyJSON)
	require.NoError(t, err)
	require.Len(t, api.Resp.Urls, 1)
	require.Equal(t, "model.onnx", api.Resp.Urls[0].Name)
	require.Equal(t, "https://example.com/model.onnx", api.Resp.Urls[0].Url)
}
