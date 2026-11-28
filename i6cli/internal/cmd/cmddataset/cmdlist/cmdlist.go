package cmdlist

import (
	"context"
	"fmt"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapidatasetlist"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'list' command under 'dataset'.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all accessible datasets",
		Long: `Queries the dataset service to list all datasets accessible to the authenticated identity.

For human user sessions, queries the server registry for all assigned datasets and permissions.
For robot profiles, lists the configured dataset.`,
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

			cfg, _ := cliconfig.ResolveConfig(profile, c)

			var datasets []*dszapidatasetlist.Dataset

			ut, _ := cli.Tokens().GetUserToken(c)
			if ut != nil && ut.UserId != "" {
				api := &dszapidatasetlist.Api{
					Req: &dszapidatasetlist.Req{
						Params: &dszapidatasetlist.Params{
							UserId: ut.UserId,
						},
					},
				}

				status, err := cli.Apis().DszListDatasets(c, api)
				if err != nil {
					return exitcode.Wrap(exitcode.General, fmt.Errorf("failed to list datasets: %w", err))
				}
				if status >= 400 {
					return exitcode.Wrap(exitcode.General, fmt.Errorf("server returned error %d listing datasets", status))
				}
				if api.Resp != nil && api.Resp.Payload != nil {
					datasets = api.Resp.Payload.Datasets
				}
			} else {
				// Robot account or local default
				if cfg != nil && cfg.DefaultDataset != "" {
					datasets = append(datasets, &dszapidatasetlist.Dataset{
						Id:   cfg.DefaultDataset,
						Name: cfg.DefaultDataset,
					})
				}
			}

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Grants"}
			rows := make([][]string, 0, len(datasets))
			quietIds := make([]string, 0, len(datasets))

			for _, ds := range datasets {
				grantsStr := ds.Grants.Format()
				if grantsStr == "" {
					grantsStr = "-"
				}
				rows = append(rows, []string{
					ds.Id,
					ds.Name,
					grantsStr,
				})
				quietIds = append(quietIds, ds.Id)
			}

			data := outfmt.Data{
				Headers: headers,
				Rows:    rows,
				QuietId: quietIds,
				JSON:    datasets,
			}

			logger.Debug(c, "dataset-list", map[string]any{
				"count": len(datasets),
			})

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}

	return cmd
}
