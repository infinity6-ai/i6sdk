package cmdcmdz_test

import (
	"context"
	"testing"

	"github.com/infinity6-ai/i6sdk/cmdz/cmdcmdz"
	"github.com/infinity6-ai/i6sdk/versionz/version"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdCmdzPrepare(t *testing.T) {
	ctx := context.Background()
	cmd := cmdcmdz.Prepare(ctx)
	require.NotNil(t, cmd)
	require.Equal(t, "cmdz", cmd.Use)
	require.Equal(t, version.Version(), cmd.Version)

	subcommands := cmd.Commands()
	require.NotEmpty(t, subcommands)

	names := make(map[string]bool)
	for _, c := range subcommands {
		names[c.Name()] = true
	}
	require.True(t, names["version"])
	require.True(t, names["js"])
}
