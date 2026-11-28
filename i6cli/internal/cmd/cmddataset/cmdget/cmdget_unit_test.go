package cmdget_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapidatasetsolutions"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmddataset/cmdget"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdDatasetGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/ds/dataset/ds-get-test/solutions" {
			w.Header().Set("Content-Type", "application/json")
			resp := &dszapidatasetsolutions.RespPayload{
				Dataset: lakedsz.MustParse("ds-get-test"),
				Solutions: []solution.Solution{
					solution.MustParse("relevance"),
					solution.MustParse("targeting"),
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	t.Run("Gets dataset solutions", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "get-dataset"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{
			DecsuiteUrl:    srv.URL,
			DefaultDataset: "ds-get-test",
		}, ctx))

		cmd := cmdget.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "ds-get-test")
		require.Contains(t, out, "relevance")
		require.Contains(t, out, "targeting")
	})
}
