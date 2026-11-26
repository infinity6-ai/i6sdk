package dszapitableget_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapitableget"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbpathz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszTableGetApi(t *testing.T) {
	api := &dszapitableget.Api{
		Req: &dszapitableget.Req{
			Params: &dszapitableget.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
				Path:    laketbpathz.MustParse("i6tb://tables.items/part.json"),
			},
		},
		Resp: &dszapitableget.Resp{
			Headers: &dszapitableget.RespHeaders{
				Location: "https://storage.googleapis.com/test",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapitableget.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-table-get", spec.Id)
}
