package cmdlogout

import (
	"context"
	"fmt"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthlogout"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'logout' command under 'auth'.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Revoke session on server and remove stored credentials from profile",
		Long: `Revokes the active user refresh token family on the server (if authenticated)
and clears all cached tokens and credentials from the active profile's tokens.db.`,
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

			// If active human session exists, notify server to revoke refresh token
			ut, _ := cli.Tokens().GetUserToken(c)
			if ut != nil && ut.UserRefreshToken != "" {
				logoutApi := &userapiauthlogout.Api{
					Req: &userapiauthlogout.Req{
						Payload: &userapiauthlogout.ReqPayload{
							RefreshToken: ut.UserRefreshToken,
						},
					},
				}
				if _, err := cli.Apis().UserAuthLogout(c, logoutApi); err != nil {
					logger.Debug(c, "logout-server-error", map[string]any{"err": err.Error()})
				}
			}

			// Clear all tokens from the local store
			if err := cli.Tokens().ClearTokens(c); err != nil {
				return fmt.Errorf("failed to clear credentials from profile store: %w", err)
			}

			logger.Info(c, "logout-success", map[string]any{
				"profile": profile,
			})

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully logged out and cleared credentials for profile %q\n", profile)
			return nil
		},
	}

	return cmd
}
