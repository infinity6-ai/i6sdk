package cmdtoken

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokens"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'token' command under 'auth'.
func New(ctx context.Context) *cobra.Command {
	var (
		offline bool
		dataset string
	)

	cmd := &cobra.Command{
		Use:   "token [<token-string> | -]",
		Short: "Ingest and activate a credentials token into active profile",
		Long: `Universally ingests and activates a credential token into the active profile's tokens.db.

Automatically detects the token type:
  - Ed25519 Dataset Robot Key: Stored as the profile machine root credential (i6_rf:<dataset>).
  - User Refresh Token: Exchanged via /api/userz/auth/refresh to establish a full user session.
  - Standalone User Access Token: Stored as a non-refreshable user access token.
  - Standalone Dataset Bearer Token: Stored as a cached bearer token (i6_ds:<dataset>).

Use '-' to read the token string from standard input (stdin).
With --offline, network verification is bypassed and the token is parsed and stored locally.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			profileFlag, _ := cmd.Flags().GetString("profile")
			profile := cliconfig.ResolveProfile(profileFlag, c)

			var tokenStr string
			if len(args) == 1 && args[0] != "-" {
				tokenStr = strings.TrimSpace(args[0])
			} else {
				// Read from stdin
				inputBytes, err := io.ReadAll(cmd.InOrStdin())
				if err != nil {
					return fmt.Errorf("failed to read token from stdin: %w", err)
				}
				tokenStr = strings.TrimSpace(string(inputBytes))
			}

			if tokenStr == "" {
				return exitcode.New(
					exitcode.Usage,
					"token string cannot be empty: provide token as argument or pipe via stdin",
				)
			}

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			res, err := cli.IngestToken(c, tokenStr, tokens.IngestOptions{
				Offline: offline,
				Dataset: strings.TrimSpace(dataset),
			})
			if err != nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("token ingestion failed: %w", err))
			}

			target := res.Dataset
			if target == "" {
				target = res.UserId
			}
			if target == "" {
				target = profile
			}

			logger.Info(c, "token-ingested", map[string]any{
				"type":    string(res.Type),
				"target":  target,
				"profile": profile,
			})

			formatFlag, _ := cmd.Flags().GetString("format")
			if formatFlag != "" {
				format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
				if err != nil {
					return err
				}

				expStr := "never"
				if res.ExpiresAt > 0 {
					expStr = time.UnixMilli(res.ExpiresAt).Format(time.RFC3339)
				}

				data := outfmt.Data{
					Headers: []string{"Property", "Value"},
					Rows: [][]string{
						{"Token Type", string(res.Type)},
						{"Target Identity", target},
						{"Expires At", expStr},
						{"Profile", profile},
					},
					QuietId: []string{target},
					JSON: map[string]any{
						"type":       string(res.Type),
						"target":     target,
						"expires_at": res.ExpiresAt,
						"profile":    profile,
						"status":     "active",
					},
				}
				return outfmt.Render(cmd.OutOrStdout(), format, data)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully ingested %s token for %s (profile: %s)\n", res.Type, target, profile)
			return nil
		},
	}

	cmd.Flags().BoolVar(&offline, "offline", false, "Ingest token locally without server network calls")
	cmd.Flags().StringVar(&dataset, "dataset", "", "Target dataset identifier override")

	return cmd
}
