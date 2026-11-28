package cmdstore_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapistorecreate"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmddataset/cmdstore"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdStore(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/ds/dataset/ds-store-test/store/list" {
			w.Header().Set("Content-Type", "application/json")
			resp := map[string]any{
				"stores": []map[string]any{
					{
						"id":   "store-123",
						"name": "Primary Store",
						"solution_domains": []string{"relevance_fashion"},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		if r.Method == "POST" && r.URL.Path == "/api/ds/dataset/ds-store-test/store/create" {
			var req dszapistorecreate.Req
			_ = json.NewDecoder(r.Body).Decode(&req)
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	t.Run("List returns stores for dataset", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "store-list"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{
			UserzUrl:       srv.URL,
			DecsuiteUrl:    srv.URL,
			DefaultDataset: "ds-store-test",
		}, ctx))

		cmd := cmdstore.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"list", "-p", profile})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "store-123")
		require.Contains(t, out, "Primary Store")
	})

	t.Run("Create registers new store", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "store-create"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{
			UserzUrl:       srv.URL,
			DecsuiteUrl:    srv.URL,
			DefaultDataset: "ds-store-test",
		}, ctx))

		cmd := cmdstore.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"create", "--name", "New Backup Store", "-p", profile})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "Successfully registered store \"New Backup Store\"")
	})
}
