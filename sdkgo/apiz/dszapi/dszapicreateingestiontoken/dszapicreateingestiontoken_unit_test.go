package dszapicreateingestiontoken_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapicreateingestiontoken"
)

func TestUnitDszCreateIngestionTokenApi(t *testing.T) {
	api := &dszapicreateingestiontoken.DszCreateIngestionTokenApi{
		Req: &dszapicreateingestiontoken.DszCreateIngestionTokenReq{
			Dataset:    "myds",
			Annotation: "Annotation",
			Channel:    "chan1",
		},
		Resp: &dszapicreateingestiontoken.DatasetToken{
			Id:             "token1",
			Jti:            "jti1",
			Dataset:        "myds",
			Annotation:     "Annotation",
			ExpirationDate: 12345678,
			Audience:       "aud",
			Grants:         []string{"write"},
			Channel:        "chan1",
			Token:          "tok-1",
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapicreateingestiontoken.DszCreateIngestionTokenApi{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	parsed := jsparser.ParseApi(api)
	require.NotNil(t, parsed)
	require.Equal(t, "dsz-create-ingestion-token", parsed.Spec.Id)
}
