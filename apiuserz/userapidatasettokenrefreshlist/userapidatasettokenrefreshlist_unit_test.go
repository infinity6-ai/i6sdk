package userapidatasettokenrefreshlist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefreshlist"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/stretchr/testify/require"
)

func TestUnitUserDatasetTokenRefreshListApi(t *testing.T) {
	api := &userapidatasettokenrefreshlist.Api{
		Req: &userapidatasettokenrefreshlist.Req{
			Params: &userapidatasettokenrefreshlist.ReqParams{
				Dataset: lakedsz.MustParse("myds"),
			},
			Headers: &userapidatasettokenrefreshlist.ReqHeaders{
				Authorization: "Bearer token",
			},
		},
		Resp: &userapidatasettokenrefreshlist.Resp{
			Payload: &userapidatasettokenrefreshlist.RespPayload{
				Tokens: []userapidatasettokenrefreshlist.RefreshTokenSummary{
					{
						Jti:     "jti_1",
						Dataset: lakedsz.MustParse("myds"),
					},
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapidatasettokenrefreshlist.Api{})
	require.Equal(t, api.Req.Params.Dataset, apiCopy.Req.Params.Dataset)

	spec := api.ApiSpec()
	require.Equal(t, "user-dataset-token-refresh-list", spec.Id)
	require.Equal(t, "List Dataset Refresh Tokens", spec.Desc().Name)
}
