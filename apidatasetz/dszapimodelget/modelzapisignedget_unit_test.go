package dszapimodelget_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapimodelget"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszModelGetApi(t *testing.T) {
	api := &dszapimodelget.Api{
		Req: &dszapimodelget.Req{
			Params: &dszapimodelget.ReqParams{
				Dataset:   lakedsz.MustParse("myds"),
				ModelName: "recsys",
			},
			Payload: &dszapimodelget.ReqPayload{
				ModelVersion: "v1",
			},
		},
		Resp: &dszapimodelget.Resp{
			Payload: &dszapimodelget.RespPayload{
				Urls: []*dszapimodelget.ModelFile{
					{
						Name: "weights.bin",
						Url:  "https://storage.googleapis.com/test",
					},
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapimodelget.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)
	require.Equal(t, api.Req.Params.ModelName, apiCopy.Req.Params.ModelName)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-model-get", spec.Id)
}
