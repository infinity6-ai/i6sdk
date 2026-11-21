package apicollectionz_test

import (
	"testing"

	"github.com/infinity6-ai/i6sdk/apicollectionz"
	"github.com/stretchr/testify/require"
)

func TestUnitApis(t *testing.T) {
	apis := apicollectionz.Apis()
	require.NotEmpty(t, apis)

	expectedIds := []string{
		"dsz-domain-get",
		"tfz-model-metadata",
		"user-auth-device-code",
		"user-auth-device-token",
		"user-auth-exchange",
		"user-auth-google",
		"user-auth-logout",
		"user-auth-me",
		"user-auth-password-reset",
		"user-auth-password-signin",
		"user-auth-password-signup",
		"user-auth-refresh",
		"user-conf-settings",
		"user-dataset-token-access-attenuate",
		"user-dataset-token-access-create",
		"user-dataset-token-refresh",
		"user-dataset-token-refresh-create",
		"user-dataset-token-refresh-delete",
		"user-dataset-token-refresh-list",
		"user-token-check",
	}

	type testScenario struct {
		apiId string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		api := apicollectionz.Api(s.apiId)
		require.NotNil(t, api)
		spec := api.ApiSpec()
		require.Equal(t, s.apiId, spec.Id)
		require.NotEmpty(t, spec.Method)
		require.NotEmpty(t, spec.Path)
	}

	for _, id := range expectedIds {
		t.Run(id, func(t *testing.T) {
			check(t, testScenario{apiId: id})
		})
	}
}

func TestUnitApiNotFound(t *testing.T) {
	require.Panics(t, func() {
		apicollectionz.Api("non-existent-api-id")
	})
}
