package dszapirecsystokenslist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapirecsystokenslist"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszListRecsysTokensApi(t *testing.T) {
	api := &dszapirecsystokenslist.Api{
		Req: &dszapirecsystokenslist.Req{
			Params: &dszapirecsystokenslist.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
		},
		Resp: &dszapirecsystokenslist.Resp{
			Payload: &dszapirecsystokenslist.RespPayload{},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapirecsystokenslist.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-list-recsys-tokens", spec.Id)
}
