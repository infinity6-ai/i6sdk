package dszapiingestiontokenslist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapiingestiontokenslist"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/dsztoken"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
)

func TestUnitDszListIngestionTokensApi(t *testing.T) {
	api := &dszapiingestiontokenslist.Api{
		Req: dszapiingestiontokenslist.DszListIngestionTokensReq{
			Dataset: lakedsz.MustParse("myds"),
		},
		Resp: dszapiingestiontokenslist.DszListIngestionTokensResp{
			Tokens: []dsztoken.DatasetToken{
				{
					Id:             "token1",
					Jti:            "jti1",
					Dataset:        lakedsz.MustParse("myds"),
					Annotation:     "Annotation",
					ExpirationDate: 12345678,
					Audience:       "aud",
					Grants:         []string{"write"},
					Channel:        "chan1",
					Token:          "tok-1",
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapiingestiontokenslist.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-list-ingestion-tokens", spec.Id)
}
