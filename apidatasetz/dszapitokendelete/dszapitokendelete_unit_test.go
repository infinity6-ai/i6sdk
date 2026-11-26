package dszapitokendelete_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapitokendelete"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszDeleteTokenApi(t *testing.T) {
	api := &dszapitokendelete.Api{
		Req: &dszapitokendelete.Req{
			Params: &dszapitokendelete.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
			Query: &dszapitokendelete.ReqQuery{
				Token: "tok_123",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapitokendelete.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-delete-token", spec.Id)
}
