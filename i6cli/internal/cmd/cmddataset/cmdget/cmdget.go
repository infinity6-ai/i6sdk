package cmdget

import (
	"context"
	"fmt"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/gox/commonz/validation/checked"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapidatasetsolutions"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'get' command under 'dataset'.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [<dataset>]",
		Short: "Display dataset metadata and solutions",
		Long: `Queries the dataset service to display metadata, active solution domains,
and configuration for the specified dataset.

If <dataset> is omitted, uses the default dataset configured for the active profile.`,
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
			cfg, _ := cliconfig.ResolveConfig(profile, c)

			argDs := ""
			if len(args) > 0 {
				argDs = args[0]
			}
			dsName, err := cliconfig.ResolveDataset(argDs, cfg, c)
			if err != nil {
				return err
			}

			var ds lakedsz.Dataset
			if err := checked.Set(&ds, dsName); err != nil {
				return exitcode.Wrap(exitcode.Usage, fmt.Errorf("invalid dataset %q: %w", dsName, err))
			}

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			api := &dszapidatasetsolutions.Api{
				Req: &dszapidatasetsolutions.Req{
					Params: &dszapidatasetsolutions.ReqParams{
						Dataset: ds,
					},
				},
			}

			status, err := cli.Apis().DszDatasetSolutions(c, api)
			if err != nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("failed to get dataset solutions: %w", err))
			}
			if status == 404 {
				return exitcode.New(exitcode.NotFound, fmt.Sprintf("dataset %q not found or access denied", dsName))
			}
			if status >= 400 {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("server returned error %d getting dataset solutions", status))
			}

			var solutions []string
			if api.Resp != nil && api.Resp.Payload != nil {
				for _, sol := range api.Resp.Payload.Solutions {
					if opt := sol.Optional(); opt.IsPresent() {
						solutions = append(solutions, opt.Must())
					}
				}
			}

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
			if err != nil {
				return err
			}

			solStr := "none"
			if len(solutions) > 0 {
				solStr = fmt.Sprintf("%v", solutions)
			}

			data := outfmt.Data{
				Headers: []string{"Property", "Value"},
				Rows: [][]string{
					{"Dataset", dsName},
					{"Active Solutions", solStr},
					{"Profile", profile},
				},
				QuietId: []string{dsName},
				Raw:     dsName,
				JSON: map[string]any{
					"dataset":   dsName,
					"solutions": solutions,
					"profile":   profile,
				},
			}

			logger.Debug(c, "dataset-get", map[string]any{
				"dataset": dsName,
			})

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}

	return cmd
}
