package cmdsdkgo_test

import (
	"context"
	"testing"

	"github.com/infinity6-ai/i6sdk/sdkgo/cmdsdkgo"
	"github.com/infinity6-ai/i6sdk/versionz/version"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdSdkGoPrepare(t *testing.T) {
	ctx := context.Background()

	type testScenario struct {
		name         string
		expectedName string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		cmd := cmdsdkgo.Prepare(ctx)
		require.NotNil(t, cmd)
		require.Equal(t, s.expectedName, cmd.Use)
		require.Equal(t, version.Version(), cmd.Version)
		require.NotEmpty(t, cmd.Commands())

		hasVersion := false
		for _, c := range cmd.Commands() {
			if c.Name() == "version" {
				hasVersion = true
			}
		}
		require.True(t, hasVersion)
	}

	t.Run("Prepare returns root command with version subcommand", func(t *testing.T) {
		check(t, testScenario{
			name:         "Prepare returns root command with version subcommand",
			expectedName: "sdkgo",
		})
	})
}
