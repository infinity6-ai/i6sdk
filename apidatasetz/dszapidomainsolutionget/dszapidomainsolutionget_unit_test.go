package dszapidomainsolutionget_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapidomainsolutionget"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
	"github.com/stretchr/testify/require"
)

func TestUnitDszDomainGetApi(t *testing.T) {
	api := &dszapidomainsolutionget.Api{
		Req: &dszapidomainsolutionget.Req{
			Params: &dszapidomainsolutionget.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
			Query: &dszapidomainsolutionget.ReqQuery{
				StoreId: storeid.MustParse("store-1"),
			},
		},
		Resp: &dszapidomainsolutionget.Resp{
			Payload: &dszapidomainsolutionget.RespPayload{
				Dataset: lakedsz.MustParse("myds"),
				StoreId: storeid.MustParse("store-1"),
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapidomainsolutionget.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-domain-get", spec.Id)
}
