package cmdkeys_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefreshcreate"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefreshlist"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmddataset/cmdkeys"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokens"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdKeys(t *testing.T) {
	t.Run("List returns table of registered keys", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" && r.URL.Path == "/api/userz/dataset-refresh-token/ds-keys-test" {
				w.Header().Set("Content-Type", "application/json")
				resp := &userapidatasettokenrefreshlist.RespPayload{
					Tokens: []userapidatasettokenrefreshlist.RefreshTokenSummary{
						{
							Jti:        "jti_test_123",
							Dataset:    lakedsz.MustParse("ds-keys-test"),
							Solution:   solution.MustParse("relevance"),
							Grants:     grants.Admin | grants.Write | grants.Read,
							CreatedBy:  "usr_admin_1",
							CreatedAt:  time.Now().UnixMilli(),
							ExpiresAt:  time.Now().Add(365 * 24 * time.Hour).UnixMilli(),
							Annotation: "CI Runner Key",
						},
					},
				}
				_ = json.NewEncoder(w).Encode(resp)
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "keys-list"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{
			UserzUrl:       srv.URL,
			DefaultDataset: "ds-keys-test",
		}, ctx))

		cmd := cmdkeys.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"list", "-p", profile})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "jti_test_123")
		require.Contains(t, out, "ds-keys-test")
		require.Contains(t, out, "CI Runner Key")
	})

	t.Run("Create provisions new machine key on server", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" && r.URL.Path == "/api/userz/dataset-refresh-token/ds-keys-test" {
				w.Header().Set("Content-Type", "application/json")
				resp := &userapidatasettokenrefreshcreate.RespPayload{
					Jti:        "jti_new_456",
					Dataset:    lakedsz.MustParse("ds-keys-test"),
					PrivKey:    "MC4CAQAwBQYDK2VwBCIEIP-fake-priv-key",
					Solution:   solution.MustParse("relevance"),
					Grants:     grants.Admin | grants.Write | grants.Read,
					ExpiresAt:  time.Now().Add(100 * 365 * 24 * time.Hour).UnixMilli(),
					Annotation: "Production Worker",
				}
				_ = json.NewEncoder(w).Encode(resp)
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "keys-create"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{
			UserzUrl:       srv.URL,
			DefaultDataset: "ds-keys-test",
		}, ctx))

		cmd := cmdkeys.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"create", "-p", profile, "--annotation", "Production Worker"})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "jti_new_456")
		require.Contains(t, out, "MC4CAQAwBQYDK2VwBCIEIP-fake-priv-key")
	})

	t.Run("Create with --save enforces mutual exclusivity on user session", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "user-session-profile"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{
			DefaultDataset: "ds-keys-test",
		}, ctx))

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		err = cli.Tokens().SetUserToken(ctx, &tokens.UserToken{
			UserId:          "usr_exclusive",
			UserAccessToken: "jwt.token",
		})
		require.NoError(t, err)

		cmd := cmdkeys.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		cmd.SetArgs([]string{"create", "-p", profile, "--save"})

		err = cmd.ExecuteContext(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "active user session exists")
	})

	t.Run("Delete revokes machine key on server by JTI", func(t *testing.T) {
		deleted := false
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "DELETE" && r.URL.Path == "/api/userz/dataset-refresh-token/ds-keys-test/jti/jti_to_delete" {
				deleted = true
				w.WriteHeader(http.StatusOK)
				return
			}
			http.NotFound(w, r)
		}))
		defer srv.Close()

		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "keys-delete"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{
			UserzUrl:       srv.URL,
			DefaultDataset: "ds-keys-test",
		}, ctx))

		cmd := cmdkeys.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"delete", "jti_to_delete", "-p", profile})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)
		require.True(t, deleted)
		require.Contains(t, outBuf.String(), "Successfully revoked machine robot key")
	})
}
