package dszapiactivesolutionslist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparserapi"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapiactivesolutionslist"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solution"
)

func TestUnitDszListActiveSolutionsApi(t *testing.T) {
	ctx := t.Context()
	api := &dszapiactivesolutionslist.Api{
		Req: &dszapiactivesolutionslist.DszListActiveSolutionsReq{
			Dataset: lakedsz.MustParse("myds"),
		},
		Resp: &dszapiactivesolutionslist.DszListActiveSolutionsResp{
			Dataset: lakedsz.MustParse("myds"),
			Solutions: []solution.Solution{
				solution.MustParse("relevance"),
				solution.MustParse("plat"),
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapiactivesolutionslist.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	parsed := jsparserapi.ParseApi(ctx, api)
	require.NotNil(t, parsed)
	require.Equal(t, "dsz-list-active-solutions", parsed.Spec.Id)
}
