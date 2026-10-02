package dszapicreatestore_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapicreatestore"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

func TestUnitDszCreateStoreApi(t *testing.T) {
	api := &dszapicreatestore.DszCreateStoreApi{
		Req: &dszapicreatestore.DszCreateStoreReq{
			Dataset: lakedsz.MustParse("myds"),
			Store: dszapicreatestore.DszStore{
				Id:   "store1",
				Name: "Store One",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapicreatestore.DszCreateStoreApi{})
	require.Equal(t, api.Req, apiCopy.Req)

	parsed := jsparser.ParseApi(api)
	require.NotNil(t, parsed)
	require.Equal(t, "dsz-create-store", parsed.Spec.Id)
}
