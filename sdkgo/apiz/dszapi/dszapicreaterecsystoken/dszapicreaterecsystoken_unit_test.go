package dszapicreaterecsystoken_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapicreaterecsystoken"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/dsztoken"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

func TestUnitDszCreateRecsysTokenApi(t *testing.T) {
	api := &dszapicreaterecsystoken.DszCreateRecsysTokenApi{
		Req: dszapicreaterecsystoken.DszCreateRecsysTokenReq{
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

	apiCopy := jsonz.MustCopy(api, &dszapicreaterecsystoken.DszCreateRecsysTokenApi{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	parsed := jsparser.ParseApi(api)
	require.NotNil(t, parsed)
	require.Equal(t, "dsz-create-recsys-token", parsed.Spec.Id)
}
