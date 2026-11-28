package cmdstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/gox/commonz/validation/checked"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapistorecreate"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapistorelist"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/dszstore"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solutiondomain"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'store' command under 'dataset'.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store",
		Short: "Manage connected cloud lake storage",
		Long:  `Lists and configures physical cloud storage buckets connected to the dataset.`,
	}

	cmd.AddCommand(newListCmd(ctx))
	cmd.AddCommand(newCreateCmd(ctx))

	return cmd
}

func newListCmd(ctx context.Context) *cobra.Command {
	var datasetFlag string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List connected stores in a dataset",
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

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			api := &dszapistorelist.Api{
				Req: &dszapistorelist.Req{
					Params: &dszapistorelist.ReqParams{
						Dataset: ds,
					},
				},
			}

			status, err := cli.Apis().DszListStores(c, api)
			if err != nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("failed to list stores: %w", err))
			}
			if status >= 400 {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("server returned error %d listing stores", status))
			}

			var stores []dszstore.DszStore
			if api.Resp != nil && api.Resp.Payload != nil {
				stores = api.Resp.Payload.Stores
			}

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
			if err != nil {
				return err
			}

			headers := []string{"Store ID", "Name", "Solutions"}
			rows := make([][]string, 0, len(stores))
			quietIds := make([]string, 0, len(stores))

			for _, s := range stores {
				solStrs := []string{}
				for _, sol := range s.SolutionDomains {
					if opt := sol.Optional(); opt.IsPresent() {
						solStrs = append(solStrs, opt.Must().String())
					}
				}
				solStr := "-"
				if len(solStrs) > 0 {
					solStr = fmt.Sprintf("%v", solStrs)
				}
				idStr := s.Id.String()

				rows = append(rows, []string{
					idStr,
					s.Name,
					solStr,
				})
				quietIds = append(quietIds, idStr)
			}

			data := outfmt.Data{
				Headers: headers,
				Rows:    rows,
				QuietId: quietIds,
				JSON:    stores,
			}

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}

	cmd.Flags().StringVar(&datasetFlag, "dataset", "", "Target dataset identifier")

	return cmd
}

func newCreateCmd(ctx context.Context) *cobra.Command {
	var (
		datasetFlag   string
		idFlag        string
		nameFlag      string
		solutionsFlag string
	)

	cmd := &cobra.Command{
		Use:   "create --id <id> [--name <name>] [--solutions <solution_domain,...>]",
		Short: "Register a new cloud store in a dataset",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			idFlag = strings.TrimSpace(idFlag)
			if idFlag == "" {
				return exitcode.New(exitcode.Usage, "--id is required")
			}

			nameFlag = strings.TrimSpace(nameFlag)
			if nameFlag == "" {
				nameFlag = idFlag
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

			var stId storeid.StoreId
			if err := checked.Set(&stId, idFlag); err != nil {
				return exitcode.Wrap(exitcode.Usage, fmt.Errorf("invalid store id %q: %w", idFlag, err))
			}

			var solDomains []solutiondomain.SolutionDomain
			for _, s := range strings.Split(solutionsFlag, ",") {
				s = strings.TrimSpace(s)
				if s == "" {
					continue
				}
				sd, err := solutiondomain.ParseString(s)
				if err != nil {
					return exitcode.Wrap(exitcode.Usage, fmt.Errorf("invalid solution domain %q (expected format like 'relevance_fashion'): %w", s, err))
				}
				solDomains = append(solDomains, sd)
			}
			if len(solDomains) == 0 {
				solDomains = append(solDomains, solutiondomain.MustParseString("relevance_fashion"))
			}

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			api := &dszapistorecreate.Api{
				Req: &dszapistorecreate.Req{
					Params: &dszapistorecreate.ReqParams{
						Dataset: ds,
					},
					Payload: &dszapistorecreate.ReqPayload{
						Store: dszstore.DszStore{
							Id:              stId,
							Name:            nameFlag,
							SolutionDomains: solDomains,
						},
					},
				},
			}

			status, err := cli.Apis().DszCreateStore(c, api)
			if err != nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("failed to create store: %w", err))
			}
			if status >= 400 {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("server returned error %d creating store", status))
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully registered store %q (id: %q) in dataset %q\n", nameFlag, idFlag, dsName)
			return nil
		},
	}

	cmd.Flags().StringVar(&datasetFlag, "dataset", "", "Target dataset identifier")
	cmd.Flags().StringVar(&idFlag, "id", "", "Unique store ID (required)")
	cmd.Flags().StringVar(&nameFlag, "name", "", "Human-readable name of the store (defaults to store ID)")
	cmd.Flags().StringVar(&solutionsFlag, "solutions", "relevance_fashion", "Comma-separated solution domains (e.g. 'relevance_fashion')")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
