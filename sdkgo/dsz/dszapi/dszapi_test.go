package dszapi_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/dszapi"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

func TestUnitDszApi(t *testing.T) {
	productApi := &dszapi.DszSolutionSelectApi{
		Req: &dszapi.DszSolutionSelectReq{
			Dataset:  lakedsz.MustParse("dsxxx"),
			Solution: "product",
		},
		Resp: &dszapi.DszSolutionSelectResp{
			Dataset:   lakedsz.MustParse("dsxxx"),
			Solutions: []string{"product"},
		},
	}

	productApi2 := jsonz.MustCopy(productApi, &dszapi.DszSolutionSelectApi{})

	require.Equal(t, productApi, productApi2)
}
