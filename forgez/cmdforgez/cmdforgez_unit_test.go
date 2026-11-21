package cmdforgez_test

import (
	"context"
	"testing"

	"github.com/infinity6-ai/i6sdk/forgez/cmdforgez"
	"github.com/infinity6-ai/i6sdk/versionz/version"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdForgezPrepare(t *testing.T) {
	ctx := context.Background()
	cmd := cmdforgez.Prepare(ctx)
	require.NotNil(t, cmd)
	require.Equal(t, "forgez", cmd.Use)
	require.Equal(t, version.Version(), cmd.Version)

	subcommands := cmd.Commands()
	require.NotEmpty(t, subcommands)

	names := make(map[string]bool)
	for _, c := range subcommands {
		names[c.Name()] = true
	}
	require.True(t, names["version"])
}
