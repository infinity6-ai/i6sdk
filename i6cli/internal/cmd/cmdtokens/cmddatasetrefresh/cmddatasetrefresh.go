package cmddatasetrefresh

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/gox/commonz/pathz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenaccessattenuate"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/infinity6-ai/i6sdk/sdkcorez/reqsign"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'dataset-refresh' command under 'tokens'.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dataset-refresh",
		Aliases: []string{"robot"},
		Short:   "Manage the machine service account Ed25519 refresh key root credential",
		Long: `Controls the machine robot Ed25519 refresh key stored in the active profile's tokens.db.

Root Credential Invariant: Mutually exclusive with user sessions. A profile holds either
a user session OR a machine robot key, never both.

Security Invariant: 'get' omits sensitive private keys. Use 'export' to view the raw key.`,
	}

	cmd.AddCommand(newGetCmd(ctx))
	cmd.AddCommand(newExportCmd(ctx))
	cmd.AddCommand(newSetCmd(ctx))
	cmd.AddCommand(newAttenuationCmd(ctx))
	cmd.AddCommand(newDeleteCmd(ctx))
	cmd.AddCommand(newBoundedCmd(ctx))

	return cmd
}

func newGetCmd(ctx context.Context) *cobra.Command {
	var datasetStr string

	cmd := &cobra.Command{
		Use:   "get [--dataset <dataset>]",
		Short: "Display robot key metadata (safely omits the private key)",
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

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			// If dataset not passed, check if there's an active key
			targetDs := strings.TrimSpace(datasetStr)
			if targetDs == "" {
				rkeys, err := cli.Tokens().GetDatasetRefreshTokensMap(c)
				if err == nil && len(rkeys) == 1 {
					for ds := range rkeys {
						targetDs = ds
					}
				}
			}
			if targetDs == "" {
				ds, err := cliconfig.ResolveDataset(datasetStr, cfg, c)
				if err != nil {
					return err
				}
				targetDs = ds
			}

			rfToken, err := cli.Tokens().GetDatasetRefreshToken(c, targetDs)
			if err != nil {
				return err
			}
			if strings.TrimSpace(rfToken) == "" {
				return exitcode.New(exitcode.NotFound, fmt.Sprintf("no robot refresh key found for dataset %q in profile %q", targetDs, profile))
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
					{"Key Type", "Ed25519 Dataset Refresh Key"},
					{"Status", "configured"},
				},
				QuietId: []string{targetDs},
				JSON: map[string]any{
					"profile":  profile,
					"dataset":  targetDs,
					"type":     "dataset_refresh",
					"status":   "configured",
				},
			}

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}

	cmd.Flags().StringVar(&datasetStr, "dataset", "", "Target dataset identifier")
	return cmd
}

func newExportCmd(ctx context.Context) *cobra.Command {
	var datasetStr string

	cmd := &cobra.Command{
		Use:   "export [--dataset <dataset>]",
		Short: "Print the active Ed25519 dataset refresh key to stdout",
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

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			targetDs := strings.TrimSpace(datasetStr)
			if targetDs == "" {
				rkeys, err := cli.Tokens().GetDatasetRefreshTokensMap(c)
				if err == nil && len(rkeys) == 1 {
					for ds := range rkeys {
						targetDs = ds
					}
				}
			}
			if targetDs == "" {
				ds, err := cliconfig.ResolveDataset(datasetStr, cfg, c)
				if err != nil {
					return err
				}
				targetDs = ds
			}

			rfToken, err := cli.Tokens().GetDatasetRefreshToken(c, targetDs)
			if err != nil {
				return err
			}
			if strings.TrimSpace(rfToken) == "" {
				return exitcode.New(exitcode.NotFound, fmt.Sprintf("no robot refresh key found for dataset %q in profile %q", targetDs, profile))
			}

			fmt.Fprintln(cmd.OutOrStdout(), rfToken)
			return nil
		},
	}

	cmd.Flags().StringVar(&datasetStr, "dataset", "", "Target dataset identifier")
	return cmd
}

func newSetCmd(ctx context.Context) *cobra.Command {
	var datasetStr string

	cmd := &cobra.Command{
		Use:   "set <key|-> [--dataset <dataset>]",
		Short: "Set or replace the robot refresh key for this profile",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ExactArgs(1),
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

			// Invariant check: ensure no active human user session exists in this profile
			ut, _ := cli.Tokens().GetUserToken(c)
			if ut != nil && ut.UserAccessToken != "" {
				return exitcode.New(
					exitcode.Usage,
					fmt.Sprintf("active user session exists for user %q in profile %q: a profile cannot hold both a user session and a robot key. Please use a separate profile (e.g. 'i6 --profile worker tokens dataset-refresh set ...')", ut.UserId, profile),
				)
			}

			cfg, err := cliconfig.ResolveConfig(profile, c)
			if err != nil {
				return err
			}

			targetDs, err := cliconfig.ResolveDataset(datasetStr, cfg, c)
			if err != nil {
				return err
			}

			var keyStr string
			if args[0] == "-" {
				data, err := io.ReadAll(cmd.InOrStdin())
				if err != nil {
					return fmt.Errorf("failed to read key from stdin: %w", err)
				}
				keyStr = strings.TrimSpace(string(data))
			} else {
				keyStr = strings.TrimSpace(args[0])
			}

			if keyStr == "" {
				return exitcode.New(exitcode.Usage, "robot key string cannot be empty")
			}

			if err := cli.Tokens().SetDatasetRefreshToken(c, targetDs, keyStr); err != nil {
				return fmt.Errorf("failed to store robot refresh key: %w", err)
			}

			logger.Info(c, "robot-key-set", map[string]any{
				"dataset": targetDs,
				"profile": profile,
			})

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully saved robot key for dataset %q in profile %q\n", targetDs, profile)
			return nil
		},
	}

	cmd.Flags().StringVar(&datasetStr, "dataset", "", "Target dataset identifier")
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
		Short: "Derive an in-flight attenuated access token using the stored robot key",
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

			parentTok, err := cli.Tokens().RefreshDatasetAccessToken(c, targetDs)
			if err != nil {
				return exitcode.Wrap(exitcode.Auth, fmt.Errorf("failed to obtain parent token using robot key: %w", err))
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
						Authorization: "Bearer " + parentTok.AccessToken,
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

func newDeleteCmd(ctx context.Context) *cobra.Command {
	var datasetStr string

	cmd := &cobra.Command{
		Use:   "delete [--dataset <dataset>]",
		Short: "Remove robot refresh key from profile",
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

			cli, err := cliconfig.NewClient(profile, c)
			if err != nil {
				return err
			}

			targetDs := strings.TrimSpace(datasetStr)
			if targetDs == "" {
				rkeys, err := cli.Tokens().GetDatasetRefreshTokensMap(c)
				if err == nil && len(rkeys) == 1 {
					for ds := range rkeys {
						targetDs = ds
					}
				}
			}
			if targetDs == "" {
				ds, err := cliconfig.ResolveDataset(datasetStr, cfg, c)
				if err != nil {
					return err
				}
				targetDs = ds
			}

			if err := cli.Tokens().RemoveDatasetRefreshToken(c, targetDs); err != nil {
				return fmt.Errorf("failed to remove robot key: %w", err)
			}

			logger.Info(c, "robot-key-deleted", map[string]any{
				"dataset": targetDs,
				"profile": profile,
			})

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully deleted robot key for dataset %q from profile %q\n", targetDs, profile)
			return nil
		},
	}

	cmd.Flags().StringVar(&datasetStr, "dataset", "", "Target dataset identifier")
	return cmd
}

func newBoundedCmd(ctx context.Context) *cobra.Command {
	boundedCmd := &cobra.Command{
		Use:   "bounded",
		Short: "Sender-constrained Proof-of-Possession operations using stored robot key",
	}

	// 1. bounded sign-request
	var (
		method     string
		path       string
		body       string
		datasetStr string
	)

	signCmd := &cobra.Command{
		Use:   "sign-request --method <METHOD> --path <PATH> [--body <str|->]",
		Short: "Generate X-I6-Signature header using the stored robot key and bound token",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			if strings.TrimSpace(method) == "" || strings.TrimSpace(path) == "" {
				return exitcode.New(exitcode.Usage, "--method and --path are required")
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

			rfKey, err := cli.Tokens().GetDatasetRefreshToken(c, targetDs)
			if err != nil || strings.TrimSpace(rfKey) == "" {
				return exitcode.New(exitcode.Auth, fmt.Sprintf("no robot key found for dataset %q in profile %q", targetDs, profile))
			}

			privKey, err := parseEd25519Key(rfKey)
			if err != nil {
				return exitcode.Wrap(exitcode.General, fmt.Errorf("failed to parse stored robot key: %w", err))
			}

			// Obtain access token for the dataset
			tok, err := cli.Tokens().GetValidDatasetAccessToken(c, targetDs)
			if err != nil {
				return exitcode.Wrap(exitcode.Auth, fmt.Errorf("failed to acquire token for signing: %w", err))
			}

			var bodyBytes []byte
			if body == "-" {
				data, err := io.ReadAll(cmd.InOrStdin())
				if err != nil {
					return fmt.Errorf("failed to read body from stdin: %w", err)
				}
				bodyBytes = data
			} else if body != "" {
				bodyBytes = []byte(body)
			}

			parsedPath, err := pathz.Parse(strings.TrimSpace(path))
			if err != nil {
				return exitcode.Wrap(exitcode.Usage, fmt.Errorf("invalid path %q: %w", path, err))
			}

			sigData := reqsign.RequestSignatureData{
				AccessToken: tok,
				Method:      strings.ToUpper(strings.TrimSpace(method)),
				Path:        parsedPath,
			}

			opts := reqsign.RequestSignatureOptions{
				PrivKey: privKey,
			}
			if len(bodyBytes) > 0 {
				hash := sha256.Sum256(bodyBytes)
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

	signCmd.Flags().StringVar(&method, "method", "", "HTTP method (GET, POST, etc.)")
	signCmd.Flags().StringVar(&path, "path", "", "HTTP request path")
	signCmd.Flags().StringVar(&body, "body", "", "Request body string or '-' to read from stdin")
	signCmd.Flags().StringVar(&datasetStr, "dataset", "", "Target dataset identifier")

	boundedCmd.AddCommand(signCmd)
	return boundedCmd
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
