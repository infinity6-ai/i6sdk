package config_test

import (
	"testing"

	"github.com/infinity6-ai/i6sdk/sdkgo/config"
	"github.com/stretchr/testify/require"
)

func TestUnitConfigPresets(t *testing.T) {
	t.Run("DefaultUrlOptions contains production endpoints", func(t *testing.T) {
		require.Equal(t, "https://decsuite.app.infinity6.ai", config.DefaultUrlOptions.Decsuite)
		require.Equal(t, "https://tfinfer.app.infinity6.ai", config.DefaultUrlOptions.Tfinfer)
		require.Equal(t, "https://userz.app.infinity6.ai", config.DefaultUrlOptions.Userz)
	})

	t.Run("SandboxUrlOptions contains sandbox endpoints", func(t *testing.T) {
		require.Equal(t, "https://decsuite.sandbox.rs.infinity6.ai", config.SandboxUrlOptions.Decsuite)
		require.Equal(t, "https://tfinfer.sandbox.lb.infinity6.ai", config.SandboxUrlOptions.Tfinfer)
		require.Equal(t, "https://userz.sandbox.rs.infinity6.ai", config.SandboxUrlOptions.Userz)
	})

	t.Run("StagingUrlOptions contains staging endpoints", func(t *testing.T) {
		require.Equal(t, "https://decsuite.staging.rs.infinity6.ai", config.StagingUrlOptions.Decsuite)
		require.Equal(t, "https://tfinfer.staging.lb.infinity6.ai", config.StagingUrlOptions.Tfinfer)
		require.Equal(t, "https://userz.staging.rs.infinity6.ai", config.StagingUrlOptions.Userz)
	})
}
