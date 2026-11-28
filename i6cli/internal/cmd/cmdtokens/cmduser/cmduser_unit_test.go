package cmduser_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdtokens/cmduser"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokens"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdUser(t *testing.T) {
	t.Run("Get displays user session metadata safely", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "meta-user"

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		exp := time.Now().Add(1 * time.Hour).UnixMilli()
		refExp := time.Now().Add(30 * 24 * time.Hour).UnixMilli()
		err = cli.Tokens().SetUserToken(ctx, &tokens.UserToken{
			UserId:                "usr_meta_test",
			UserAccessToken:       "secret.access.token",
			UserRefreshToken:      "secret.refresh.token",
			ExpiresAt:             exp,
			RefreshTokenExpiresAt: &refExp,
		})
		require.NoError(t, err)

		cmd := cmduser.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "table", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "get"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "usr_meta_test")
		require.Contains(t, out, "active")
		require.NotContains(t, out, "secret.access.token")
		require.NotContains(t, out, "secret.refresh.token")
	})

	t.Run("Expose outputs raw access token to stdout", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "expose-user"

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		err = cli.Tokens().SetUserToken(ctx, &tokens.UserToken{
			UserId:          "usr_expose_test",
			UserAccessToken: "jwt.raw.access.token.123",
		})
		require.NoError(t, err)

		cmd := cmduser.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "expose"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)
		require.Equal(t, "jwt.raw.access.token.123\n", outBuf.String())
	})

	t.Run("Export outputs raw refresh token to stdout", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "export-user"

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		err = cli.Tokens().SetUserToken(ctx, &tokens.UserToken{
			UserId:           "usr_export_test",
			UserAccessToken:  "jwt.acc",
			UserRefreshToken: "jwt.raw.refresh.token.456",
		})
		require.NoError(t, err)

		cmd := cmduser.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "export"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)
		require.Equal(t, "jwt.raw.refresh.token.456\n", outBuf.String())
	})

	t.Run("Delete evicts session from tokens.db", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "delete-user"

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		err = cli.Tokens().SetUserToken(ctx, &tokens.UserToken{
			UserId:          "usr_delete_test",
			UserAccessToken: "jwt.acc",
		})
		require.NoError(t, err)

		cmd := cmduser.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "delete"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)
		require.Contains(t, outBuf.String(), "usr_delete_test")

		ut, err := cli.Tokens().GetUserToken(ctx)
		require.NoError(t, err)
		require.Nil(t, ut)
	})

	t.Run("Fails get when unauthenticated", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())

		cmd := cmduser.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", "empty-prof", "Profile")
		cmd.SetArgs([]string{"-p", "empty-prof", "get"})

		err := cmd.ExecuteContext(ctx)
		require.Error(t, err)
		require.Equal(t, exitcode.Auth, exitcode.Code(err))
		require.Contains(t, strings.ToLower(err.Error()), "no user session")
	})
}
