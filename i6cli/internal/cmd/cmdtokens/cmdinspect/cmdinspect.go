package cmdinspect

import (
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapitokencheck"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/infinity6-ai/i6sdk/sdkcorez/reqsign"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// InspectionReport holds the deep diagnostic audit of a single stored token.
type InspectionReport struct {
	CredentialSource string                         `json:"credential_source"`
	Target           string                         `json:"target"`
	TokenKind        string                         `json:"token_kind"`
	Valid            bool                           `json:"valid"`
	Status           string                         `json:"status"`
	Summary          string                         `json:"summary"`
	ExpiresAt        string                         `json:"expires_at,omitempty"`
	TTLRemaining     string                         `json:"ttl_remaining,omitempty"`
	Grants           string                         `json:"grants,omitempty"`
	Hint             string                         `json:"hint,omitempty"`
	Diagnostics      *userapitokencheck.Diagnostics `json:"diagnostics,omitempty"`
	RawPayload       *userapitokencheck.RespPayload `json:"raw_payload,omitempty"`
}

// New constructs the 'inspect' command under 'tokens'.
func New(ctx context.Context) *cobra.Command {
	var refreshed bool

	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Deep inspection and debugging of all credentials stored in active profile",
		Long: `Performs a deep cryptographic and validation audit of every credential currently
stored in the active profile's tokens.db.

Inspects:
  - Active user session access token (cryptographically verified via /api/userz/token/check)
  - Active user refresh token (audited against local session lifecycle state)
  - Active machine robot refresh key assertion (cryptographically proven via /api/userz/token/check)
  - All cached derived dataset access tokens (cryptographically verified via /api/userz/token/check)

With --refreshed, automatically refreshes expired or missing tokens before inspecting.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
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

			var reports []InspectionReport

			// 1. Inspect User Access Token & Refresh Token
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
					expStr := ""
					if p.Diagnostics.ExpiresAt > 0 {
						expStr = time.Unix(p.Diagnostics.ExpiresAt, 0).Format(time.RFC3339)
					}
					status := "valid"
					if !p.Valid {
						status = p.ErrorCode
						if status == "" {
							status = "invalid"
						}
					}

					reports = append(reports, InspectionReport{
						CredentialSource: "user_session",
						Target:           ut.UserId,
						TokenKind:        p.TokenKind,
						Valid:            p.Valid,
						Status:           status,
						Summary:          p.Summary,
						ExpiresAt:        expStr,
						TTLRemaining:     remStr,
						Hint:             p.Diagnostics.Hint,
						Diagnostics:      &p.Diagnostics,
						RawPayload:       p,
					})
				}

				// 1b. Inspect User Refresh Token
				if ut.UserRefreshToken != "" {
					remStr := "unlimited"
					expStr := ""
					isExpired := false
					var expSec int64
					var timeRemSec int64

					if ut.RefreshTokenExpiresAt != nil && *ut.RefreshTokenExpiresAt > 0 {
						expTime := time.UnixMilli(*ut.RefreshTokenExpiresAt)
						expStr = expTime.Format(time.RFC3339)
						expSec = *ut.RefreshTokenExpiresAt / 1000
						remDur := time.Until(expTime)
						timeRemSec = int64(remDur.Seconds())
						if remDur <= 0 {
							isExpired = true
							remStr = "expired"
						} else {
							remStr = remDur.Round(time.Second).String()
						}
					}

					status := "valid"
					summary := "Opaque session refresh token stored locally (enables silent renewal)"
					hint := ""
					if isExpired {
						status = "expired"
						summary = "Session refresh token has expired"
						hint = "Re-authenticate using 'i6 auth login' to establish a new session."
					}

					reports = append(reports, InspectionReport{
						CredentialSource: "user_refresh_token",
						Target:           ut.UserId,
						TokenKind:        "user_refresh_token",
						Valid:            !isExpired,
						Status:           status,
						Summary:          summary,
						ExpiresAt:        expStr,
						TTLRemaining:     remStr,
						Hint:             hint,
						Diagnostics: &userapitokencheck.Diagnostics{
							UserId:               ut.UserId,
							ExpiresAt:            expSec,
							TimeRemainingSeconds: timeRemSec,
							IsExpired:            isExpired,
							Hint:                 hint,
						},
					})
				}
			}

			// 2. Inspect Robot Refresh Keys via assertion
			rkeys, _ := cli.Tokens().GetDatasetRefreshTokensMap(c)
			for ds, rfKey := range rkeys {
				privKey, err := parseEd25519Key(rfKey)
				if err != nil {
					reports = append(reports, InspectionReport{
						CredentialSource: "dataset_refresh_key",
						Target:           ds,
						TokenKind:        "ed25519_key",
						Valid:            false,
						Status:           "malformed_key",
						Summary:          fmt.Sprintf("Failed to parse private key: %v", err),
					})
					continue
				}

				var nonce [16]byte
				exp := uint64(time.Now().Unix() + 300)
				assertion, err := reqsign.SignAssertion(privKey, "POST", "/api/userz/refresh-dataset", "", nil, exp, nonce, false)
				if err != nil {
					reports = append(reports, InspectionReport{
						CredentialSource: "dataset_refresh_key",
						Target:           ds,
						TokenKind:        "ed25519_assertion",
						Valid:            false,
						Status:           "assertion_signing_failed",
						Summary:          fmt.Sprintf("Failed to sign assertion: %v", err),
					})
					continue
				}

				checkApi := &userapitokencheck.Api{
					Req: &userapitokencheck.Req{
						Payload: &userapitokencheck.ReqPayload{
							Token: assertion.Pack(),
						},
					},
				}

				if _, err := cli.Apis().UserTokenCheck(c, checkApi); err == nil && checkApi.Resp != nil && checkApi.Resp.Payload != nil {
					p := checkApi.Resp.Payload
					status := "valid"
					if !p.Valid {
						status = p.ErrorCode
						if status == "" {
							status = "invalid"
						}
					}
					reports = append(reports, InspectionReport{
						CredentialSource: "dataset_refresh_key",
						Target:           ds,
						TokenKind:        p.TokenKind,
						Valid:            p.Valid,
						Status:           status,
						Summary:          p.Summary,
						TTLRemaining:     "unlimited",
						Hint:             p.Diagnostics.Hint,
						Diagnostics:      &p.Diagnostics,
						RawPayload:       p,
					})
				} else {
					reports = append(reports, InspectionReport{
						CredentialSource: "dataset_refresh_key",
						Target:           ds,
						TokenKind:        "ed25519_key",
						Valid:            true,
						Status:           "configured",
						Summary:          "Valid 32-byte Ed25519 seed configured",
					})
				}
			}

			// 3. Inspect Cached Dataset Access Tokens
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
					expStr := ""
					if p.Diagnostics.ExpiresAt > 0 {
						expStr = time.Unix(p.Diagnostics.ExpiresAt, 0).Format(time.RFC3339)
					}
					status := "valid"
					if !p.Valid {
						status = p.ErrorCode
						if status == "" {
							status = "invalid"
						}
					}

					grantsStr := grants.Grant(p.Diagnostics.Grants).Format()

					reports = append(reports, InspectionReport{
						CredentialSource: "dataset_access_token",
						Target:           ds,
						TokenKind:        p.TokenKind,
						Valid:            p.Valid,
						Status:           status,
						Summary:          p.Summary,
						ExpiresAt:        expStr,
						TTLRemaining:     remStr,
						Grants:           grantsStr,
						Hint:             p.Diagnostics.Hint,
						Diagnostics:      &p.Diagnostics,
						RawPayload:       p,
					})
				}
			}

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
			if err != nil {
				return err
			}

			headers := []string{"Source", "Target", "Kind", "Valid", "Status", "TTL Remaining", "Grants", "Summary"}
			rows := make([][]string, 0, len(reports))
			quietIds := make([]string, 0, len(reports))

			for _, r := range reports {
				validStr := "no"
				if r.Valid {
					validStr = "yes"
				}
				remStr := r.TTLRemaining
				if remStr == "" {
					remStr = "-"
				}
				gStr := r.Grants
				if gStr == "" {
					gStr = "-"
				}

				rows = append(rows, []string{
					r.CredentialSource,
					r.Target,
					r.TokenKind,
					validStr,
					r.Status,
					remStr,
					gStr,
					r.Summary,
				})
				quietIds = append(quietIds, fmt.Sprintf("%s:%s", r.Target, r.Status))
			}

			data := outfmt.Data{
				Headers: headers,
				Rows:    rows,
				QuietId: quietIds,
				JSON:    reports,
			}

			logger.Debug(c, "tokens-inspect", map[string]any{
				"profile": profile,
				"count":   len(reports),
			})

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}

	cmd.Flags().BoolVar(&refreshed, "refreshed", false, "Automatically refresh expired tokens before inspecting")
	return cmd
}

func parseEd25519Key(keyInput string) (ed25519.PrivateKey, error) {
	keyInput = strings.TrimSpace(keyInput)
	rawBytes, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(keyInput, "="))
	if err != nil {
		rawBytes, err = base64.StdEncoding.DecodeString(keyInput)
		if err != nil {
			rawBytes, err = hex.DecodeString(keyInput)
			if err != nil {
				rawBytes = []byte(keyInput)
			}
		}
	}
	if len(rawBytes) == ed25519.PrivateKeySize {
		return ed25519.PrivateKey(rawBytes), nil
	}
	if len(rawBytes) == ed25519.SeedSize {
		return ed25519.NewKeyFromSeed(rawBytes), nil
	}
	if parsed, err := x509.ParsePKCS8PrivateKey(rawBytes); err == nil {
		if edKey, ok := parsed.(ed25519.PrivateKey); ok {
			return edKey, nil
		}
	}
	return nil, fmt.Errorf("invalid ed25519 key length: expected 32 (seed) or 64 bytes, got %d", len(rawBytes))
}
