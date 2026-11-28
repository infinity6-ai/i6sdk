package cmddatasetaccess

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenaccessattenuate"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapitokencheck"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'dataset-access' command under 'tokens'.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dataset-access",
		Aliases: []string{"bearer"},
		Short:   "Manage short-lived dataset access tokens cached in profile",
		Long: `Manages short-lived bearer access tokens cached in the active profile's tokens.db.
A profile can cache tokens across multiple datasets concurrently.`,
	}

	cmd.AddCommand(newListCmd(ctx))
	cmd.AddCommand(newGetCmd(ctx))
	cmd.AddCommand(newExposeCmd(ctx))
	cmd.AddCommand(newAttenuationCmd(ctx))
	cmd.AddCommand(newCheckCmd(ctx))
	cmd.AddCommand(newDeleteCmd(ctx))

	return cmd
}

func newListCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all cached dataset access tokens in profile",
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

			dtokens, err := cli.Tokens().GetDatasetTokensMap(c)
			if err != nil {
				return err
			}

			headers := []string{"Dataset", "Status", "TTL Remaining", "Grants", "Solution"}
			rows := make([][]string, 0, len(dtokens))
			quietIds := make([]string, 0, len(dtokens))

			for ds, dt := range dtokens {
				status := "active"
				remStr := "active"
				if dt.ExpiresAt > 0 {
					rem := time.Until(time.UnixMilli(dt.ExpiresAt))
					if rem <= 0 {
						status = "expired"
						remStr = "expired"
					} else {
						remStr = (rem.Truncate(time.Second)).String()
					}
				}

				gStr := grants.Grant(dt.Grants).Format()
				solStr := dt.Solution
				if solStr == "" {
					solStr = "(all)"
				}

				rows = append(rows, []string{
					ds,
					status,
					remStr,
					gStr,
					solStr,
				})
				quietIds = append(quietIds, ds)
			}

			data := outfmt.Data{
				Headers: headers,
				Rows:    rows,
				QuietId: quietIds,
				JSON:    dtokens,
			}

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
			if err != nil {
				return err
			}

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}
}

func newGetCmd(ctx context.Context) *cobra.Command {
	var (
		datasetStr string
		refreshed  bool
	)

	cmd := &cobra.Command{
		Use:   "get [--dataset <dataset>]",
		Short: "Display dataset access token metadata (safely omits secret token)",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			profileFlag, _ := cmd.Flags().GetString("profile")
			profile := cliconfig.ResolveProfile(profileFlag, c)

			cfg, err := cliconfig.ResolveConfig(profile, c)
			if err != nil {
				return err
			}

			targetDs, err := cliconfig.ResolveDataset(datasetStr, cfg, c)
			if err != nil {
				return err
			}

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			if refreshed {
				if _, err := cli.Tokens().GetValidDatasetAccessToken(c, targetDs); err != nil {
					return exitcode.Wrap(exitcode.Auth, fmt.Errorf("failed to refresh dataset token: %w", err))
				}
			}

			dt, err := cli.Tokens().GetDatasetToken(c, targetDs)
			if err != nil {
				return err
			}
			if dt == nil {
				return exitcode.New(exitcode.NotFound, fmt.Sprintf("no cached access token found for dataset %q in profile %q", targetDs, profile))
			}

			status := "active"
			if cli.Tokens().IsDatasetTokenExpired(dt) {
				status = "expired"
			}

			expStr := "never"
			if dt.ExpiresAt > 0 {
				expStr = time.UnixMilli(dt.ExpiresAt).Format(time.RFC3339)
			}

			gStr := grants.Grant(dt.Grants).Format()
			solStr := dt.Solution
			if solStr == "" {
				solStr = "(all)"
			}

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
			if err != nil {
				return err
			}

			data := outfmt.Data{
				Headers: []string{"Property", "Value"},
				Rows: [][]string{
					{"Profile", profile},
					{"Dataset", targetDs},
					{"Status", status},
					{"Expires At", expStr},
					{"Grants", gStr},
					{"Solution", solStr},
				},
				QuietId: []string{targetDs},
				JSON: map[string]any{
					"profile":    profile,
					"dataset":    targetDs,
					"status":     status,
					"expires_at": dt.ExpiresAt,
					"grants":     gStr,
					"solution":   dt.Solution,
				},
			}

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}

	cmd.Flags().StringVar(&datasetStr, "dataset", "", "Target dataset identifier")
	cmd.Flags().BoolVar(&refreshed, "refreshed", false, "Automatically refresh if missing or expired before displaying")
	return cmd
}

func newExposeCmd(ctx context.Context) *cobra.Command {
	var (
		datasetStr string
		refreshed  bool
	)

	cmd := &cobra.Command{
		Use:   "expose [--dataset <dataset>]",
		Short: "Print the raw Bearer JWT access token to stdout",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			profileFlag, _ := cmd.Flags().GetString("profile")
			profile := cliconfig.ResolveProfile(profileFlag, c)

			cfg, err := cliconfig.ResolveConfig(profile, c)
			if err != nil {
				return err
			}

			targetDs, err := cliconfig.ResolveDataset(datasetStr, cfg, c)
			if err != nil {
				return err
			}

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			if refreshed {
				tok, err := cli.Tokens().GetValidDatasetAccessToken(c, targetDs)
				if err != nil {
					return exitcode.Wrap(exitcode.Auth, fmt.Errorf("failed to get valid dataset access token: %w", err))
				}
				fmt.Fprintln(cmd.OutOrStdout(), tok)
				return nil
			}

			dt, err := cli.Tokens().GetDatasetToken(c, targetDs)
			if err != nil {
				return err
			}
			if dt == nil || dt.AccessToken == "" {
				return exitcode.New(exitcode.NotFound, fmt.Sprintf("no cached access token found for dataset %q in profile %q", targetDs, profile))
			}

			fmt.Fprintln(cmd.OutOrStdout(), dt.AccessToken)
			return nil
		},
	}

	cmd.Flags().StringVar(&datasetStr, "dataset", "", "Target dataset identifier")
	cmd.Flags().BoolVar(&refreshed, "refreshed", false, "Automatically refresh if missing or expired before exposing")
	return cmd
}

func newAttenuationCmd(ctx context.Context) *cobra.Command {
	var (
		datasetStr  string
		grantsStr   string
		solutionStr string
		expSec      int64
	)

	cmd := &cobra.Command{
		Use:   "attenuation [--dataset <dataset>]",
		Short: "Derive an in-flight child access token with reduced permissions from cached token",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			profileFlag, _ := cmd.Flags().GetString("profile")
			profile := cliconfig.ResolveProfile(profileFlag, c)

			cfg, err := cliconfig.ResolveConfig(profile, c)
			if err != nil {
				return err
			}

			targetDs, err := cliconfig.ResolveDataset(datasetStr, cfg, c)
			if err != nil {
				return err
			}

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			parentTok, err := cli.Tokens().GetValidDatasetAccessToken(c, targetDs)
			if err != nil {
				return exitcode.Wrap(exitcode.Auth, fmt.Errorf("failed to obtain parent token for dataset %q: %w", targetDs, err))
			}

			var reqGrants grants.Grant
			if strings.TrimSpace(grantsStr) != "" {
				parsedG, err := grants.Parse(grantsStr)
				if err != nil {
					return exitcode.New(exitcode.Usage, fmt.Sprintf("invalid grants %q: %v", grantsStr, err))
				}
				reqGrants = parsedG
			}

			var reqSol solution.Solution
			if strings.TrimSpace(solutionStr) != "" {
				parsedS, err := solution.ParseOptional(solutionStr)
				if err != nil {
					return exitcode.New(exitcode.Usage, fmt.Sprintf("invalid solution %q: %v", solutionStr, err))
				}
				reqSol = parsedS
			}

			attenApi := &userapidatasettokenaccessattenuate.Api{
				Req: &userapidatasettokenaccessattenuate.Req{
					Headers: &userapidatasettokenaccessattenuate.ReqHeaders{
						Authorization: "Bearer " + parentTok,
					},
					Payload: &userapidatasettokenaccessattenuate.ReqPayload{
						Subsolution: reqSol,
						Subgrants:   reqGrants,
						ExpSeconds:  expSec,
					},
				},
			}

			if _, err := cli.Apis().UserDatasetTokenAccessAttenuate(c, attenApi); err != nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("token attenuation failed on server: %w", err))
			}
			if attenApi.Resp == nil || attenApi.Resp.Payload == nil {
				return fmt.Errorf("empty attenuation response from server")
			}

			fmt.Fprintln(cmd.OutOrStdout(), attenApi.Resp.Payload.AccessToken)
			return nil
		},
	}

	cmd.Flags().StringVar(&datasetStr, "dataset", "", "Target dataset identifier")
	cmd.Flags().StringVar(&grantsStr, "grants", "", "Attenuated grants mask or names (e.g. read,write)")
	cmd.Flags().StringVar(&solutionStr, "solution", "", "Attenuated solution restriction")
	cmd.Flags().Int64Var(&expSec, "exp", 0, "Requested lifetime in seconds (max 24 hours)")

	return cmd
}

func newCheckCmd(ctx context.Context) *cobra.Command {
	var (
		datasetStr string
		refreshed  bool
	)

	cmd := &cobra.Command{
		Use:   "check [--dataset <dataset>]",
		Short: "Cryptographically validate cached token against /api/userz/token/check",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			profileFlag, _ := cmd.Flags().GetString("profile")
			profile := cliconfig.ResolveProfile(profileFlag, c)

			cfg, err := cliconfig.ResolveConfig(profile, c)
			if err != nil {
				return err
			}

			targetDs, err := cliconfig.ResolveDataset(datasetStr, cfg, c)
			if err != nil {
				return err
			}

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			var tokStr string
			if refreshed {
				fresh, err := cli.Tokens().GetValidDatasetAccessToken(c, targetDs)
				if err != nil {
					return exitcode.Wrap(exitcode.Auth, fmt.Errorf("failed to refresh dataset token: %w", err))
				}
				tokStr = fresh
			} else {
				dt, err := cli.Tokens().GetDatasetToken(c, targetDs)
				if err != nil {
					return err
				}
				if dt == nil || dt.AccessToken == "" {
					return exitcode.New(exitcode.NotFound, fmt.Sprintf("no cached token found for dataset %q", targetDs))
				}
				tokStr = dt.AccessToken
			}

			checkApi := &userapitokencheck.Api{
				Req: &userapitokencheck.Req{
					Payload: &userapitokencheck.ReqPayload{
						Token: tokStr,
					},
				},
			}

			if _, err := cli.Apis().UserTokenCheck(c, checkApi); err != nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("token validation failed: %w", err))
			}
			if checkApi.Resp == nil || checkApi.Resp.Payload == nil {
				return fmt.Errorf("empty check response from server")
			}

			p := checkApi.Resp.Payload
			status := "valid"
			if !p.Valid {
				status = p.ErrorCode
				if status == "" {
					status = "invalid"
				}
			}

			remStr := "expired"
			if p.Diagnostics.TimeRemainingSeconds > 0 {
				remStr = (time.Duration(p.Diagnostics.TimeRemainingSeconds) * time.Second).String()
			}

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
			if err != nil {
				return err
			}

			data := outfmt.Data{
				Headers: []string{"Property", "Value"},
				Rows: [][]string{
					{"Dataset", targetDs},
					{"Valid", fmt.Sprintf("%t", p.Valid)},
					{"Status", status},
					{"Summary", p.Summary},
					{"TTL Remaining", remStr},
					{"Grants", grants.Grant(p.Diagnostics.Grants).Format()},
				},
				QuietId: []string{status},
				JSON:    p,
			}

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}

	cmd.Flags().StringVar(&datasetStr, "dataset", "", "Target dataset identifier")
	cmd.Flags().BoolVar(&refreshed, "refreshed", false, "Automatically refresh if missing or expired before checking")
	return cmd
}

func newDeleteCmd(ctx context.Context) *cobra.Command {
	var datasetStr string

	cmd := &cobra.Command{
		Use:   "delete [--dataset <dataset>]",
		Short: "Evict the cached bearer access token for target dataset from profile",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			profileFlag, _ := cmd.Flags().GetString("profile")
			profile := cliconfig.ResolveProfile(profileFlag, c)

			cfg, err := cliconfig.ResolveConfig(profile, c)
			if err != nil {
				return err
			}

			targetDs, err := cliconfig.ResolveDataset(datasetStr, cfg, c)
			if err != nil {
				return err
			}

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			if err := cli.Tokens().RemoveDatasetToken(c, targetDs); err != nil {
				return fmt.Errorf("failed to evict dataset token: %w", err)
			}

			logger.Info(c, "dataset-token-evicted", map[string]any{
				"dataset": targetDs,
				"profile": profile,
			})

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully evicted cached access token for dataset %q from profile %q\n", targetDs, profile)
			return nil
		},
	}

	cmd.Flags().StringVar(&datasetStr, "dataset", "", "Target dataset identifier")
	return cmd
}
