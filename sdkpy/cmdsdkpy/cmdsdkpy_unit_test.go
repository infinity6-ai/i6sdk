package cmdsdkpy_test

import (
	"context"
	"testing"

	"github.com/infinity6-ai/i6sdk/sdkpy/cmdsdkpy"
	"github.com/infinity6-ai/i6sdk/sdkpy/version"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdSdkPyPrepare(t *testing.T) {
	ctx := context.Background()
	cmd := cmdsdkpy.Prepare(ctx)
	require.NotNil(t, cmd)
	require.Equal(t, "sdkpy", cmd.Use)
	require.Equal(t, version.Version(), cmd.Version)

	subcommands := cmd.Commands()
	require.NotEmpty(t, subcommands)

	names := make(map[string]bool)
	for _, c := range subcommands {
		names[c.Name()] = true
	}
	require.True(t, names["version"])
	require.True(t, names["generate"])
}
