package cmdstatus_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdauth/cmdstatus"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokens"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdStatus(t *testing.T) {
	t.Run("Displays unauthenticated state", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())

		cmd := cmdstatus.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", "main", "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "Unauthenticated")
		require.Contains(t, out, "main")
	})

	t.Run("Displays user session details", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())

		profile := "logged-in-user"
		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		exp := time.Now().Add(1 * time.Hour).UnixMilli()
		refExp := time.Now().Add(30 * 24 * time.Hour).UnixMilli()
		err = cli.Tokens().SetUserToken(ctx, &tokens.UserToken{
			UserId:                "usr_abc123",
			UserAccessToken:       "valid.access.jwt",
			UserRefreshToken:      "valid.refresh.jwt",
			ExpiresAt:             exp,
			RefreshTokenExpiresAt: &refExp,
		})
		require.NoError(t, err)

		cmd := cmdstatus.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "json", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "--format", "json"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, `"identity": "User Session"`)
		require.Contains(t, out, `"target_id": "usr_abc123"`)
		require.Contains(t, out, `"status": "active"`)
	})

	t.Run("Displays robot key identity with quiet format", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())

		profile := "service-account"
		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		err = cli.Tokens().SetDatasetRefreshToken(ctx, "fashion-prod", "robot-privkey-123")
		require.NoError(t, err)

		cmd := cmdstatus.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "quiet", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "--format", "quiet"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		require.Equal(t, "fashion-prod\n", outBuf.String())
	})
}
