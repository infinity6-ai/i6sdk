package config_test

import (
	"testing"

	"github.com/infinity6-ai/i6sdk/sdkgo/config"
	"github.com/stretchr/testify/require"
)

func TestUnitConfigPresets(t *testing.T) {
	t.Run("ProdUrls contains production endpoints and is the default", func(t *testing.T) {
		require.Equal(t, config.ProdDecsuiteUrl, config.ProdUrls.Decsuite)
		require.Equal(t, config.ProdTfinferUrl, config.ProdUrls.Tfinfer)
		require.Equal(t, config.ProdUserzUrl, config.ProdUrls.Userz)

		require.Equal(t, config.ProdDecsuiteUrl, config.DefaultDecsuiteUrl)
		require.Equal(t, config.ProdTfinferUrl, config.DefaultTfinferUrl)
		require.Equal(t, config.ProdUserzUrl, config.DefaultUserzUrl)
	})

	t.Run("SandboxUrls contains sandbox endpoints", func(t *testing.T) {
		require.Equal(t, config.SandboxDecsuiteUrl, config.SandboxUrls.Decsuite)
		require.Equal(t, config.SandboxTfinferUrl, config.SandboxUrls.Tfinfer)
		require.Equal(t, config.SandboxUserzUrl, config.SandboxUrls.Userz)
	})

	t.Run("StagingUrls contains staging endpoints", func(t *testing.T) {
		require.Equal(t, config.StagingDecsuiteUrl, config.StagingUrls.Decsuite)
		require.Equal(t, config.StagingTfinferUrl, config.StagingUrls.Tfinfer)
		require.Equal(t, config.StagingUserzUrl, config.StagingUrls.Userz)
	})
}
