package dszapiingestiontokenslist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapiingestiontokenslist"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszListIngestionTokensApi(t *testing.T) {
	api := &dszapiingestiontokenslist.Api{
		Req: &dszapiingestiontokenslist.Req{
			Params: &dszapiingestiontokenslist.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
		},
		Resp: &dszapiingestiontokenslist.Resp{
			Payload: &dszapiingestiontokenslist.RespPayload{},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapiingestiontokenslist.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-list-ingestion-tokens", spec.Id)
}
