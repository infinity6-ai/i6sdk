package userapidatasettokenrefreshlist_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/userapi/userapidatasettokenrefreshlist"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solution"
)

func TestUnitUserDatasetTokenRefreshListApi(t *testing.T) {
	api := &userapidatasettokenrefreshlist.Api{
		Req: userapidatasettokenrefreshlist.Req{
			Dataset:    lakedsz.MustParse("ds64l7qauv4l6512kf0c"),
			AuthHeader: "Bearer eyJ...",
		},
		Resp: userapidatasettokenrefreshlist.Resp{
			Tokens: []userapidatasettokenrefreshlist.RefreshTokenSummary{
				{
					Jti:        "jti-1",
					Dataset:    lakedsz.MustParse("ds64l7qauv4l6512kf0c"),
					Solution:   solution.MustParse("relevance"),
					Grants:     12,
					CreatedBy:  "u-1",
					CreatedAt:  1792900000,
					LastUsedAt: 1792900500,
					ExpiresAt:  1892900000,
					Revoked:    false,
					Annotation: "Airflow",
				},
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapidatasettokenrefreshlist.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	spec := api.ApiSpec()
	require.Equal(t, "user-dataset-token-refresh-list", spec.Id)
}
