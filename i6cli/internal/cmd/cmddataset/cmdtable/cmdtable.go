package cmdtable

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/gox/commonz/validation/checked"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapitableget"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapitablelist"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbpathz"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'table' command under 'dataset'.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "table",
		Short: "Manage dataset tables and files",
		Long:  `Inspect and browse lake table schemas, partitions, and files within a dataset.`,
	}

	cmd.AddCommand(newListCmd(ctx))
	cmd.AddCommand(newGetCmd(ctx))

	return cmd
}

func newListCmd(ctx context.Context) *cobra.Command {
	var (
		datasetFlag string
		recursive   bool
		stat        bool
		limit       int
	)

	cmd := &cobra.Command{
		Use:   "list [<path>]",
		Short: "List tables and files in a dataset",
		Long: `Lists lake files and tables at the specified path within a dataset.
If no path is provided, lists the root tables of the dataset.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			profileFlag, _ := cmd.Flags().GetString("profile")
			profile := cliconfig.ResolveProfile(profileFlag, c)
			cfg, _ := cliconfig.ResolveConfig(profile, c)

			dsName, err := cliconfig.ResolveDataset(datasetFlag, cfg, c)
			if err != nil {
				return err
			}

			var ds lakedsz.Dataset
			if err := checked.Set(&ds, dsName); err != nil {
				return exitcode.Wrap(exitcode.Usage, fmt.Errorf("invalid dataset %q: %w", dsName, err))
			}

			var targetPath laketbpathz.Path
			if len(args) > 0 {
				if err := checked.Set(&targetPath, args[0]); err != nil {
					return exitcode.Wrap(exitcode.Usage, fmt.Errorf("invalid path %q: %w", args[0], err))
				}
			}

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			api := &dszapitablelist.Api{
				Req: &dszapitablelist.Req{
					Params: &dszapitablelist.ReqParams{
						Dataset: ds,
					},
					Query: &dszapitablelist.ReqQuery{
						Path:      targetPath,
						Recursive: recursive,
						Stat:      stat,
						Limit:     limit,
					},
				},
			}

			status, err := cli.Apis().DszTableList(c, api)
			if err != nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("failed to list tables: %w", err))
			}
			if status == 404 {
				return exitcode.New(exitcode.NotFound, fmt.Sprintf("dataset or path not found: %q", dsName))
			}
			if status >= 400 {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("server returned error %d listing tables", status))
			}

			var files []*dszapitablelist.File
			if api.Resp != nil && api.Resp.Payload != nil {
				files = api.Resp.Payload.Files
			}

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
			if err != nil {
				return err
			}

			headers := []string{"Path", "Size", "Modified", "Signed URL"}
			rows := make([][]string, 0, len(files))
			quietIds := make([]string, 0, len(files))

			for _, f := range files {
				pathStr := f.Path.Get()
				sizeStr := "-"
				modStr := "-"
				if f.Stat != nil {
					if strings.HasSuffix(pathStr, "/") {
						sizeStr = "(DIR)"
					} else {
						sizeStr = humanize.Bytes(uint64(f.Stat.Size))
					}
					if f.Stat.UpdatedAt != nil {
						modStr = f.Stat.UpdatedAt.Format(time.RFC3339)
					}
				}
				rows = append(rows, []string{
					pathStr,
					sizeStr,
					modStr,
					f.SignedUrl,
				})
				quietIds = append(quietIds, pathStr)
			}

			data := outfmt.Data{
				Headers: headers,
				Rows:    rows,
				QuietId: quietIds,
				JSON:    files,
			}

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}

	cmd.Flags().StringVar(&datasetFlag, "dataset", "", "Target dataset identifier")
	cmd.Flags().BoolVar(&recursive, "recursive", false, "List recursively")
	cmd.Flags().BoolVar(&stat, "stat", true, "Include file stats (size, modified time)")
	cmd.Flags().IntVar(&limit, "limit", 100, "Limit the number of results")

	return cmd
}

func newGetCmd(ctx context.Context) *cobra.Command {
	var datasetFlag string

	cmd := &cobra.Command{
		Use:   "get <path>",
		Short: "Get presigned URL for a table file",
		Long:  `Generates a presigned GET URL for a specific table file.`,
		Args:  cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			profileFlag, _ := cmd.Flags().GetString("profile")
			profile := cliconfig.ResolveProfile(profileFlag, c)
			cfg, _ := cliconfig.ResolveConfig(profile, c)

			dsName, err := cliconfig.ResolveDataset(datasetFlag, cfg, c)
			if err != nil {
				return err
			}

			var ds lakedsz.Dataset
			if err := checked.Set(&ds, dsName); err != nil {
				return exitcode.Wrap(exitcode.Usage, fmt.Errorf("invalid dataset %q: %w", dsName, err))
			}

			var targetPath laketbpathz.Path
			if err := checked.Set(&targetPath, args[0]); err != nil {
				return exitcode.Wrap(exitcode.Usage, fmt.Errorf("invalid path %q: %w", args[0], err))
			}

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			api := &dszapitableget.Api{
				Req: &dszapitableget.Req{
					Params: &dszapitableget.ReqParams{
						Dataset: ds,
						Path:    targetPath,
					},
				},
			}

			status, err := cli.Apis().DszTableGet(c, api)
			if err != nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("failed to get table file: %w", err))
			}
			if status == 404 {
				return exitcode.New(exitcode.NotFound, fmt.Sprintf("file %q not found", targetPath.Get()))
			}

			// In redirect style API, successful presigned URLs return HTTP 307 with Location header
			location := ""
			if api.Resp != nil && api.Resp.Headers != nil {
				location = api.Resp.Headers.Location
			}
			if location == "" {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("server did not return a Location header for signed URL (status %d)", status))
			}

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatRaw)
			if err != nil {
				return err
			}

			data := outfmt.Data{
				Headers: []string{"Property", "Value"},
				Rows: [][]string{
					{"Path", targetPath.Get()},
					{"Signed URL", location},
				},
				QuietId: []string{location},
				Raw:     location,
				JSON: map[string]string{
					"path":       targetPath.Get(),
					"signed_url": location,
				},
			}

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}

	cmd.Flags().StringVar(&datasetFlag, "dataset", "", "Target dataset identifier")

	return cmd
}
