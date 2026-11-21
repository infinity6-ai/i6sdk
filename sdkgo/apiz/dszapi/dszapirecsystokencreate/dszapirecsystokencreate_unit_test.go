package dszapirecsystokencreate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapirecsystokencreate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/dsztoken"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
)

func TestUnitDszCreateRecsysTokenApi(t *testing.T) {
	api := &dszapirecsystokencreate.Api{
		Req: dszapirecsystokencreate.DszCreateRecsysTokenReq{
			Dataset:    lakedsz.MustParse("myds"),
			Annotation: "Annotation",
			Channel:    "chan1",
		},
		Resp: dsztoken.DatasetToken{
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
	}

	apiCopy := jsonz.MustCopy(api, &dszapirecsystokencreate.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-create-recsys-token", spec.Id)
}
