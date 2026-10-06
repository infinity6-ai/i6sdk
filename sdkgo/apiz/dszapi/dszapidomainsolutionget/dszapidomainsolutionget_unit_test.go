package dszapidomainsolutionget_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapidomainsolutionget"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/domain"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solution"
)

func TestUnitDszSolutionDomainGetApi(t *testing.T) {
	api := &dszapidomainsolutionget.Api{
		Req: &dszapidomainsolutionget.DszSolutionDomainGetReq{
			Dataset: lakedsz.MustParse("myds"),
			Store:   "store1",
		},
		Resp: &dszapidomainsolutionget.DszSolutionDomainGetResp{
			Dataset: lakedsz.MustParse("myds"),
			Solutions: []dszapidomainsolutionget.SolutionDomain{
				{
					Solution: solution.MustParse("relevance"),
					Domain:   domain.MustParse("fashion"),
				},
			},
			Store: "store1",
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapidomainsolutionget.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	parsed := jsparser.ParseApi(api)
	require.NotNil(t, parsed)
	require.Equal(t, "dsz-domain-get", parsed.Spec.Id)
}
