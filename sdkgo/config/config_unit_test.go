package config_test

import (
	"testing"

	"github.com/infinity6-ai/i6sdk/sdkgo/config"
	"github.com/stretchr/testify/require"
)

func TestUnitConfigPresets(t *testing.T) {
	t.Run("SandboxUrls contains sandbox endpoints", func(t *testing.T) {
		require.Equal(t, config.DefaultDecsuiteUrl, config.SandboxUrls.Decsuite)
		require.Equal(t, config.DefaultTfinferUrl, config.SandboxUrls.Tfinfer)
		require.Equal(t, config.DefaultUserzUrl, config.SandboxUrls.Userz)
	})

	t.Run("StagingUrls contains staging endpoints", func(t *testing.T) {
		require.Equal(t, config.StagingDecsuiteUrl, config.StagingUrls.Decsuite)
		require.Equal(t, config.StagingTfinferUrl, config.StagingUrls.Tfinfer)
		require.Equal(t, config.StagingUserzUrl, config.StagingUrls.Userz)
	})
}
