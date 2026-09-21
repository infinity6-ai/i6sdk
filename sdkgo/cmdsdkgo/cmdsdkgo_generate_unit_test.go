//go:build !testcover

package cmdsdkgo_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/cmdsdkgo"
)

// runCmd executa a árvore de comandos com os args dados e devolve a saída e o erro.
func runCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	rootCmd := cmdsdkgo.Prepare(t.Context())
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return out.String(), err
}

func TestUnitGenerateJs(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "apigenz")

	out, err := runCmd(t, "generate", "js", "--out", dir)
	require.NoError(t, err)

	for _, name := range []string{"dsz-domain-select.js", "ingestz-get-url.js", "pipez-start.js"} {
		require.FileExists(t, filepath.Join(dir, name))
		require.Contains(t, out, filepath.Join(dir, name))
	}
	require.Contains(t, out, "generated 3 apis into "+dir)

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.Len(t, entries, 3, "exactly one file per api, nothing else")
}

func TestUnitGenerateJsFailure(t *testing.T) {
	type testScenario struct {
		args   func(t *testing.T) []string
		errMsg string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		_, err := runCmd(t, s.args(t)...)
		require.Error(t, err)
		require.Contains(t, err.Error(), s.errMsg)
	}

	t.Run("unknown generator", func(t *testing.T) {
		check(t, testScenario{
			args:   func(t *testing.T) []string { return []string{"generate", "nope"} },
			errMsg: "unknown command",
		})
	})

	t.Run("out is below a regular file", func(t *testing.T) {
		check(t, testScenario{
			args: func(t *testing.T) []string {
				file := filepath.Join(t.TempDir(), "file")
				require.NoError(t, os.WriteFile(file, nil, 0o644))
				return []string{"generate", "js", "--out", filepath.Join(file, "sub")}
			},
			errMsg: "failed to generate js clients",
		})
	})
}
