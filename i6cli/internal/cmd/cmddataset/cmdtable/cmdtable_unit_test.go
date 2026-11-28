package cmdtable_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/infinity6-ai/gox/fsz/fsz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapitablelist"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmddataset/cmdtable"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbpathz"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdTable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/ds/dataset/ds-table-test/table-list" {
			now := time.Now()
			w.Header().Set("Content-Type", "application/json")
			resp := &dszapitablelist.RespPayload{
				Files: []*dszapitablelist.File{
					{
						Path: laketbpathz.MustParse("i6tb://results.relevance-fashion-catalog/part.json"),
						Stat: &fsz.FileStat{
							Size:      1024,
							UpdatedAt: &now,
						},
					},
					{
						Path: laketbpathz.MustParse("i6tb://results.relevance-fashion-catalog/sub-dir/"),
						Stat: &fsz.FileStat{
							Size:      0,
							UpdatedAt: &now,
						},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		if r.URL.Path == "/api/ds/dataset/ds-table-test/table-file/i6tb://results.relevance-fashion-catalog/part.json" || r.URL.Path == "/api/ds/dataset/ds-table-test/table-file/results.relevance-fashion-catalog/part.json" {
			w.Header().Set("Location", "https://signed.url/results.relevance-fashion-catalog/part.json")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	t.Run("List returns tables and files", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "table-list"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{
			UserzUrl:       srv.URL,
			DecsuiteUrl:    srv.URL,
			DefaultDataset: "ds-table-test",
		}, ctx))

		cmd := cmdtable.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"list", "-p", profile})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "i6tb://results.relevance-fashion-catalog/part.json")
		require.Contains(t, out, "1.0 kB")
		require.Contains(t, out, "i6tb://results.relevance-fashion-catalog/sub-dir/")
		require.Contains(t, out, "(DIR)")
	})

	t.Run("Get returns signed URL location", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "table-get"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{
			UserzUrl:       srv.URL,
			DecsuiteUrl:    srv.URL,
			DefaultDataset: "ds-table-test",
		}, ctx))

		cmd := cmdtable.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"get", "i6tb://results.relevance-fashion-catalog/part.json", "-p", profile})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "https://signed.url/results.relevance-fashion-catalog/part.json")
	})
}
