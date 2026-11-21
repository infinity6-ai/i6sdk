package dszapitableget_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapitableget"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbpathz"
)

func TestUnitDszApiTableGet(t *testing.T) {
	api := &dszapitableget.Api{
		Req: dszapitableget.Req{
			Dataset: lakedsz.MustParse("myds"),
			Path:    laketbpathz.MustParse("i6tb://mydb.mytable/part-1"),
		},
		Resp: dszapitableget.Resp{
			Location: "https://example.com/part-1",
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapitableget.Api{})

	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)
	require.Equal(t, "https://example.com/part-1", apiCopy.Resp.Location)
}

func TestUnitDszApiTableGetParseApi(t *testing.T) {
	api := &dszapitableget.Api{}
	spec := api.ApiSpec()
	require.Equal(t, "dsz-table-get", spec.Id)
}
