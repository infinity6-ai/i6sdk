package userapiconfsettings_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiconfsettings"
	"github.com/stretchr/testify/require"
)

func TestUnitUserConfSettingsApi(t *testing.T) {
	api := &userapiconfsettings.Api{
		Req: &userapiconfsettings.Req{},
		Resp: &userapiconfsettings.Resp{
			Payload: &userapiconfsettings.RespPayload{
				RefreshTokenMaxIdleSeconds:   86400,
				UserAccessTokenTTLSeconds:    3600,
				DatasetAccessTokenTTLSeconds: 3600,
				DeviceCodeTTLSeconds:         600,
				SignedRequestMaxTTLSeconds:   300,
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapiconfsettings.Api{})
	require.Equal(t, api.Resp.Payload.RefreshTokenMaxIdleSeconds, apiCopy.Resp.Payload.RefreshTokenMaxIdleSeconds)

	spec := api.ApiSpec()
	require.Equal(t, "user-conf-settings", spec.Id)
	require.Equal(t, "Public Settings Discovery", spec.Desc().Name)
}
