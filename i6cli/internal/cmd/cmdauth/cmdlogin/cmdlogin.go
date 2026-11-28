package cmdlogin

import (
	"context"
	"fmt"
	"time"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthdevicecode"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthdevicetoken"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokens"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'login' command under 'auth'.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate via browser using OAuth 2.0 Device Code flow",
		Long: `Initiates an interactive OAuth 2.0 Device Code authorization flow.

Displays a one-time user code and browser verification URL, then securely polls
the server until authorization is completed. Upon success, stores the user session
(access and refresh tokens) in the active profile's tokens.db.

Root Identity Invariant: If a machine robot key is already active in this profile,
login fails to prevent identity and privilege ambiguity.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			profileFlag, _ := cmd.Flags().GetString("profile")
			profile := cliconfig.ResolveProfile(profileFlag, c)

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			// Invariant check: ensure no active robot refresh key exists in this profile
			rkeys, err := cli.Tokens().GetDatasetRefreshTokensMap(c)
			if err == nil && len(rkeys) > 0 {
				var datasets []string
				for ds := range rkeys {
					datasets = append(datasets, ds)
				}
				return exitcode.New(
					exitcode.Usage,
					fmt.Sprintf("active robot key exists for dataset %v in profile %q: a profile cannot hold both a user session and a robot key. Please use a separate profile (e.g. 'i6 --profile developer auth login')", datasets, profile),
				)
			}

			// 1. Request device code
			devApi := &userapiauthdevicecode.Api{
				Req: &userapiauthdevicecode.Req{
					Payload: &userapiauthdevicecode.ReqPayload{
						App: "i6cli",
					},
				},
			}

			if _, err := cli.Apis().UserAuthDeviceCode(c, devApi); err != nil {
				return fmt.Errorf("failed to initiate device authorization: %w", err)
			}
			if devApi.Resp == nil || devApi.Resp.Payload == nil {
				return fmt.Errorf("empty device authorization response from server")
			}

			p := devApi.Resp.Payload
			if p.DeviceCode == "" {
				return fmt.Errorf("empty device code received from server")
			}
			verifyUrl := p.VerificationUriComplete
			if verifyUrl == "" {
				verifyUrl = p.VerificationUri
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nTo sign in, open this URL in a browser:\n  %s\n\nEnter code: %s\n\nWaiting for authorization...", verifyUrl, p.UserCode)

			// 2. Poll for authorization
			interval := time.Duration(p.IntervalSec) * time.Second
			if interval < 100*time.Millisecond {
				interval = 500 * time.Millisecond
			}

			timeoutSec := p.ExpiresInSec
			if timeoutSec <= 0 {
				timeoutSec = 300
			}
			deadline := time.Now().Add(time.Duration(timeoutSec) * time.Second)

			for {
				select {
				case <-c.Done():
					return c.Err()
				case <-time.After(interval):
				}

				if time.Now().After(deadline) {
					return exitcode.New(exitcode.Auth, "device authorization code expired")
				}

				pollApi := &userapiauthdevicetoken.Api{
					Req: &userapiauthdevicetoken.Req{
						Payload: &userapiauthdevicetoken.ReqPayload{
							DeviceCode: p.DeviceCode,
						},
					},
				}

				if _, err := cli.Apis().UserAuthDeviceToken(c, pollApi); err != nil {
					logger.Debug(c, "device-token-poll-error", map[string]any{"err": err.Error()})
					continue
				}
				if pollApi.Resp == nil || pollApi.Resp.Payload == nil {
					continue
				}

				switch pollApi.Resp.Payload.Status {
				case "pending":
					continue
				case "authorized":
					payload := pollApi.Resp.Payload
					ut := &tokens.UserToken{
						UserAccessToken:       payload.UserAccessToken,
						UserRefreshToken:      payload.UserRefreshToken,
						UserId:                payload.UserId,
						ExpiresAt:             payload.ExpiresAt,
						RefreshTokenExpiresAt: &payload.RefreshTokenExpiresAt,
					}
					if err := cli.Tokens().SetUserToken(c, ut); err != nil {
						return fmt.Errorf("failed to save user token session: %w", err)
					}
					logger.Info(c, "login-success", map[string]any{
						"user_id": payload.UserId,
						"profile": profile,
					})
					fmt.Fprintf(cmd.OutOrStdout(), "\nSuccessfully authenticated as user %s (profile: %s)\n", payload.UserId, profile)
					return nil
				case "denied":
					return exitcode.New(exitcode.Auth, "device authorization request was denied by user")
				case "expired":
					return exitcode.New(exitcode.Auth, "device authorization code has expired")
				default:
					continue
				}
			}
		},
	}

	return cmd
}
