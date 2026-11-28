package cmdset_test

import (
	"bytes"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdconfig/cmdset"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/sdkgo/config"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdSet(t *testing.T) {
	ctx := t.Context()

	type testScenario struct {
		args        []string
		wantErrCode int
		wantDataset string
		wantUserz   string
		wantDec     string
		wantTf      string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		// Seed initial configuration
		initialCfg := &cliconfig.Config{DefaultDataset: "orig-ds"}
		require.NoError(t, cliconfig.SaveConfig("main", initialCfg, ctx))

		cmd := cmdset.New(ctx)
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
		if s.wantUserz != "" {
			require.Equal(t, s.wantUserz, cfg.UserzUrl)
		}
		if s.wantDec != "" {
			require.Equal(t, s.wantDec, cfg.DecsuiteUrl)
		}
		if s.wantTf != "" {
			require.Equal(t, s.wantTf, cfg.TfinferUrl)
		}
	}

	t.Run("Updates box preset to staging setting all 3 URLs", func(t *testing.T) {
		check(t, testScenario{
			args:        []string{"box", "staging"},
			wantDataset: "orig-ds",
			wantUserz:   config.StagingUrlOptions.Userz,
			wantDec:     config.StagingUrlOptions.Decsuite,
			wantTf:      config.StagingUrlOptions.Tfinfer,
		})
	})

	t.Run("Updates default_dataset", func(t *testing.T) {
		check(t, testScenario{
			args:        []string{"default_dataset", "new-fashion-ds"},
			wantDataset: "new-fashion-ds",
		})
	})

	t.Run("Updates dataset using dataset alias", func(t *testing.T) {
		check(t, testScenario{
			args:        []string{"dataset", "aliased-ds"},
			wantDataset: "aliased-ds",
		})
	})

	t.Run("Updates userz_url", func(t *testing.T) {
		check(t, testScenario{
			args:        []string{"userz_url", "https://auth.internal.corp"},
			wantDataset: "orig-ds",
			wantUserz:   "https://auth.internal.corp",
		})
	})

	t.Run("Updates decsuite_url", func(t *testing.T) {
		check(t, testScenario{
			args:        []string{"decsuite_url", "https://lake.internal.corp"},
			wantDataset: "orig-ds",
			wantDec:     "https://lake.internal.corp",
		})
	})

	t.Run("Updates tfinfer_url", func(t *testing.T) {
		check(t, testScenario{
			args:        []string{"tfinfer_url", "https://infer.internal.corp"},
			wantDataset: "orig-ds",
			wantTf:      "https://infer.internal.corp",
		})
	})

	t.Run("Fails on invalid box value", func(t *testing.T) {
		check(t, testScenario{
			args:        []string{"box", "invalid-box-val"},
			wantErrCode: exitcode.Usage,
		})
	})

	t.Run("Fails on unknown key", func(t *testing.T) {
		check(t, testScenario{
			args:        []string{"unknown_field", "value"},
			wantErrCode: exitcode.Usage,
		})
	})

	t.Run("Fails when insufficient arguments provided", func(t *testing.T) {
		check(t, testScenario{
			args:        []string{"dataset"},
			wantErrCode: exitcode.General,
		})
	})
}
