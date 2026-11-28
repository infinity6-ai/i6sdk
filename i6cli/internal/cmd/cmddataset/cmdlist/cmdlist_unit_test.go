package cmdlist_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapidatasetlist"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmddataset/cmdlist"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokens"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdDatasetList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/ds/user-id/usr_list_123/datasets" {
			w.Header().Set("Content-Type", "application/json")
			resp := &dszapidatasetlist.Payload{
				Datasets: []*dszapidatasetlist.Dataset{
					{
						Id:     "ds-list-test-1",
						Name:   "List Test 1",
						Grants: grants.Admin | grants.Write | grants.Read,
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	t.Run("Lists datasets for active user session", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "list-user"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{
			UserzUrl:    srv.URL,
			DecsuiteUrl: srv.URL,
		}, ctx))

		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		err = cli.Tokens().SetUserToken(ctx, &tokens.UserToken{
			UserId:          "usr_list_123",
			UserAccessToken: "jwt.token",
		})
		require.NoError(t, err)

		cmd := cmdlist.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "ds-list-test-1")
		require.Contains(t, out, "List Test 1")
		require.Contains(t, out, "admin|write|read")
	})

	t.Run("Lists default dataset for robot profile", func(t *testing.T) {
		ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())
		profile := "list-robot"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{
			DefaultDataset: "ds-robot-default",
		}, ctx))

		cmd := cmdlist.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.PersistentFlags().String("format", "", "Format")

		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "ds-robot-default")
	})
}
