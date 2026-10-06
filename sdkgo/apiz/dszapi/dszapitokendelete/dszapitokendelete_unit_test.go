package dszapitokendelete_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapitokendelete"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
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

	parsed := jsparser.ParseApi(api)
	require.NotNil(t, parsed)
	require.Equal(t, "dsz-delete-token", parsed.Spec.Id)
}
