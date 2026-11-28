package cmddatasetrefresh_test

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdtokens/cmddatasetrefresh"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokens"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdDatasetRefresh(t *testing.T) {
	var privKeyB64 string
	for {
		_, priv, err := ed25519.GenerateKey(nil)
		require.NoError(t, err)
		privKeyB64 = base64.RawURLEncoding.EncodeToString(priv.Seed())
		if !strings.HasPrefix(privKeyB64, "-") {
			break
		}
	}

	t.Run("Set saves robot key for dataset", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "robot-prof"

		cmd := cmddatasetrefresh.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "set", "--dataset", "fashion-ds", privKeyB64})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)
		require.Contains(t, outBuf.String(), "fashion-ds")

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)
		key, err := cli.Tokens().GetDatasetRefreshToken(ctx, "fashion-ds")
		require.NoError(t, err)
		require.Equal(t, privKeyB64, key)
	})

	t.Run("Set fails when user session is active in profile", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "human-prof"

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		err = cli.Tokens().SetUserToken(ctx, &tokens.UserToken{
			UserId:          "usr_active",
			UserAccessToken: "jwt.access",
		})
		require.NoError(t, err)

		cmd := cmddatasetrefresh.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.SetArgs([]string{"-p", profile, "set", "--dataset", "any-ds", privKeyB64})

		err = cmd.ExecuteContext(ctx)
		require.Error(t, err)
		require.Equal(t, exitcode.Usage, exitcode.Code(err))
		require.Contains(t, err.Error(), "active user session exists")
	})

	t.Run("Get displays metadata safely without revealing secret key", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "robot-prof"

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		err = cli.Tokens().SetDatasetRefreshToken(ctx, "safe-ds", privKeyB64)
		require.NoError(t, err)

		cmd := cmddatasetrefresh.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "table", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "get", "--dataset", "safe-ds"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "safe-ds")
		require.Contains(t, out, "configured")
		require.NotContains(t, out, privKeyB64)
	})

	t.Run("Export prints raw key string to stdout", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "export-prof"

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		err = cli.Tokens().SetDatasetRefreshToken(ctx, "export-ds", privKeyB64)
		require.NoError(t, err)

		cmd := cmddatasetrefresh.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "export", "--dataset", "export-ds"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)
		require.Equal(t, privKeyB64+"\n", outBuf.String())
	})

	t.Run("Delete evicts robot key from store", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "del-prof"

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		err = cli.Tokens().SetDatasetRefreshToken(ctx, "del-ds", privKeyB64)
		require.NoError(t, err)

		cmd := cmddatasetrefresh.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "delete", "--dataset", "del-ds"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)
		require.Contains(t, outBuf.String(), "del-ds")

		rkeys, err := cli.Tokens().GetDatasetRefreshTokensMap(ctx)
		require.NoError(t, err)
		require.Empty(t, rkeys)
	})
}
