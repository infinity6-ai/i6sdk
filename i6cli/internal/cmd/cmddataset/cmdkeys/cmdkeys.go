package cmdkeys

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/gox/commonz/validation/checked"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefreshcreate"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefreshdelete"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefreshlist"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'keys' command under 'dataset'.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Aliases: []string{"robot-keys", "refresh-keys"},
		Short:   "Manage remote server machine robot keys for a dataset",
		Long: `Manage remote server-side machine robot Ed25519 refresh keys for a dataset.

These keys are provisioned and stored in the server registry and are used by CI/CD
runners, background workers, and automation services to authenticate as a robot account.`,
	}

	cmd.AddCommand(newListCmd(ctx))
	cmd.AddCommand(newCreateCmd(ctx))
	cmd.AddCommand(newDeleteCmd(ctx))

	return cmd
}

func newListCmd(ctx context.Context) *cobra.Command {
	var datasetFlag string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all registered machine robot keys for a dataset",
		Long:  `Queries the server to list all registered machine robot keys (JTIs, solutions, grants, expirations) provisioned for the dataset.`,
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

			api := &userapidatasettokenrefreshlist.Api{
				Req: &userapidatasettokenrefreshlist.Req{
					Params: &userapidatasettokenrefreshlist.ReqParams{
						Dataset: ds,
					},
				},
			}

			status, err := cli.Apis().UserDatasetTokenRefreshList(c, api)
			if err != nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("failed to list machine keys: %w", err))
			}
			if status >= 400 {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("server returned error %d listing machine keys", status))
			}

			var tokens []userapidatasettokenrefreshlist.RefreshTokenSummary
			if api.Resp != nil && api.Resp.Payload != nil {
				tokens = api.Resp.Payload.Tokens
			}

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
			if err != nil {
				return err
			}

			headers := []string{"JTI", "Dataset", "Solution", "Grants", "Created By", "Last Used", "Expires", "Status", "Annotation"}
			rows := make([][]string, 0, len(tokens))
			quietIds := make([]string, 0, len(tokens))

			for _, t := range tokens {
				lastUsed := "-"
				if t.LastUsedAt > 0 {
					lastUsed = time.UnixMilli(t.LastUsedAt).Format(time.RFC3339)
				}
				expires := "unlimited"
				if t.ExpiresAt > 0 {
					expires = time.UnixMilli(t.ExpiresAt).Format(time.RFC3339)
				}
				status := "active"
				if t.Revoked {
					status = "revoked"
				} else if t.ExpiresAt > 0 && time.Now().UnixMilli() > t.ExpiresAt {
					status = "expired"
				}

				solStr := formatSolution(t.Solution)
				dsStr := formatDataset(t.Dataset)

				rows = append(rows, []string{
					t.Jti,
					dsStr,
					solStr,
					t.Grants.Format(),
					t.CreatedBy,
					lastUsed,
					expires,
					status,
					t.Annotation,
				})
				quietIds = append(quietIds, t.Jti)
			}

			data := outfmt.Data{
				Headers: headers,
				Rows:    rows,
				QuietId: quietIds,
				JSON:    tokens,
			}

			logger.Debug(c, "dataset-keys-list", map[string]any{
				"dataset": dsName,
				"count":   len(tokens),
			})

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}

	cmd.Flags().StringVar(&datasetFlag, "dataset", "", "Target dataset identifier (or default_dataset / I6_DATASET)")
	return cmd
}

func newCreateCmd(ctx context.Context) *cobra.Command {
	var (
		datasetFlag string
		solFlag     string
		grantsFlag  string
		expSeconds  int64
		annotation  string
		outPath     string
		saveFlag    bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Provision a new machine robot Ed25519 key on the server",
		Long: `Provisions a new machine robot Ed25519 refresh key on the server for the dataset.

The server generates a new Ed25519 keypair, stores the public key record (JTI),
and returns the one-time private key artifact.

Security Warning: The private key is returned only once upon creation.
Store it securely or pass --save to automatically persist it to the active profile's tokens.db.`,
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

			if solFlag == "" {
				solFlag = "relevance"
			}
			var sol solution.Solution
			if err := checked.Set(&sol, solFlag); err != nil {
				return exitcode.Wrap(exitcode.Usage, fmt.Errorf("invalid solution %q: %w", solFlag, err))
			}

			parsedGrants := grants.Admin | grants.Write | grants.Read
			if grantsFlag != "" {
				g, err := grants.Parse(grantsFlag)
				if err != nil {
					return exitcode.Wrap(exitcode.Usage, fmt.Errorf("invalid grants %q: %w", grantsFlag, err))
				}
				parsedGrants = g
			}

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			if saveFlag {
				ut, _ := cli.Tokens().GetUserToken(c)
				if ut != nil && ut.UserId != "" {
					return exitcode.New(
						exitcode.Usage,
						fmt.Sprintf("active user session exists for user %q in profile %q: a profile cannot hold both a user session and a robot key. Please use a separate profile (e.g. 'i6 --profile worker dataset keys create --save ...')", ut.UserId, profile),
					)
				}
			}

			createApi := &userapidatasettokenrefreshcreate.Api{
				Req: &userapidatasettokenrefreshcreate.Req{
					Params: &userapidatasettokenrefreshcreate.ReqParams{
						Dataset: ds,
					},
					Payload: &userapidatasettokenrefreshcreate.ReqPayload{
						Solution:   sol,
						Grants:     parsedGrants,
						ExpSeconds: expSeconds,
						Annotation: strings.TrimSpace(annotation),
					},
				},
			}

			status, err := cli.Apis().UserDatasetTokenRefreshCreate(c, createApi)
			if err != nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("failed to create machine key: %w", err))
			}
			if status >= 400 || createApi.Resp == nil || createApi.Resp.Payload == nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("server returned error %d creating machine key", status))
			}

			p := createApi.Resp.Payload

			if saveFlag {
				if err := cli.Tokens().SetDatasetRefreshToken(c, dsName, p.PrivKey); err != nil {
					return exitcode.Wrap(exitcode.General, fmt.Errorf("failed to save machine key to profile tokens.db: %w", err))
				}
				logger.Info(c, "machine-key-saved", map[string]any{
					"profile": profile,
					"dataset": dsName,
					"jti":     p.Jti,
				})
			}

			if outPath != "" && outPath != "-" {
				if err := os.WriteFile(outPath, []byte(p.PrivKey), 0600); err != nil {
					return exitcode.Wrap(exitcode.General, fmt.Errorf("failed to write private key to %q: %w", outPath, err))
				}
			}

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
			if err != nil {
				return err
			}

			expStr := "unlimited"
			if p.ExpiresAt > 0 {
				expStr = time.UnixMilli(p.ExpiresAt).Format(time.RFC3339)
			}

			data := outfmt.Data{
				Headers: []string{"Property", "Value"},
				Rows: [][]string{
					{"JTI", p.Jti},
					{"Dataset", formatDataset(p.Dataset)},
					{"Solution", formatSolution(p.Solution)},
					{"Grants", p.Grants.Format()},
					{"Expires At", expStr},
					{"Private Key", p.PrivKey},
					{"Annotation", p.Annotation},
				},
				QuietId: []string{p.PrivKey},
				Raw:     p.PrivKey,
				JSON:    p,
			}

			if saveFlag {
				fmt.Fprintf(cmd.ErrOrStderr(), "Successfully provisioned and saved machine robot key for dataset %q to profile %q\n", dsName, profile)
			}

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}

	cmd.Flags().StringVar(&datasetFlag, "dataset", "", "Target dataset identifier (or default_dataset / I6_DATASET)")
	cmd.Flags().StringVar(&solFlag, "solution", "relevance", "Solution domain restriction")
	cmd.Flags().StringVar(&grantsFlag, "grants", "", "Grants bitmask (e.g. 'admin,write,read' or '0x0E')")
	cmd.Flags().Int64Var(&expSeconds, "exp", 0, "Expiration in seconds (0 = 100 years)")
	cmd.Flags().StringVar(&annotation, "annotation", "", "Description note for this key")
	cmd.Flags().StringVar(&outPath, "out", "", "Optional file path to write private key to")
	cmd.Flags().BoolVar(&saveFlag, "save", false, "Save the private key into the active profile's tokens.db")

	return cmd
}

func newDeleteCmd(ctx context.Context) *cobra.Command {
	var datasetFlag string

	cmd := &cobra.Command{
		Use:   "delete <jti>",
		Short: "Revoke a machine robot key on the server by JTI",
		Long:  `Instantly and permanently revokes a machine robot key from the server registry by its JTI public key identifier.`,
		Args:  cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			jti := strings.TrimSpace(args[0])
			if jti == "" {
				return exitcode.New(exitcode.Usage, "jti is required")
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

			delApi := &userapidatasettokenrefreshdelete.Api{
				Req: &userapidatasettokenrefreshdelete.Req{
					Params: &userapidatasettokenrefreshdelete.ReqParams{
						Dataset: ds,
						Jti:     jti,
					},
				},
			}

			status, err := cli.Apis().UserDatasetTokenRefreshDelete(c, delApi)
			if err != nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("failed to revoke machine key: %w", err))
			}
			if status >= 400 {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("server returned error %d revoking machine key", status))
			}

			logger.Info(c, "machine-key-revoked", map[string]any{
				"dataset": dsName,
				"jti":     jti,
			})

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully revoked machine robot key %q for dataset %q on server\n", jti, dsName)
			return nil
		},
	}

	cmd.Flags().StringVar(&datasetFlag, "dataset", "", "Target dataset identifier (or default_dataset / I6_DATASET)")
	return cmd
}

func formatSolution(s solution.Solution) string {
	if opt := s.Optional(); opt.IsPresent() {
		return opt.Must()
	}
	return "all"
}

func formatDataset(d lakedsz.Dataset) string {
	if opt := d.Optional(); opt.IsPresent() {
		return opt.Must()
	}
	return ""
}
