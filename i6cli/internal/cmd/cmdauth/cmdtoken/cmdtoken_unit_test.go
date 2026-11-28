package cmdtoken_test

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdauth/cmdtoken"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdToken(t *testing.T) {
	ctx := t.Context()

	// Generate a valid 32-byte Ed25519 private key seed for offline robot key test
	var privKeyB64 string
	for {
		_, priv, err := ed25519.GenerateKey(nil)
		require.NoError(t, err)
		privKeyB64 = base64.RawURLEncoding.EncodeToString(priv.Seed())
		if !strings.HasPrefix(privKeyB64, "-") {
			break
		}
	}

	t.Run("Ingests Ed25519 robot refresh key offline", func(t *testing.T) {
		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		profile := "robot-worker"
		cmd := cmdtoken.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{
			"-p", profile,
			"--offline",
			"--dataset", "test-fashion",
			privKeyB64,
		})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)
		require.Contains(t, outBuf.String(), "test-fashion")

		// Verify key was saved in profile
		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)
		rkeys, err := cli.Tokens().GetDatasetRefreshTokensMap(ctx)
		require.NoError(t, err)
		require.Contains(t, rkeys, "test-fashion")
	})

	t.Run("Ingests token piped via stdin", func(t *testing.T) {
		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		profile := "stdin-profile"
		cmd := cmdtoken.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "json", "Format")

		inBuf := strings.NewReader(privKeyB64)
		outBuf := new(bytes.Buffer)
		cmd.SetIn(inBuf)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{
			"-p", profile,
			"--offline",
			"--dataset", "piped-ds",
			"-",
		})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, `"target": "piped-ds"`)
		require.Contains(t, out, `"type": "dataset_refresh"`)
	})

	t.Run("Fails when token input is empty", func(t *testing.T) {
		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		cmd := cmdtoken.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", "main", "Profile")
		cmd.PersistentFlags().String("format", "", "Format")
		cmd.SetIn(strings.NewReader("   "))
		cmd.SetArgs([]string{"-"})

		err = cmd.ExecuteContext(ctx)
		require.Error(t, err)
		require.Equal(t, exitcode.Usage, exitcode.Code(err))
		require.Contains(t, err.Error(), "token string cannot be empty")
	})
}
