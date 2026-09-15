package dszapi_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/dszapi"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

func TestUnitDszApi(t *testing.T) {
	productApi := &dszapi.DszProductSelectApi{
		Req: &dszapi.DszProductSelectReq{
			Dataset: lakedsz.MustParse("dsxxx"),
			Product: "product",
		},
		Resp: &dszapi.DszProductSelectResp{
			Dataset:  lakedsz.MustParse("dsxxx"),
			Products: []string{"product"},
		},
	}

	productApi2 := jsonz.MustCopy(productApi, &dszapi.DszProductSelectApi{})

	require.Equal(t, productApi, productApi2)
}
