package cmdinspect_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/infinity6-ai/i6sdk/apiuserz/userapitokencheck"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdauth/cmdinspect"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokens"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdAuthInspect(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/userz/token/check" {
			var payload userapitokencheck.ReqPayload
			_ = json.NewDecoder(r.Body).Decode(&payload)
			w.Header().Set("Content-Type", "application/json")

			tok := payload.Token
			if tok == "arbitrary-valid-token" || tok == "stored-session-token" {
				resp := &userapitokencheck.RespPayload{
					Valid:     true,
					TokenKind: "user_jwt",
					Summary:   "User JWT verified against user service key",
					Diagnostics: userapitokencheck.Diagnostics{
						UserId:               "usr_inspect_999",
						TimeRemainingSeconds: 3600,
						ExpiresAt:            time.Now().Add(1 * time.Hour).Unix(),
					},
				}
				_ = json.NewEncoder(w).Encode(resp)
				return
			}

			resp := &userapitokencheck.RespPayload{
				Valid:     false,
				TokenKind: "user_jwt",
				ErrorCode: "invalid_signature",
				Summary:   "Invalid token signature",
			}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	t.Run("Inspects arbitrary token passed as argument", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "auth-inspect-arg"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{
			UserzUrl: srv.URL,
		}, ctx))

		cmd := cmdinspect.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"arbitrary-valid-token", "-p", profile})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "usr_inspect_999")
		require.Contains(t, out, "user_jwt")
		require.Contains(t, out, "yes")
	})

	t.Run("Inspects arbitrary token piped via stdin", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "auth-inspect-stdin"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{
			UserzUrl: srv.URL,
		}, ctx))

		cmd := cmdinspect.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "json", "Format")

		cmd.SetIn(strings.NewReader("arbitrary-valid-token\n"))
		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-", "-p", profile, "--format", "json"})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, `"valid": true`)
		require.Contains(t, out, `"user_id": "usr_inspect_999"`)
	})

	t.Run("Inspects active profile session when argument omitted", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "auth-inspect-profile"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{
			UserzUrl: srv.URL,
		}, ctx))

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		err = cli.Tokens().SetUserToken(ctx, &tokens.UserToken{
			UserId:          "usr_inspect_999",
			UserAccessToken: "stored-session-token",
		})
		require.NoError(t, err)

		cmd := cmdinspect.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "usr_inspect_999")
		require.Contains(t, out, "user_jwt")
	})

	t.Run("Returns error when no argument and no active session", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "empty-profile"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{
			UserzUrl: srv.URL,
		}, ctx))

		cmd := cmdinspect.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		cmd.SetArgs([]string{"-p", profile})

		err := cmd.ExecuteContext(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no active user session found")
	})
}
