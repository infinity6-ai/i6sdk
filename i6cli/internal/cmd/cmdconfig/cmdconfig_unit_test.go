package cmdconfig_test

import (
	"bytes"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdconfig"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdConfig(t *testing.T) {
	ctx := t.Context()

	t.Run("Config command registers all subcommands", func(t *testing.T) {
		cmd := cmdconfig.New(ctx)
		require.NotNil(t, cmd)
		require.Equal(t, "config", cmd.Use)

		subcommands := make(map[string]bool)
		for _, sub := range cmd.Commands() {
			subcommands[sub.Name()] = true
		}

		require.True(t, subcommands["init"])
		require.True(t, subcommands["show"])
		require.True(t, subcommands["set"])
	})

	t.Run("Config help output displays available commands", func(t *testing.T) {
		cmd := cmdconfig.New(ctx)
		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"--help"})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "init")
		require.Contains(t, out, "show")
		require.Contains(t, out, "set")
	})
}
