package cmdlogin_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthdevicecode"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthdevicetoken"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdauth/cmdlogin"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdLogin(t *testing.T) {
	ctx := t.Context()

	t.Run("Fails when machine robot key already exists in profile", func(t *testing.T) {
		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		profile := "worker-profile"
		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)

		// Seed a robot refresh key
		err = cli.Tokens().SetDatasetRefreshToken(ctx, "fashion", "seed123")
		require.NoError(t, err)

		cmd := cmdlogin.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.SetArgs([]string{"-p", profile})

		execErr := cmd.ExecuteContext(ctx)
		require.Error(t, execErr)
		require.Equal(t, exitcode.Usage, exitcode.Code(execErr))
		require.Contains(t, execErr.Error(), "active robot key exists")
	})

	t.Run("Authorizes successfully through device code flow", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.URL.Path {
			case "/api/userz/auth/device/code":
				resp := &userapiauthdevicecode.RespPayload{
					DeviceCode:              "dev-abc",
					UserCode:                "CODE-999",
					VerificationUri:         "http://test.local/device",
					VerificationUriComplete: "http://test.local/device?code=CODE-999",
					ExpiresInSec:            60,
					IntervalSec:             0,
				}
				_ = json.NewEncoder(w).Encode(resp)
			case "/api/userz/auth/device/token":
				resp := &userapiauthdevicetoken.RespPayload{
					Status:                 "authorized",
					UserAccessToken:        "jwt.access.mock",
					UserRefreshToken:       "jwt.refresh.mock",
					UserId:                 "user_12345",
					ExpiresAt:              1800000000000,
					RefreshTokenExpiresAt: 1900000000000,
				}
				_ = json.NewEncoder(w).Encode(resp)
			default:
				http.NotFound(w, r)
			}
		}))
		defer srv.Close()

		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		profile := "human-user"
		cfg := &cliconfig.Config{UserzUrl: srv.URL}
		require.NoError(t, cliconfig.SaveConfig(profile, cfg, ctx))

		cmd := cmdlogin.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"-p", profile})

		err = cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "http://test.local/device?code=CODE-999")
		require.Contains(t, out, "user_12345")

		// Verify tokens were saved in store
		cli, err := cliconfig.NewClient(profile, ctx)
		require.NoError(t, err)
		ut, err := cli.Tokens().GetUserToken(ctx)
		require.NoError(t, err)
		require.NotNil(t, ut)
		require.Equal(t, "user_12345", ut.UserId)
		require.Equal(t, "jwt.access.mock", ut.UserAccessToken)
	})

	t.Run("Fails with Auth exit code when denied", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.URL.Path {
			case "/api/userz/auth/device/code":
				resp := &userapiauthdevicecode.RespPayload{
					DeviceCode:   "dev-denied",
					UserCode:     "DENY-000",
					ExpiresInSec: 60,
					IntervalSec:  0,
				}
				_ = json.NewEncoder(w).Encode(resp)
			case "/api/userz/auth/device/token":
				resp := &userapiauthdevicetoken.RespPayload{
					Status: "denied",
				}
				_ = json.NewEncoder(w).Encode(resp)
			default:
				http.NotFound(w, r)
			}
		}))
		defer srv.Close()

		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		profile := "deny-profile"
		require.NoError(t, cliconfig.SaveConfig(profile, &cliconfig.Config{UserzUrl: srv.URL}, ctx))

		cmd := cmdlogin.New(ctx)
		cmd.PersistentFlags().StringP("profile", "p", profile, "Profile")
		cmd.SetArgs([]string{"-p", profile})

		err = cmd.ExecuteContext(ctx)
		require.Error(t, err)
		require.Equal(t, exitcode.Auth, exitcode.Code(err))
		require.Contains(t, err.Error(), "denied")
	})
}
