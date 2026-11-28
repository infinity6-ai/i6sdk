package cmdlogout_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthlogout"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdauth/cmdlogout"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokens"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdLogout(t *testing.T) {
	ctx := t.Context()

	t.Run("Logs out when unauthenticated clears store without error", func(t *testing.T) {
		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		cmd := cmdlogout.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", "main", "Profile")
		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", "main"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)
		require.Contains(t, outBuf.String(), "cleared credentials for profile \"main\"")
	})

	t.Run("Logs out user session calls server and clears tokens", func(t *testing.T) {
		serverCalled := false
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/userz/auth/logout" {
				serverCalled = true
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(&userapiauthlogout.RespPayload{Success: true})
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		profile := "active-session"
		cfg := &cliconfig.Config{UserzUrl: srv.URL}
		require.NoError(t, cliconfig.SaveConfig(profile, cfg, ctx))

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		err = cli.Tokens().SetUserToken(ctx, &tokens.UserToken{
			UserId:           "user_to_logout",
			UserAccessToken:  "access.jwt",
			UserRefreshToken: "refresh.jwt",
		})
		require.NoError(t, err)

		cmd := cmdlogout.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)
		require.True(t, serverCalled)

		// Verify tokens were evicted from store
		ut, err := cli.Tokens().GetUserToken(ctx)
		require.NoError(t, err)
		require.Nil(t, ut)
	})

	t.Run("Logs out robot key clears stored credentials", func(t *testing.T) {
		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		profile := "robot-profile"
		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		err = cli.Tokens().SetDatasetRefreshToken(ctx, "ds-to-clear", "key123")
		require.NoError(t, err)

		cmd := cmdlogout.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		rkeys, err := cli.Tokens().GetDatasetRefreshTokensMap(ctx)
		require.NoError(t, err)
		require.Empty(t, rkeys)
	})
}
