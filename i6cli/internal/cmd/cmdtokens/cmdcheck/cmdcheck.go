package cmdcheck

import (
	"context"
	"time"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapitokencheck"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// CheckItem holds the audit result for a single credential.
type CheckItem struct {
	Type          string `json:"type"`
	Target        string `json:"target"`
	Valid         bool   `json:"valid"`
	Status        string `json:"status"`
	TimeRemaining string `json:"time_remaining"`
	Grants        string `json:"grants,omitempty"`
	Hint          string `json:"hint,omitempty"`
}

// New constructs the 'check' command under 'tokens'.
func New(ctx context.Context) *cobra.Command {
	var refreshed bool

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Perform a diagnostic audit of all credentials in profile",
		Long: `Audits all credentials stored in the active profile's tokens.db against
the server diagnostic validation endpoint (/api/userz/token/check).

Checks:
  - Active root credential (User Session or Robot Refresh Key)
  - All cached derived dataset access tokens

With --refreshed, automatically refreshes expired or missing tokens before auditing.`,
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

			var items []CheckItem

			// 1. Audit User Root Credential
			ut, _ := cli.Tokens().GetUserToken(c)
			if ut != nil {
				tokStr := ut.UserAccessToken
				if refreshed {
					if fresh, err := cli.Tokens().GetValidAccessToken(c); err == nil && fresh != "" {
						tokStr = fresh
					}
				}

				checkApi := &userapitokencheck.Api{
					Req: &userapitokencheck.Req{
						Payload: &userapitokencheck.ReqPayload{
							Token: tokStr,
						},
					},
				}
				if _, err := cli.Apis().UserTokenCheck(c, checkApi); err == nil && checkApi.Resp != nil && checkApi.Resp.Payload != nil {
					p := checkApi.Resp.Payload
					remStr := "expired"
					if p.Diagnostics.TimeRemainingSeconds > 0 {
						remStr = (time.Duration(p.Diagnostics.TimeRemainingSeconds) * time.Second).String()
					}
					status := "valid"
					if !p.Valid {
						status = p.ErrorCode
						if status == "" {
							status = "invalid"
						}
					}
					items = append(items, CheckItem{
						Type:          "User Access Token",
						Target:        ut.UserId,
						Valid:         p.Valid,
						Status:        status,
						TimeRemaining: remStr,
						Hint:          p.Diagnostics.Hint,
					})
				} else {
					status := "active"
					if cli.Tokens().IsAccessTokenExpired(ut) {
						status = "expired"
					}
					items = append(items, CheckItem{
						Type:   "User Access Token (local)",
						Target: ut.UserId,
						Valid:  status == "active",
						Status: status,
					})
				}
			}

			// 2. Audit Robot Keys
			rkeys, _ := cli.Tokens().GetDatasetRefreshTokensMap(c)
			for ds := range rkeys {
				items = append(items, CheckItem{
					Type:          "Ed25519 Dataset Refresh Key",
					Target:        ds,
					Valid:         true,
					Status:        "configured",
					TimeRemaining: "unlimited",
				})
			}

			// 3. Audit Cached Dataset Access Tokens
			dtokens, _ := cli.Tokens().GetDatasetTokensMap(c)
			for ds, dt := range dtokens {
				tokStr := dt.AccessToken
				if refreshed {
					if fresh, err := cli.Tokens().GetValidDatasetAccessToken(c, ds); err == nil && fresh != "" {
						tokStr = fresh
					}
				}

				checkApi := &userapitokencheck.Api{
					Req: &userapitokencheck.Req{
						Payload: &userapitokencheck.ReqPayload{
							Token: tokStr,
						},
					},
				}
				if _, err := cli.Apis().UserTokenCheck(c, checkApi); err == nil && checkApi.Resp != nil && checkApi.Resp.Payload != nil {
					p := checkApi.Resp.Payload
					remStr := "expired"
					if p.Diagnostics.TimeRemainingSeconds > 0 {
						remStr = (time.Duration(p.Diagnostics.TimeRemainingSeconds) * time.Second).String()
					}
					status := "valid"
					if !p.Valid {
						status = p.ErrorCode
						if status == "" {
							status = "invalid"
						}
					}
					grantsStr := grants.Grant(p.Diagnostics.Grants).Format()

					items = append(items, CheckItem{
						Type:          "Dataset Access Token",
						Target:        ds,
						Valid:         p.Valid,
						Status:        status,
						TimeRemaining: remStr,
						Grants:        grantsStr,
						Hint:          p.Diagnostics.Hint,
					})
				} else {
					status := "active"
					if cli.Tokens().IsDatasetTokenExpired(dt) {
						status = "expired"
					}
					items = append(items, CheckItem{
						Type:   "Dataset Access Token (local)",
						Target: ds,
						Valid:  status == "active",
						Status: status,
					})
				}
			}

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
			if err != nil {
				return err
			}

			headers := []string{"Type", "Target", "Valid", "Status", "TTL Remaining", "Grants"}
			rows := make([][]string, 0, len(items))
			quietIds := make([]string, 0, len(items))

			for _, it := range items {
				validStr := "no"
				if it.Valid {
					validStr = "yes"
				}
				gStr := it.Grants
				if gStr == "" {
					gStr = "-"
				}
				remStr := it.TimeRemaining
				if remStr == "" {
					remStr = "-"
				}
				rows = append(rows, []string{
					it.Type,
					it.Target,
					validStr,
					it.Status,
					remStr,
					gStr,
				})
				quietIds = append(quietIds, it.Target)
			}

			data := outfmt.Data{
				Headers: headers,
				Rows:    rows,
				QuietId: quietIds,
				JSON:    items,
			}

			logger.Debug(c, "tokens-check", map[string]any{
				"profile": profile,
				"count":   len(items),
			})

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}

	cmd.Flags().BoolVar(&refreshed, "refreshed", false, "Automatically refresh expired tokens before auditing")
	return cmd
}
