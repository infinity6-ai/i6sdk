package dszapiingestiontokencreate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapiingestiontokencreate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszCreateIngestionTokenApi(t *testing.T) {
	api := &dszapiingestiontokencreate.Api{
		Req: &dszapiingestiontokencreate.Req{
			Params: &dszapiingestiontokencreate.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
			Payload: &dszapiingestiontokencreate.ReqPayload{
				Annotation: "test token",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapiingestiontokencreate.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-create-ingestion-token", spec.Id)
}
