package cmduser

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenaccessattenuate"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthlogout"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'user' command under 'tokens'.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage the authenticated human user session root credential",
		Long:  `Controls the human user session credentials stored in the active profile's tokens.db.`,
	}

	cmd.AddCommand(newGetCmd(ctx))
	cmd.AddCommand(newExposeCmd(ctx))
	cmd.AddCommand(newExportCmd(ctx))
	cmd.AddCommand(newAttenuationCmd(ctx))
	cmd.AddCommand(newDeleteCmd(ctx))

	return cmd
}

func newGetCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Display user session metadata (safely omits secret tokens)",
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

			ut, err := cli.Tokens().GetUserToken(c)
			if err != nil {
				return err
			}
			if ut == nil {
				return exitcode.New(exitcode.Auth, fmt.Sprintf("no user session found in profile %q", profile))
			}

			status := "active"
			if cli.Tokens().IsAccessTokenExpired(ut) {
				status = "access token expired"
			}

			accessExp := "never"
			if ut.ExpiresAt > 0 {
				accessExp = time.UnixMilli(ut.ExpiresAt).Format(time.RFC3339)
			}

			refreshExp := "never"
			if ut.RefreshTokenExpiresAt != nil && *ut.RefreshTokenExpiresAt > 0 {
				refreshExp = time.UnixMilli(*ut.RefreshTokenExpiresAt).Format(time.RFC3339)
			}

			hasRefresh := "no"
			if ut.UserRefreshToken != "" {
				hasRefresh = "yes"
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
					{"User ID", ut.UserId},
					{"Session Status", status},
					{"Access Token Expiry", accessExp},
					{"Refresh Token Expiry", refreshExp},
					{"Refreshable", hasRefresh},
				},
				QuietId: []string{ut.UserId},
				JSON: map[string]any{
					"profile":               profile,
					"user_id":               ut.UserId,
					"status":                status,
					"access_token_expiry":   accessExp,
					"refresh_token_expiry":  refreshExp,
					"refreshable":           ut.UserRefreshToken != "",
				},
			}

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}
}

func newExposeCmd(ctx context.Context) *cobra.Command {
	var refreshed bool

	cmd := &cobra.Command{
		Use:   "expose",
		Short: "Print the raw active JWT access token to stdout",
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

			if refreshed {
				tok, err := cli.Tokens().GetValidAccessToken(c)
				if err != nil {
					return exitcode.Wrap(exitcode.Auth, fmt.Errorf("failed to get valid access token: %w", err))
				}
				fmt.Fprintln(cmd.OutOrStdout(), tok)
				return nil
			}

			ut, err := cli.Tokens().GetUserToken(c)
			if err != nil {
				return err
			}
			if ut == nil || ut.UserAccessToken == "" {
				return exitcode.New(exitcode.Auth, fmt.Sprintf("no active user access token found in profile %q", profile))
			}

			fmt.Fprintln(cmd.OutOrStdout(), ut.UserAccessToken)
			return nil
		},
	}

	cmd.Flags().BoolVar(&refreshed, "refreshed", false, "Automatically refresh the token if expired before exposing")
	return cmd
}

func newExportCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "export",
		Short: "Print the active user refresh token to stdout",
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

			ut, err := cli.Tokens().GetUserToken(c)
			if err != nil {
				return err
			}
			if ut == nil || ut.UserRefreshToken == "" {
				return exitcode.New(exitcode.Auth, fmt.Sprintf("no user refresh token found in profile %q", profile))
			}

			fmt.Fprintln(cmd.OutOrStdout(), ut.UserRefreshToken)
			return nil
		},
	}
}

func newAttenuationCmd(ctx context.Context) *cobra.Command {
	var (
		datasetStr  string
		grantsStr   string
		solutionStr string
		expSec      int64
	)

	cmd := &cobra.Command{
		Use:   "attenuation --dataset <dataset>",
		Short: "Mint an in-flight attenuated child access token for target dataset",
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

			// 1. Obtain parent dataset access token for target dataset
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

			childTok := attenApi.Resp.Payload.AccessToken
			fmt.Fprintln(cmd.OutOrStdout(), childTok)
			return nil
		},
	}

	cmd.Flags().StringVar(&datasetStr, "dataset", "", "Target dataset identifier")
	cmd.Flags().StringVar(&grantsStr, "grants", "", "Attenuated grants mask or names (e.g. read,write)")
	cmd.Flags().StringVar(&solutionStr, "solution", "", "Attenuated solution restriction")
	cmd.Flags().Int64Var(&expSec, "exp", 0, "Requested lifetime in seconds (max 24 hours)")

	return cmd
}

func newDeleteCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Revoke session on server and delete user credentials from profile",
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

			ut, _ := cli.Tokens().GetUserToken(c)
			if ut == nil {
				fmt.Fprintf(cmd.OutOrStdout(), "No active user session in profile %q\n", profile)
				return nil
			}

			if ut.UserRefreshToken != "" {
				logoutApi := &userapiauthlogout.Api{
					Req: &userapiauthlogout.Req{
						Payload: &userapiauthlogout.ReqPayload{
							RefreshToken: ut.UserRefreshToken,
						},
					},
				}
				_, _ = cli.Apis().UserAuthLogout(c, logoutApi)
			}

			if err := cli.Tokens().RemoveUserToken(c); err != nil {
				return fmt.Errorf("failed to remove user token: %w", err)
			}

			logger.Info(c, "user-session-deleted", map[string]any{
				"profile": profile,
				"user_id": ut.UserId,
			})

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully deleted user session for %s from profile %q\n", ut.UserId, profile)
			return nil
		},
	}
}
