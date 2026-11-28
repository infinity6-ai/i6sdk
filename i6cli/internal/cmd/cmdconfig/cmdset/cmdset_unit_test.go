package cmdset_test

import (
	"bytes"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdconfig/cmdset"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdSet(t *testing.T) {
	ctx := t.Context()

	type testScenario struct {
		args        []string
		wantErrCode int
		wantBox     string
		wantDataset string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		// Seed initial configuration
		initialCfg := &cliconfig.Config{Box: cliconfig.BoxProd, DefaultDataset: "orig-ds"}
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
		require.Equal(t, s.wantBox, cfg.Box)
		require.Equal(t, s.wantDataset, cfg.DefaultDataset)
	}

	t.Run("Updates box environment to staging", func(t *testing.T) {
		check(t, testScenario{
			args:        []string{"box", "staging"},
			wantBox:     cliconfig.BoxStaging,
			wantDataset: "orig-ds",
		})
	})

	t.Run("Updates default_dataset", func(t *testing.T) {
		check(t, testScenario{
			args:        []string{"default_dataset", "new-fashion-ds"},
			wantBox:     cliconfig.BoxProd,
			wantDataset: "new-fashion-ds",
		})
	})

	t.Run("Updates dataset using dataset alias", func(t *testing.T) {
		check(t, testScenario{
			args:        []string{"dataset", "aliased-ds"},
			wantBox:     cliconfig.BoxProd,
			wantDataset: "aliased-ds",
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
			args:        []string{"box"},
			wantErrCode: exitcode.General, // cobra flag/arg validation returns generic error
		})
	})
}
