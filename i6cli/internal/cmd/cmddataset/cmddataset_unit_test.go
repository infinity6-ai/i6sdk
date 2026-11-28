package cmddataset_test

import (
	"bytes"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmddataset"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdDataset(t *testing.T) {
	ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())

	t.Run("Dataset command group has core subcommands", func(t *testing.T) {
		cmd := cmddataset.New(ctx)
		subcommands := make(map[string]bool)
		for _, sub := range cmd.Commands() {
			subcommands[sub.Name()] = true
		}
		require.True(t, subcommands["list"])
		require.True(t, subcommands["get"])
		require.True(t, subcommands["table"])
		require.True(t, subcommands["store"])
		require.True(t, subcommands["keys"])
	})

	t.Run("Dataset help output displays subcommands", func(t *testing.T) {
		cmd := cmddataset.New(ctx)
		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"--help"})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "list")
		require.Contains(t, out, "get")
		require.Contains(t, out, "table")
		require.Contains(t, out, "store")
		require.Contains(t, out, "keys")
	})
}
