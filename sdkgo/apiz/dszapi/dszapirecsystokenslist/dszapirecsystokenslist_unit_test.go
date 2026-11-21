package dszapirecsystokenslist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapirecsystokenslist"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/dsztoken"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
)

func TestUnitDszListRecsysTokensApi(t *testing.T) {
	api := &dszapirecsystokenslist.Api{
		Req: dszapirecsystokenslist.DszListRecsysTokensReq{
			Dataset: lakedsz.MustParse("myds"),
		},
		Resp: dszapirecsystokenslist.DszListRecsysTokensResp{
			Tokens: []dsztoken.DatasetToken{
				{
					Id:             "token1",
					Jti:            "jti1",
					Dataset:        lakedsz.MustParse("myds"),
					Annotation:     "Annotation",
					ExpirationDate: 12345678,
					Audience:       "aud",
					Grants:         []string{"read"},
					Channel:        "chan1",
					Token:          "tok-1",
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapirecsystokenslist.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-list-recsys-tokens", spec.Id)
}
