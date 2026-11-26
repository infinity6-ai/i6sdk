package dszapisignedurl_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapisignedurl"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszSignedUrlApi(t *testing.T) {
	api := &dszapisignedurl.Api{
		Req: &dszapisignedurl.Req{
			Params: &dszapisignedurl.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
			Payload: &dszapisignedurl.ReqPayload{
				Paths: []*dszapisignedurl.Path{},
			},
		},
		Resp: &dszapisignedurl.Resp{
			Payload: &dszapisignedurl.RespPayload{
				Urls: []*dszapisignedurl.Url{
					{Url: "https://storage.googleapis.com/test"},
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapisignedurl.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-signed-url", spec.Id)
}
