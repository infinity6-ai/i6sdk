package dszapiingestiontokencreate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapiingestiontokencreate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/dsztoken"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
)

func TestUnitDszCreateIngestionTokenApi(t *testing.T) {
	api := &dszapiingestiontokencreate.Api{
		Req: dszapiingestiontokencreate.DszCreateIngestionTokenReq{
			Dataset:    lakedsz.MustParse("myds"),
			Annotation: "Annotation",
		},
		Resp: dsztoken.DatasetToken{
			Id:             "token1",
			Jti:            "jti1",
			Dataset:        lakedsz.MustParse("myds"),
			Annotation:     "Annotation",
			ExpirationDate: 12345678,
			Audience:       "aud",
			Grants:         []string{"write"},
			Token:          "tok-1",
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapiingestiontokencreate.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-create-ingestion-token", spec.Id)
}
