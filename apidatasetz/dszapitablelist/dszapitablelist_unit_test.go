package dszapitablelist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapitablelist"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszTableListApi(t *testing.T) {
	api := &dszapitablelist.Api{
		Req: &dszapitablelist.Req{
			Params: &dszapitablelist.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
		},
		Resp: &dszapitablelist.Resp{
			Payload: &dszapitablelist.RespPayload{
				Files: []*dszapitablelist.File{},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapitablelist.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-table-list", spec.Id)
}
