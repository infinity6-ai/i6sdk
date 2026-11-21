package userapiconfsettings_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/userapi/userapiconfsettings"
)

func TestUnitUserConfSettingsApi(t *testing.T) {
	api := &userapiconfsettings.Api{
		Req: userapiconfsettings.Req{},
		Resp: userapiconfsettings.Resp{
			RefreshTokenMaxIdleSeconds:  7776000,
			UserAccessTokenTTLSeconds:    3600,
			DatasetAccessTokenTTLSeconds: 3600,
			DeviceCodeTTLSeconds:        900,
			SignedRequestMaxTTLSeconds:  86400,
		},
	}

	spec := api.ApiSpec()
	require.Equal(t, "user-conf-settings", spec.Id)
	require.Equal(t, "GET", spec.Method)
	require.Equal(t, "/api/userz/conf/settings.json", spec.Path)
	require.Equal(t, "userz", spec.Desc().Labels["servname"])

	refs := api.ApiDataRefs()
	require.NotNil(t, refs)
	require.Nil(t, refs.ReqBody)
	require.NotNil(t, refs.RespBody)

	apiCopy := jsonz.MustCopy(api, &userapiconfsettings.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)
}
