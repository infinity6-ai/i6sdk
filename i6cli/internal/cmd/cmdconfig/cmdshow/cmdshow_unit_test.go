package cmdshow_test

import (
	"bytes"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdconfig/cmdshow"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdShow(t *testing.T) {
	ctx := t.Context()

	t.Run("Shows uninitialized profile details in table format", func(t *testing.T) {
		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		cmd := cmdshow.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", "main", "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		output := outBuf.String()
		require.Contains(t, output, "Profile")
		require.Contains(t, output, "main")
		require.Contains(t, output, "Userz Service")
		require.Contains(t, output, "not initialized")
	})

	t.Run("Shows initialized profile in json format", func(t *testing.T) {
		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		profile := "custom-worker"
		cfg := &cliconfig.Config{
			DefaultDataset: "ds-staging-1",
			UserzUrl:       "https://custom-userz.corp",
		}
		require.NoError(t, cliconfig.SaveConfig(profile, cfg, ctx))

		cmd := cmdshow.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "json", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "--format", "json"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		output := outBuf.String()
		require.Contains(t, output, `"profile": "custom-worker"`)
		require.Contains(t, output, `"default_dataset": "ds-staging-1"`)
		require.Contains(t, output, `"status": "initialized"`)
		require.Contains(t, output, `"userz": "https://custom-userz.corp"`)
	})

	t.Run("Shows profile in quiet format", func(t *testing.T) {
		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		cmd := cmdshow.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", "qa-profile", "Profile")
		cmd.PersistentFlags().String("format", "quiet", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", "qa-profile", "--format", "quiet"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)
		require.Equal(t, "qa-profile\n", outBuf.String())
	})
}
