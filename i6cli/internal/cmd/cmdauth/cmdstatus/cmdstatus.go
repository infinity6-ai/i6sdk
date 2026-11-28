package cmdstatus

import (
	"context"
	"fmt"
	"time"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'status' command under 'auth'.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"me"},
		Short:   "Display active authentication identity, tenant, and token expiry",
		Long: `Displays the active authentication state for the selected profile.

Shows identity type (User Session vs. Robot Service Account vs. Unauthenticated),
associated User ID or Dataset identifier, expiration times, and cached token status.`,
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

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
			if err != nil {
				return err
			}

			// 1. Check for Human User Session
			ut, _ := cli.Tokens().GetUserToken(c)
			// 2. Check for Machine Robot Refresh Key
			rkeys, _ := cli.Tokens().GetDatasetRefreshTokensMap(c)
			// 3. Check for Cached Dataset Access Tokens
			dtokens, _ := cli.Tokens().GetDatasetTokensMap(c)

			var (
				identity   = "Unauthenticated"
				targetId   = "(none)"
				statusStr  = "unauthenticated"
				accessExp  = "(none)"
				refreshExp = "(none)"
				rawVal     = ""
			)

			if ut != nil {
				identity = "User Session"
				targetId = ut.UserId
				if targetId == "" {
					targetId = "(user)"
				}
				if cli.Tokens().IsAccessTokenExpired(ut) {
					statusStr = "access token expired"
				} else {
					statusStr = "active"
				}
				if ut.ExpiresAt > 0 {
					accessExp = time.UnixMilli(ut.ExpiresAt).Format(time.RFC3339)
				}
				if ut.RefreshTokenExpiresAt != nil && *ut.RefreshTokenExpiresAt > 0 {
					refreshExp = time.UnixMilli(*ut.RefreshTokenExpiresAt).Format(time.RFC3339)
				}
				rawVal = ut.UserId
			} else if len(rkeys) > 0 {
				identity = "Machine Robot Account (Ed25519)"
				statusStr = "active"
				var dsList []string
				for ds := range rkeys {
					dsList = append(dsList, ds)
				}
				targetId = fmt.Sprintf("%v", dsList)
				if len(dsList) == 1 {
					targetId = dsList[0]
					rawVal = dsList[0]
				}
			}

			cachedCount := len(dtokens)

			data := outfmt.Data{
				Headers: []string{"Property", "Value"},
				Rows: [][]string{
					{"Profile", profile},
					{"Identity", identity},
					{"Target ID", targetId},
					{"Session Status", statusStr},
					{"Access Token Expiry", accessExp},
					{"Refresh Token Expiry", refreshExp},
					{"Cached Dataset Tokens", fmt.Sprintf("%d", cachedCount)},
				},
				QuietId: []string{rawVal},
				Raw:     rawVal,
				JSON: map[string]any{
					"profile":               profile,
					"identity":              identity,
					"target_id":             targetId,
					"status":                statusStr,
					"access_token_expiry":   accessExp,
					"refresh_token_expiry":  refreshExp,
					"cached_dataset_tokens": cachedCount,
				},
			}

			if rawVal == "" && format == outfmt.FormatQuiet {
				data.QuietId = []string{"unauthenticated"}
			}

			logger.Debug(c, "auth-status", map[string]any{
				"profile":  profile,
				"identity": identity,
				"status":   statusStr,
			})

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}

	return cmd
}
