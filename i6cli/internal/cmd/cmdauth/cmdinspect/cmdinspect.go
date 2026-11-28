package cmdinspect

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapitokencheck"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'inspect' command under 'auth'.
func New(ctx context.Context) *cobra.Command {
	var refreshed bool

	cmd := &cobra.Command{
		Use:   "inspect [<token-string | ->]",
		Short: "Perform deep cryptographic and identity diagnostics on a user token",
		Long: `Performs a deep cryptographic and identity audit of a User JWT access token
against the user service diagnostic mirror (/api/userz/token/check).

If <token-string> (or '-' for stdin) is provided, audits that arbitrary token statelessly
without modifying or reading tokens.db.

If omitted, audits the human user access token stored in the active configuration profile.
With --refreshed, automatically refreshes an expired session before inspecting.`,
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

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			var tokenStr string
			if len(args) > 0 {
				arg := strings.TrimSpace(args[0])
				if arg == "-" {
					data, err := io.ReadAll(cmd.InOrStdin())
					if err != nil {
						return fmt.Errorf("failed to read token from stdin: %w", err)
					}
					tokenStr = strings.TrimSpace(string(data))
				} else {
					tokenStr = arg
				}
			}

			if tokenStr == "" {
				if refreshed {
					if fresh, err := cli.Tokens().GetValidAccessToken(c); err == nil && fresh != "" {
						tokenStr = fresh
					}
				}
				if tokenStr == "" {
					ut, err := cli.Tokens().GetUserToken(c)
					if err != nil || ut == nil || ut.UserAccessToken == "" {
						return exitcode.New(
							exitcode.Auth,
							fmt.Sprintf("no active user session found in profile %q: provide a token argument or authenticate via 'i6 auth login'", profile),
						)
					}
					tokenStr = ut.UserAccessToken
				}
			}

			checkApi := &userapitokencheck.Api{
				Req: &userapitokencheck.Req{
					Payload: &userapitokencheck.ReqPayload{
						Token: tokenStr,
					},
				},
			}

			status, err := cli.Apis().UserTokenCheck(c, checkApi)
			if err != nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("failed to verify user token: %w", err))
			}
			if status >= 400 || checkApi.Resp == nil || checkApi.Resp.Payload == nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("server returned error %d checking user token", status))
			}

			p := checkApi.Resp.Payload

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
			if err != nil {
				return err
			}

			remStr := "expired"
			if p.Diagnostics.TimeRemainingSeconds > 0 {
				remStr = (time.Duration(p.Diagnostics.TimeRemainingSeconds) * time.Second).String()
			}
			expStr := ""
			if p.Diagnostics.ExpiresAt > 0 {
				expStr = time.Unix(p.Diagnostics.ExpiresAt, 0).Format(time.RFC3339)
			}
			statusStr := "valid"
			if !p.Valid {
				statusStr = p.ErrorCode
				if statusStr == "" {
					statusStr = "invalid"
				}
			}

			userId := p.Diagnostics.UserId
			if userId == "" {
				userId = "(unknown)"
			}

			validStr := "no"
			if p.Valid {
				validStr = "yes"
			}

			rows := [][]string{
				{"Valid", validStr},
				{"Token Kind", p.TokenKind},
				{"User ID", userId},
				{"Status", statusStr},
				{"Summary", p.Summary},
			}

			if expStr != "" {
				rows = append(rows, []string{"Expires At", expStr})
			}
			rows = append(rows, []string{"TTL Remaining", remStr})

			if p.Diagnostics.Hint != "" {
				rows = append(rows, []string{"Hint", p.Diagnostics.Hint})
			}

			data := outfmt.Data{
				Headers: []string{"Property", "Value"},
				Rows:    rows,
				QuietId: []string{fmt.Sprintf("%s:%s", userId, statusStr)},
				Raw:     tokenStr,
				JSON:    p,
			}

			logger.Debug(c, "auth-inspect", map[string]any{
				"valid":      p.Valid,
				"token_kind": p.TokenKind,
				"user_id":    userId,
			})

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}

	cmd.Flags().BoolVar(&refreshed, "refreshed", false, "Automatically refresh expired user session before inspecting")

	return cmd
}
