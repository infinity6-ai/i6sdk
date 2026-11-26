package dszapistorelist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapistorelist"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszListStoresApi(t *testing.T) {
	api := &dszapistorelist.Api{
		Req: &dszapistorelist.Req{
			Params: &dszapistorelist.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
		},
		Resp: &dszapistorelist.Resp{
			Payload: &dszapistorelist.RespPayload{
				Dataset: lakedsz.MustParse("myds"),
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapistorelist.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-list-stores", spec.Id)
}
