package dszapi_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/domain"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solution"
)

func TestUnitDszApi(t *testing.T) {
	productApi := &dszapi.DszSolutionDomainSelectApi{
		Req: &dszapi.DszSolutionDomainSelectReq{
			Dataset:  lakedsz.MustParse("myds"),
			Solution: solution.MustParse("raas"),
			Domain:   domain.MustParse("fashion"),
			Store:    "store1",
		},
		Resp: &dszapi.DszSolutionDomainSelectReq{
			Dataset:  lakedsz.MustParse("myds"),
			Solution: solution.MustParse("raas"),
			Domain:   domain.MustParse("fashion"),
			Store:    "store1",
		},
	}

	productApi2 := jsonz.MustCopy(productApi, &dszapi.DszSolutionDomainSelectApi{})

	require.Equal(t, *productApi.Req, *productApi2.Req)
}
