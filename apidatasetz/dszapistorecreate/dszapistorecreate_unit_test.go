package dszapistorecreate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapistorecreate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszCreateStoreApi(t *testing.T) {
	api := &dszapistorecreate.Api{
		Req: &dszapistorecreate.Req{
			Params: &dszapistorecreate.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapistorecreate.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-create-store", spec.Id)
}
