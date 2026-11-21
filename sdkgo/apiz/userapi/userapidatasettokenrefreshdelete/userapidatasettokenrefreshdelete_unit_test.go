package userapidatasettokenrefreshdelete_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/userapi/userapidatasettokenrefreshdelete"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

func TestUnitUserDatasetTokenRefreshDeleteApi(t *testing.T) {
	api := &userapidatasettokenrefreshdelete.Api{
		Req: userapidatasettokenrefreshdelete.Req{
			Dataset:    lakedsz.MustParse("ds64l7qauv4l6512kf0c"),
			Jti:        "jti-1234",
			AuthHeader: "Bearer eyJ...",
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapidatasettokenrefreshdelete.Api{})
	require.Equal(t, api.Req, apiCopy.Req)

	spec := api.ApiSpec()
	require.Equal(t, "user-dataset-token-refresh-delete", spec.Id)
}
