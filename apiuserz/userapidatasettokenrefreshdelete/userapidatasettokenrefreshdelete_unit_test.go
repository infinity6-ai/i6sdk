package userapidatasettokenrefreshdelete_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefreshdelete"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitUserDatasetTokenRefreshDeleteApi(t *testing.T) {
	api := &userapidatasettokenrefreshdelete.Api{
		Req: &userapidatasettokenrefreshdelete.Req{
			Params: &userapidatasettokenrefreshdelete.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
				Jti:     "jti_123",
			},
			Headers: &userapidatasettokenrefreshdelete.ReqHeaders{
				Authorization: "Bearer token",
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapidatasettokenrefreshdelete.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)
	require.Equal(t, api.Req.Params.Jti, apiCopy.Req.Params.Jti)

	spec := api.ApiSpec()
	require.Equal(t, "user-dataset-token-refresh-delete", spec.Id)
	require.Equal(t, "Revoke Dataset Refresh Token", spec.Desc().Name)
}
