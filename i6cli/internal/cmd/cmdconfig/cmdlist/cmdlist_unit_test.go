package cmdlist_test

import (
	"bytes"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdconfig/cmdlist"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdList(t *testing.T) {
	t.Run("Lists profiles in table format with active marker", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())

		require.NoError(t, cliconfig.SaveConfig("main", &cliconfig.Config{DefaultDataset: "ds-main"}, ctx))
		require.NoError(t, cliconfig.SaveConfig("worker", &cliconfig.Config{DefaultDataset: "ds-worker"}, ctx))

		cmd := cmdlist.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", "main", "Profile")
		cmd.PersistentFlags().String("format", "table", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", "main"})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "PROFILE")
		require.Contains(t, out, "main")
		require.Contains(t, out, "worker")
		require.Contains(t, out, "ds-main")
		require.Contains(t, out, "ds-worker")
		require.Contains(t, out, "*")
	})

	t.Run("Lists profiles in quiet format", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())

		require.NoError(t, cliconfig.SaveConfig("alpha", &cliconfig.Config{}, ctx))
		require.NoError(t, cliconfig.SaveConfig("beta", &cliconfig.Config{}, ctx))

		cmd := cmdlist.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", "alpha", "Profile")
		cmd.PersistentFlags().String("format", "quiet", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"--format", "quiet"})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "alpha\n")
		require.Contains(t, out, "beta\n")
	})

	t.Run("Lists profiles in json format", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())

		require.NoError(t, cliconfig.SaveConfig("test-json", &cliconfig.Config{DefaultDataset: "json-ds"}, ctx))

		cmd := cmdlist.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", "test-json", "Profile")
		cmd.PersistentFlags().String("format", "json", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"--format", "json"})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, `"name": "test-json"`)
		require.Contains(t, out, `"active": true`)
		require.Contains(t, out, `"default_dataset": "json-ds"`)
	})
}
