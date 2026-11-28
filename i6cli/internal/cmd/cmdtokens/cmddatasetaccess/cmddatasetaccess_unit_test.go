package cmddatasetaccess_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdtokens/cmddatasetaccess"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokens"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdDatasetAccess(t *testing.T) {
	t.Run("List displays all cached dataset tokens", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "cached-tokens-profile"

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		exp := time.Now().Add(2 * time.Hour).UnixMilli()
		require.NoError(t, cli.Tokens().SetDatasetToken(ctx, "ds-one", &tokens.DatasetToken{
			AccessToken: "bearer.token.one",
			ExpiresAt:   exp,
			Grants:      grants.Read | grants.Write,
		}))
		require.NoError(t, cli.Tokens().SetDatasetToken(ctx, "ds-two", &tokens.DatasetToken{
			AccessToken: "bearer.token.two",
			ExpiresAt:   exp,
			Grants:      grants.Read,
		}))

		cmd := cmddatasetaccess.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "table", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "list"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "ds-one")
		require.Contains(t, out, "ds-two")
		require.Contains(t, out, "write|read")
	})

	t.Run("Get displays token metadata safely without revealing secret token", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "safe-get-profile"

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		exp := time.Now().Add(1 * time.Hour).UnixMilli()
		require.NoError(t, cli.Tokens().SetDatasetToken(ctx, "my-target-ds", &tokens.DatasetToken{
			AccessToken: "super.secret.token.string",
			ExpiresAt:   exp,
			Grants:      grants.Read,
		}))

		cmd := cmddatasetaccess.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "table", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "get", "--dataset", "my-target-ds"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "my-target-ds")
		require.Contains(t, out, "active")
		require.NotContains(t, out, "super.secret.token.string")
	})

	t.Run("Expose outputs raw token string to stdout", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "expose-profile"

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		require.NoError(t, cli.Tokens().SetDatasetToken(ctx, "target-ds", &tokens.DatasetToken{
			AccessToken: "exact.raw.jwt.token",
		}))

		cmd := cmddatasetaccess.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "expose", "--dataset", "target-ds"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)
		require.Equal(t, "exact.raw.jwt.token\n", outBuf.String())
	})

	t.Run("Delete evicts cached token from store", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "evict-profile"

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		require.NoError(t, cli.Tokens().SetDatasetToken(ctx, "evict-ds", &tokens.DatasetToken{
			AccessToken: "token-to-evict",
		}))

		cmd := cmddatasetaccess.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "delete", "--dataset", "evict-ds"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)
		require.Contains(t, outBuf.String(), "evict-ds")

		dt, err := cli.Tokens().GetDatasetToken(ctx, "evict-ds")
		require.NoError(t, err)
		require.Nil(t, dt)
	})
}
