package dszapitokendelete_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapitokendelete"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
)

func TestUnitDszDeleteTokenApi(t *testing.T) {
	api := &dszapitokendelete.Api{
		Req: dszapitokendelete.DszDeleteTokenReq{
			Dataset: lakedsz.MustParse("myds"),
			Token:   "tok-123",
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapitokendelete.Api{})
	require.Equal(t, api.Req, apiCopy.Req)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-delete-token", spec.Id)
}
