package cmdraw_test

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/infinity6-ai/i6sdk/apiuserz/userapitokencheck"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdtokens/cmdraw"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdRaw(t *testing.T) {
	var privKeyB64 string
	for {
		_, priv, err := ed25519.GenerateKey(nil)
		require.NoError(t, err)
		privKeyB64 = base64.RawURLEncoding.EncodeToString(priv.Seed())
		if !strings.HasPrefix(privKeyB64, "-") {
			break
		}
	}

	t.Run("Raw check calls server statelessly and displays diagnostics", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/userz/token/check" {
				w.Header().Set("Content-Type", "application/json")
				resp := &userapitokencheck.RespPayload{
					Valid:     true,
					TokenKind: "dataset_access_token",
					Summary:   "Valid bearer token",
					Diagnostics: userapitokencheck.Diagnostics{
						Dataset:              "stateless-ds",
						TimeRemainingSeconds: 600,
					},
				}
				_ = json.NewEncoder(w).Encode(resp)
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "raw-prof"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{UserzUrl: srv.URL}, ctx))

		cmd := cmdraw.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "table", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile, "check", "arbitrary.raw.token.string"})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "stateless-ds")
		require.Contains(t, out, "valid")
	})

	t.Run("Raw sign-request generates valid signature envelope", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())

		cmd := cmdraw.New(ctx)
		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{
			"sign-request",
			"--key", privKeyB64,
			"--token", "sample-access-token",
			"--method", "POST",
			"--path", "/api/ingest/upload",
		})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		sig := strings.TrimSpace(outBuf.String())
		require.NotEmpty(t, sig)
		// Composite signature header has <envelope>.<token> format
		require.Contains(t, sig, ".")
		require.Contains(t, sig, "sample-access-token")
	})

	t.Run("Raw sign-request fails when required flags are missing", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())

		cmd := cmdraw.New(ctx)
		cmd.SetArgs([]string{
			"sign-request",
			"--key", privKeyB64,
		})

		err := cmd.ExecuteContext(ctx)
		require.Error(t, err)
		require.Equal(t, exitcode.Usage, exitcode.Code(err))
	})
}
