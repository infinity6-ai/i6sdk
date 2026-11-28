package cmdraw

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/gox/commonz/pathz"
	"github.com/infinity6-ai/gox/commonz/validation/checked"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenaccessattenuate"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefresh"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapitokencheck"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/infinity6-ai/i6sdk/sdkcorez/reqsign"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'raw' command under 'tokens'.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "raw",
		Short: "Stateless zero-disk cryptographic and token exchange operations",
		Long: `Performs pure in-memory cryptographic and token exchange operations.
These commands never read from or write to the local tokens.db store.`,
	}

	cmd.AddCommand(newCheckCmd(ctx))
	cmd.AddCommand(newDatasetBearerCmd(ctx))
	cmd.AddCommand(newAttenuationCmd(ctx))
	cmd.AddCommand(newSignRequestCmd(ctx))

	return cmd
}

func newCheckCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "check <token-string | ->",
		Short: "Statelessly inspect and validate an arbitrary token string",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			var tokenStr string
			if args[0] == "-" {
				data, err := io.ReadAll(cmd.InOrStdin())
				if err != nil {
					return fmt.Errorf("failed to read token from stdin: %w", err)
				}
				tokenStr = strings.TrimSpace(string(data))
			} else {
				tokenStr = strings.TrimSpace(args[0])
			}

			if tokenStr == "" {
				return exitcode.New(exitcode.Usage, "token string cannot be empty")
			}

			profileFlag, _ := cmd.Flags().GetString("profile")
			profile := cliconfig.ResolveProfile(profileFlag, c)

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			checkApi := &userapitokencheck.Api{
				Req: &userapitokencheck.Req{
					Payload: &userapitokencheck.ReqPayload{
						Token: tokenStr,
					},
				},
			}

			if _, err := cli.Apis().UserTokenCheck(c, checkApi); err != nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("token check call failed: %w", err))
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

			target := p.Diagnostics.UserId
			if target == "" {
				target = p.Diagnostics.Dataset
			}
			if target == "" {
				target = "(none)"
			}

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
			if err != nil {
				return err
			}

			data := outfmt.Data{
				Headers: []string{"Property", "Value"},
				Rows: [][]string{
					{"Valid", fmt.Sprintf("%t", p.Valid)},
					{"Token Kind", p.TokenKind},
					{"Target", target},
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
}

func newDatasetBearerCmd(ctx context.Context) *cobra.Command {
	var (
		rfKey   string
		dataset string
		expSec  int64
	)

	cmd := &cobra.Command{
		Use:   "dataset-bearer --refresh-token <key|-> --dataset <dataset>",
		Short: "Statelessly exchange an Ed25519 robot key for a short-lived Bearer token",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			if strings.TrimSpace(dataset) == "" {
				return exitcode.New(exitcode.Usage, "--dataset is required")
			}

			var keyStr string
			if rfKey == "-" {
				data, err := io.ReadAll(cmd.InOrStdin())
				if err != nil {
					return fmt.Errorf("failed to read key from stdin: %w", err)
				}
				keyStr = strings.TrimSpace(string(data))
			} else {
				keyStr = strings.TrimSpace(rfKey)
			}

			if keyStr == "" {
				return exitcode.New(exitcode.Usage, "--refresh-token is required")
			}

			privKey, err := parseEd25519Key(keyStr)
			if err != nil {
				return exitcode.Wrap(exitcode.Usage, fmt.Errorf("invalid ed25519 private key: %w", err))
			}

			profileFlag, _ := cmd.Flags().GetString("profile")
			profile := cliconfig.ResolveProfile(profileFlag, c)

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			if expSec <= 0 {
				expSec = 300
			}

			var nonce [16]byte
			exp := uint64(time.Now().Unix() + expSec)
			assertion, err := reqsign.SignAssertion(privKey, "POST", "/api/userz/refresh-dataset", "", nil, exp, nonce, false)
			if err != nil {
				return fmt.Errorf("failed to sign refresh assertion: %w", err)
			}

			var ds lakedsz.Dataset
			if err := checked.Set(&ds, dataset); err != nil {
				return exitcode.New(exitcode.Usage, fmt.Sprintf("invalid dataset identifier %q: %v", dataset, err))
			}

			refApi := &userapidatasettokenrefresh.Api{
				Req: &userapidatasettokenrefresh.Req{
					Headers: &userapidatasettokenrefresh.ReqHeaders{
						SignatureHeader: assertion.Pack(),
					},
					Payload: &userapidatasettokenrefresh.ReqPayload{
						Dataset:    ds,
						ExpSeconds: expSec,
						Unbound:    true,
					},
				},
			}

			if _, err := cli.Apis().UserDatasetTokenRefresh(c, refApi); err != nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("token exchange failed: %w", err))
			}
			if refApi.Resp == nil || refApi.Resp.Payload == nil {
				return fmt.Errorf("empty refresh response from server")
			}

			fmt.Fprintln(cmd.OutOrStdout(), refApi.Resp.Payload.AccessToken)
			return nil
		},
	}

	cmd.Flags().StringVar(&rfKey, "refresh-token", "", "Ed25519 dataset refresh private key (or '-' for stdin)")
	cmd.Flags().StringVar(&dataset, "dataset", "", "Target dataset identifier")
	cmd.Flags().Int64Var(&expSec, "exp", 300, "Requested lifetime in seconds")

	return cmd
}

func newAttenuationCmd(ctx context.Context) *cobra.Command {
	var (
		grantsStr   string
		solutionStr string
		expSec      int64
	)

	cmd := &cobra.Command{
		Use:   "attenuation [<token|->]",
		Short: "Statelessly derive an attenuated child token from a parent access token",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			var parentTok string
			if len(args) == 1 && args[0] != "-" {
				parentTok = strings.TrimSpace(args[0])
			} else {
				data, err := io.ReadAll(cmd.InOrStdin())
				if err != nil {
					return fmt.Errorf("failed to read parent token from stdin: %w", err)
				}
				parentTok = strings.TrimSpace(string(data))
			}

			if parentTok == "" {
				return exitcode.New(exitcode.Usage, "parent token cannot be empty")
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

			profileFlag, _ := cmd.Flags().GetString("profile")
			profile := cliconfig.ResolveProfile(profileFlag, c)

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
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
				return exitcode.Wrap(exitcode.General, fmt.Errorf("token attenuation failed: %w", err))
			}
			if attenApi.Resp == nil || attenApi.Resp.Payload == nil {
				return fmt.Errorf("empty attenuation response from server")
			}

			fmt.Fprintln(cmd.OutOrStdout(), attenApi.Resp.Payload.AccessToken)
			return nil
		},
	}

	cmd.Flags().StringVar(&grantsStr, "grants", "", "Attenuated grants mask or names (e.g. read,write)")
	cmd.Flags().StringVar(&solutionStr, "solution", "", "Attenuated solution restriction")
	cmd.Flags().Int64Var(&expSec, "exp", 0, "Requested lifetime in seconds")

	return cmd
}

func newSignRequestCmd(ctx context.Context) *cobra.Command {
	var (
		keyStr   string
		tokenStr string
		method   string
		path     string
		body     string
	)

	cmd := &cobra.Command{
		Use:   "sign-request --key <key|-> --token <token|-> --method <METHOD> --path <PATH> [--body <str|->]",
		Short: "Statelessly generate X-I6-Signature header for a request",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(method) == "" || strings.TrimSpace(path) == "" {
				return exitcode.New(exitcode.Usage, "--method and --path are required")
			}

			var finalKey string
			if keyStr == "-" {
				data, err := io.ReadAll(cmd.InOrStdin())
				if err != nil {
					return fmt.Errorf("failed to read key from stdin: %w", err)
				}
				finalKey = strings.TrimSpace(string(data))
			} else {
				finalKey = strings.TrimSpace(keyStr)
			}
			if finalKey == "" {
				return exitcode.New(exitcode.Usage, "--key is required")
			}

			privKey, err := parseEd25519Key(finalKey)
			if err != nil {
				return exitcode.Wrap(exitcode.Usage, fmt.Errorf("invalid ed25519 key: %w", err))
			}

			var finalTok string
			if tokenStr == "-" {
				data, err := io.ReadAll(cmd.InOrStdin())
				if err != nil {
					return fmt.Errorf("failed to read token from stdin: %w", err)
				}
				finalTok = strings.TrimSpace(string(data))
			} else {
				finalTok = strings.TrimSpace(tokenStr)
			}
			if finalTok == "" {
				return exitcode.New(exitcode.Usage, "--token is required for request signing")
			}

			parsedPath, err := pathz.Parse(strings.TrimSpace(path))
			if err != nil {
				return exitcode.Wrap(exitcode.Usage, fmt.Errorf("invalid path %q: %w", path, err))
			}

			sigData := reqsign.RequestSignatureData{
				AccessToken: finalTok,
				Method:      strings.ToUpper(strings.TrimSpace(method)),
				Path:        parsedPath,
			}

			opts := reqsign.RequestSignatureOptions{
				PrivKey: privKey,
			}

			if body != "" && body != "-" {
				hash := sha256.Sum256([]byte(body))
				sigData.BodyHash = hash[:]
			} else {
				opts.UnsignedBody = true
			}

			sigHeader, err := reqsign.BuildRequestSignature(sigData, opts)
			if err != nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("failed to build request signature: %w", err))
			}

			fmt.Fprintln(cmd.OutOrStdout(), sigHeader)
			return nil
		},
	}

	cmd.Flags().StringVar(&keyStr, "key", "", "Ed25519 private key string (or '-' for stdin)")
	cmd.Flags().StringVar(&tokenStr, "token", "", "Bound dataset access token (or '-' for stdin)")
	cmd.Flags().StringVar(&method, "method", "", "HTTP method (GET, POST, etc.)")
	cmd.Flags().StringVar(&path, "path", "", "HTTP request path")
	cmd.Flags().StringVar(&body, "body", "", "Request body string")

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
	return nil, fmt.Errorf("invalid ed25519 key length: expected 32 (seed) or 64 bytes, got %d", len(rawBytes))
}
