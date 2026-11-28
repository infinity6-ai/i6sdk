package cmdcheck_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/infinity6-ai/i6sdk/apiuserz/userapitokencheck"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdtokens/cmdcheck"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokens"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdCheck(t *testing.T) {
	t.Run("Audits empty profile returns empty table", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())

		cmd := cmdcheck.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", "main", "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", "main"})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "TYPE")
		require.Contains(t, out, "TARGET")
	})

	t.Run("Audits user session with server diagnostics in JSON format", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/userz/token/check" {
				w.Header().Set("Content-Type", "application/json")
				resp := &userapitokencheck.RespPayload{
					Valid:     true,
					TokenKind: "user_jwt",
					Summary:   "Valid user session",
					Diagnostics: userapitokencheck.Diagnostics{
						UserId:               "usr_checked",
						TimeRemainingSeconds: 3600,
					},
				}
				_ = json.NewEncoder(w).Encode(resp)
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "audit-user"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{UserzUrl: srv.URL}, ctx))

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		err = cli.Tokens().SetUserToken(ctx, &tokens.UserToken{
			UserId:          "usr_checked",
			UserAccessToken: "test.jwt.token",
			ExpiresAt:       time.Now().Add(1 * time.Hour).UnixMilli(),
		})
		require.NoError(t, err)

		cmd := cmdcheck.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "json", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "--format", "json"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, `"target": "usr_checked"`)
		require.Contains(t, out, `"valid": true`)
		require.Contains(t, out, `"status": "valid"`)
	})

	t.Run("Audits robot key in table format", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "audit-robot"

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		err = cli.Tokens().SetDatasetRefreshToken(ctx, "fashion-lake", "ed25519-seed")
		require.NoError(t, err)

		cmd := cmdcheck.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "table", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "Ed25519 Dataset Refresh Key")
		require.Contains(t, out, "fashion-lake")
		require.Contains(t, out, "configured")
	})
}
