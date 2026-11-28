package cmdinspect_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/infinity6-ai/i6sdk/apiuserz/userapitokencheck"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdtokens/cmdinspect"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokens"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdInspect(t *testing.T) {
	t.Run("Inspects empty profile returns empty table", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())

		cmd := cmdinspect.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", "main", "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", "main"})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "SOURCE")
		require.Contains(t, out, "TARGET")
	})

	t.Run("Inspects user session in json format with diagnostics", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/userz/token/check" {
				w.Header().Set("Content-Type", "application/json")
				resp := &userapitokencheck.RespPayload{
					Valid:     true,
					TokenKind: "user_jwt",
					Summary:   "User JWT verified against userz key",
					Diagnostics: userapitokencheck.Diagnostics{
						UserId:               "usr_inspect_123",
						TimeRemainingSeconds: 7200,
						ExpiresAt:            1800000000,
					},
				}
				_ = json.NewEncoder(w).Encode(resp)
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "inspect-user"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{UserzUrl: srv.URL}, ctx))

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		err = cli.Tokens().SetUserToken(ctx, &tokens.UserToken{
			UserId:          "usr_inspect_123",
			UserAccessToken: "jwt.inspect.token",
			ExpiresAt:       time.Now().Add(2 * time.Hour).UnixMilli(),
		})
		require.NoError(t, err)

		cmd := cmdinspect.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "json", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "--format", "json"})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, `"credential_source": "user_session"`)
		require.Contains(t, out, `"target": "usr_inspect_123"`)
		require.Contains(t, out, `"token_kind": "user_jwt"`)
		require.Contains(t, out, `"valid": true`)
		require.Contains(t, out, `"summary": "User JWT verified against userz key"`)
	})
}
