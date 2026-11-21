package dszapitablelist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/gox/fsz/fsz/fszsignmethod"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapitablelist"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbpathz"
	"github.com/stretchr/testify/require"
)

func TestUnitDszApiTableList(t *testing.T) {
	api := &dszapitablelist.Api{
		Req: dszapitablelist.Req{
			Dataset:    lakedsz.MustParse("myds"),
			Path:       laketbpathz.MustParse("i6tb://mydb.mytable"),
			SignMethod: fszsignmethod.SignMethodGet,
			Recursive:  true,
			Stat:       true,
			FromCursor: "cur-123",
			Limit:      50,
		},
		Resp: dszapitablelist.Resp{
			Files: []*dszapitablelist.File{
				{
					Path:      laketbpathz.MustParse("i6tb://mydb.mytable/part-1"),
					SignedUrl: "https://example.com/part-1",
				},
			},
			NextCursor: "cur-456",
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapitablelist.Api{})

	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)
	require.Equal(t, "cur-123", apiCopy.Req.FromCursor)
	require.Equal(t, "cur-456", apiCopy.Resp.NextCursor)
}

func TestUnitDszApiTableListParseApi(t *testing.T) {
	api := &dszapitablelist.Api{}
	spec := api.ApiSpec()
	require.Equal(t, "dsz-table-list", spec.Id)
}
