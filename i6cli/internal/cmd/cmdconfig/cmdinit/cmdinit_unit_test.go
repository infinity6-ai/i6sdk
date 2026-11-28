package cmdinit_test

import (
	"bytes"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdconfig/cmdinit"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/sdkgo/config"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdInit(t *testing.T) {
	ctx := t.Context()

	type testScenario struct {
		args        []string
		wantErrCode int
		wantDataset string
		wantUserz   string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		cmd := cmdinit.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", "main", "Profile")
		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs(s.args)

		execErr := cmd.ExecuteContext(ctx)
		if s.wantErrCode != 0 {
			require.Error(t, execErr)
			require.Equal(t, s.wantErrCode, exitcode.Code(execErr))
			return
		}
		require.NoError(t, execErr)

		cfg, err := cliconfig.ResolveConfig("main", ctx)
		require.NoError(t, err)
		require.Equal(t, s.wantDataset, cfg.DefaultDataset)
		require.Equal(t, s.wantUserz, cfg.UserzUrl)
	}

	t.Run("Initializes default prod configuration", func(t *testing.T) {
		check(t, testScenario{
			args:        []string{},
			wantDataset: "",
			wantUserz:   config.DefaultUrlOptions.Userz,
		})
	})

	t.Run("Initializes staging configuration with dataset and staging URLs", func(t *testing.T) {
		check(t, testScenario{
			args:        []string{"--box", "staging", "--dataset", "my-fashion-ds"},
			wantDataset: "my-fashion-ds",
			wantUserz:   config.StagingUrlOptions.Userz,
		})
	})

	t.Run("Initializes sandbox configuration", func(t *testing.T) {
		check(t, testScenario{
			args:        []string{"--box", "sandbox"},
			wantDataset: "",
			wantUserz:   config.SandboxUrlOptions.Userz,
		})
	})

	t.Run("Fails on invalid box preset", func(t *testing.T) {
		check(t, testScenario{
			args:        []string{"--box", "unknown-box"},
			wantErrCode: exitcode.Usage,
		})
	})

	t.Run("Fails when already exists without force flag", func(t *testing.T) {
		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		// First initialization succeeds
		cmd1 := cmdinit.New(ctx)
		cmd1.PersistentFlags().StringP("profile", "p", "main", "Profile")
		require.NoError(t, cmd1.ExecuteContext(ctx))

		// Second initialization fails without force
		cmd2 := cmdinit.New(ctx)
		cmd2.PersistentFlags().StringP("profile", "p", "main", "Profile")
		err = cmd2.ExecuteContext(ctx)
		require.Error(t, err)
		require.Equal(t, exitcode.Usage, exitcode.Code(err))
		require.Contains(t, err.Error(), "already exists")

		// Third initialization succeeds with force
		cmd3 := cmdinit.New(ctx)
		cmd3.PersistentFlags().StringP("profile", "p", "main", "Profile")
		cmd3.SetArgs([]string{"-f", "--dataset", "overwritten-ds"})
		require.NoError(t, cmd3.ExecuteContext(ctx))

		cfg, err := cliconfig.ResolveConfig("main", ctx)
		require.NoError(t, err)
		require.Equal(t, "overwritten-ds", cfg.DefaultDataset)
	})
}
