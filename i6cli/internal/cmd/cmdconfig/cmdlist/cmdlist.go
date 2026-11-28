package cmdlist

import (
	"context"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'list' command under 'config'.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all local configuration profiles",
		Long: `Lists all local configuration profiles discovered in the base configuration directory.
The currently active profile is marked with an asterisk (*).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			profileFlag, _ := cmd.Flags().GetString("profile")
			profiles, err := cliconfig.ListProfiles(profileFlag, c)
			if err != nil {
				return err
			}

			headers := []string{"Active", "Profile", "Default Dataset", "Userz Service", "Config", "Tokens"}
			rows := make([][]string, 0, len(profiles))
			quietIds := make([]string, 0, len(profiles))

			for _, p := range profiles {
				activeMarker := ""
				if p.Active {
					activeMarker = "*"
				}

				ds := p.DefaultDataset
				if ds == "" {
					ds = "(none)"
				}

				cfgStr := "no"
				if p.HasConfig {
					cfgStr = "yes"
				}

				tokStr := "no"
				if p.HasTokens {
					tokStr = "yes"
				}

				rows = append(rows, []string{
					activeMarker,
					p.Name,
					ds,
					p.UserzUrl,
					cfgStr,
					tokStr,
				})
				quietIds = append(quietIds, p.Name)
			}

			data := outfmt.Data{
				Headers: headers,
				Rows:    rows,
				QuietId: quietIds,
				JSON:    profiles,
			}

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
			if err != nil {
				return err
			}

			logger.Debug(c, "config-list", map[string]any{
				"count":  len(profiles),
				"format": string(format),
			})

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}

	return cmd
}
