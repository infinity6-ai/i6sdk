package cmdcli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/cmdcli"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/versionz/version"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdcli(t *testing.T) {
	t.Run("Prepare constructs root command with version and persistent flags", func(t *testing.T) {
		cmd := cmdcli.Prepare(t.Context())
		require.NotNil(t, cmd)
		require.Equal(t, "i6", cmd.Use)
		require.NotEmpty(t, cmd.Version)

		// Persistent flags
		profileFlag := cmd.PersistentFlags().Lookup("profile")
		require.NotNil(t, profileFlag)
		require.Equal(t, "p", profileFlag.Shorthand)
		require.Equal(t, "main", profileFlag.DefValue)

		formatFlag := cmd.PersistentFlags().Lookup("format")
		require.NotNil(t, formatFlag)
		require.Empty(t, formatFlag.DefValue)
	})

	t.Run("Version command prints version to stdout", func(t *testing.T) {
		cmd := cmdcli.Prepare(t.Context())
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetArgs([]string{"version"})

		code := cmdcli.ExecuteWithCode(t.Context(), cmd)
		require.Equal(t, exitcode.Success, code)
		require.Equal(t, version.Version()+"\n", buf.String())
	})

	t.Run("ExecuteWithCode handles errors and extracts exit code", func(t *testing.T) {
		cmd := cmdcli.Prepare(t.Context())
		cmd.AddCommand(&cobra.Command{
			Use: "fail-usage",
			RunE: func(cmd *cobra.Command, args []string) error {
				return exitcode.New(exitcode.Usage, "simulated usage error")
			},
		})

		errBuf := new(bytes.Buffer)
		cmd.SetErr(errBuf)
		cmd.SetArgs([]string{"fail-usage"})

		code := cmdcli.ExecuteWithCode(t.Context(), cmd)
		require.Equal(t, exitcode.Usage, code)
		require.Contains(t, errBuf.String(), "simulated usage error")
	})
}

func TestUnitCompletion(t *testing.T) {
	type testScenario struct {
		shell       string
		wantErrCode int
		wantKeyword string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		cmd := cmdcli.Prepare(t.Context())
		outBuf := new(bytes.Buffer)
		errBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetErr(errBuf)
		cmd.SetArgs([]string{"completion", s.shell})

		code := cmdcli.ExecuteWithCode(t.Context(), cmd)
		if s.wantErrCode != 0 {
			require.Equal(t, s.wantErrCode, code)
			return
		}
		require.Equal(t, exitcode.Success, code)
		require.Contains(t, strings.ToLower(outBuf.String()), strings.ToLower(s.wantKeyword))
	}

	t.Run("Bash completion generates valid script", func(t *testing.T) {
		check(t, testScenario{
			shell:       "bash",
			wantKeyword: "bash completion",
		})
	})

	t.Run("Zsh completion generates valid script", func(t *testing.T) {
		check(t, testScenario{
			shell:       "zsh",
			wantKeyword: "compdef",
		})
	})

	t.Run("Fish completion generates valid script", func(t *testing.T) {
		check(t, testScenario{
			shell:       "fish",
			wantKeyword: "complete",
		})
	})

	t.Run("PowerShell completion generates valid script", func(t *testing.T) {
		check(t, testScenario{
			shell:       "powershell",
			wantKeyword: "Register-ArgumentCompleter",
		})
	})

	t.Run("Unsupported shell returns Usage error code", func(t *testing.T) {
		check(t, testScenario{
			shell:       "unsupported-shell",
			wantErrCode: exitcode.Usage,
		})
	})
}
