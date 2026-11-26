package dszapirecsystokencreate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapirecsystokencreate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszCreateRecsysTokenApi(t *testing.T) {
	api := &dszapirecsystokencreate.Api{
		Req: &dszapirecsystokencreate.Req{
			Params: &dszapirecsystokencreate.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
			Payload: &dszapirecsystokencreate.ReqPayload{
				Annotation: "recsys",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapirecsystokencreate.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-create-recsys-token", spec.Id)
}
