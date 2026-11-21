package dszapistorecreate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapistorecreate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/dszstore"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
)

func TestUnitDszCreateStoreApi(t *testing.T) {
	api := &dszapistorecreate.Api{
		Req: &dszapistorecreate.DszCreateStoreReq{
			Dataset: lakedsz.MustParse("myds"),
			Store: dszstore.DszStore{
				Id:   storeid.MustParse("store1"),
				Name: "Store One",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &dszapistorecreate.Api{})
	require.Equal(t, api.Req, apiCopy.Req)

	spec := api.ApiSpec()
	require.Equal(t, "dsz-create-store", spec.Id)
}
