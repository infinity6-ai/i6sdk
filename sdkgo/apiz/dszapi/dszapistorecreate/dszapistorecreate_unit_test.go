package dszapistorecreate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparserapi"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapistorecreate"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/dszstore"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

func TestUnitDszCreateStoreApi(t *testing.T) {
	api := &dszapistorecreate.Api{
		Req: &dszapistorecreate.DszCreateStoreReq{
			Dataset: lakedsz.MustParse("myds"),
			StoreId: dszstore.DszStore{
				Id:   "store1",
				Name: "Store One",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapistorecreate.Api{})
	require.Equal(t, api.Req, apiCopy.Req)

	parsed := jsparserapi.ParseApi(api)
	require.NotNil(t, parsed)
	require.Equal(t, "dsz-create-store", parsed.Spec.Id)
}
