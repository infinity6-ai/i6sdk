package cmdsign

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Prepare mounts the `sign` command and its subcommands to the parent CLI.
func Prepare(parent *cobra.Command) *cobra.Command {
	signCmd := &cobra.Command{
		Use:   "sign",
		Short: "sign HTTP requests and refresh assertions using pure Go standard library",
		Long: `Infinity6 Request Signing (Proof-of-Possession / DPoP) CLI.

This command demonstrates and executes request signing using ONLY the Go standard
library without external cryptographic dependencies.

Subcommands:
  request     Sign an HTTP request for downstream microservices (X-I6-Signature)
  assertion   Sign a standalone refresh assertion for POST /api/userz/refresh-dataset
  keygen      Generate a fresh Ed25519 keypair (Private Key and Public Key/JTI)`,
	}
	parent.AddCommand(signCmd)

	signCmd.AddCommand(newSignRequestCmd())
	signCmd.AddCommand(newSignAssertionCmd())
	signCmd.AddCommand(newKeygenCmd())

	return signCmd
}

func newSignRequestCmd() *cobra.Command {
	var (
		keyStr       string
		method       string
		rawUrl       string
		rawPath      string
		rawQuery     string
		rawBody      string
		accessToken  string
		expEpoch     int64
		ttlStr       string
		nonceStr     string
		unsignedBody bool
		quiet        bool
		showCurl     bool
	)

	cmd := &cobra.Command{
		Use:   "request",
		Short: "sign an HTTP request for downstream microservices",
		Example: `  # Sign a GET request with an access token:
  go run ./main.go sign request --key <base64url_priv_key> --url "https://api.infinity6.ai/api/datasetmanager/data?view=summary" --token <access_token>

  # Sign a POST request with JSON body from file:
  go run ./main.go sign request --key @privkey.txt -m POST -u /api/relevance/infer -b @payload.json -t <access_token>

  # Output only the raw header value for shell scripts:
  SIG=$(go run ./main.go sign request -k @key.txt -u /api/data -t <token> --quiet)
  curl -H "X-I6-Signature: $SIG" https://api.infinity6.ai/api/data`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if keyStr == "" {
				return fmt.Errorf("missing required flag: --key / -k")
			}

			privKey, err := ParsePrivateKey(keyStr)
			if err != nil {
				return fmt.Errorf("invalid private key: %w", err)
			}

			// Parse URL if provided
			path := rawPath
			query := rawQuery
			host := ""
			scheme := "https"

			if rawUrl != "" {
				parsedUrl, err := url.Parse(rawUrl)
				if err != nil {
					return fmt.Errorf("invalid url: %w", err)
				}
				if parsedUrl.Host != "" {
					host = parsedUrl.Host
				}
				if parsedUrl.Scheme != "" {
					scheme = parsedUrl.Scheme
				}
				if parsedUrl.Path != "" {
					path = parsedUrl.Path
				}
				if parsedUrl.RawQuery != "" {
					query = parsedUrl.RawQuery
				}
			}

			if path == "" {
				path = "/"
			}

			// Read body if from file
			bodyBytes := []byte(rawBody)
			if strings.HasPrefix(rawBody, "@") {
				fBytes, err := os.ReadFile(rawBody[1:])
				if err != nil {
					return fmt.Errorf("failed to read body file: %w", err)
				}
				bodyBytes = fBytes
			}

			// Calculate expiration
			var exp uint64
			if expEpoch > 0 {
				exp = uint64(expEpoch)
			} else {
				ttl := 5 * time.Minute
				if ttlStr != "" {
					parsedTtl, err := time.ParseDuration(ttlStr)
					if err != nil {
						return fmt.Errorf("invalid ttl duration: %w", err)
					}
					ttl = parsedTtl
				}
				exp = uint64(time.Now().Add(ttl).Unix())
			}

			// Nonce
			var nonce [16]byte
			switch strings.ToLower(strings.TrimSpace(nonceStr)) {
			case "", "zero", "zeros":
				// 16 zero bytes (time-windowed replay mode)
				nonce = [16]byte{}
			case "random":
				rnd, err := GenerateRandomNonce()
				if err != nil {
					return err
				}
				nonce = rnd
			default:
				// Hex encoded 16-byte nonce
				b, err := hex.DecodeString(nonceStr)
				if err != nil || len(b) != 16 {
					return fmt.Errorf("invalid nonce: must be 'random', 'zero', or a 32-character hex string")
				}
				copy(nonce[:], b)
			}

			sig, err := SignRequest(privKey, method, path, query, bodyBytes, accessToken, exp, nonce, unsignedBody)
			if err != nil {
				return err
			}

			if quiet {
				fmt.Println(sig.HeaderValue)
				return nil
			}

			fmt.Println("================================================================================")
			fmt.Println("                       INFINITY6 REQUEST SIGNING RESULT                        ")
			fmt.Println("================================================================================")
			fmt.Printf("HTTP Method:       %s\n", strings.ToUpper(method))
			fmt.Printf("Normalized Path:   %s\n", path)
			fmt.Printf("Query String:      %s\n", query)
			fmt.Printf("Expires At (Unix): %d (%s)\n", sig.Exp, time.Unix(int64(sig.Exp), 0).UTC().Format(time.RFC3339))
			fmt.Printf("Nonce (Hex):       %s\n", hex.EncodeToString(sig.Nonce[:]))
			fmt.Printf("Unsigned Body:     %t\n", unsignedBody)
			fmt.Println("--------------------------------------------------------------------------------")
			fmt.Println("CANONICAL STRING (7 lines signed by Ed25519):")
			for i, line := range strings.Split(sig.CanonicalString, "\n") {
				lineName := canonicalLineDesc(i)
				fmt.Printf("  [%d] %-18s: %q\n", i+1, lineName, line)
			}
			fmt.Println("--------------------------------------------------------------------------------")
			fmt.Println("X-I6-Signature Header:")
			fmt.Printf("  %s\n", sig.HeaderValue)
			fmt.Println("================================================================================")

			if showCurl {
				targetHost := host
				if targetHost == "" {
					targetHost = "api.infinity6.ai"
				}
				fullUrl := fmt.Sprintf("%s://%s%s", scheme, targetHost, path)
				if query != "" {
					fullUrl += "?" + query
				}

				fmt.Println("\nExample curl command:")
				fmt.Printf("curl -X %s %q \\\n", strings.ToUpper(method), fullUrl)
				fmt.Printf("  -H %q", fmt.Sprintf("X-I6-Signature: %s", sig.HeaderValue))
				if len(bodyBytes) > 0 {
					fmt.Printf(" \\\n  -H %q \\\n  -d %q", "Content-Type: application/json", string(bodyBytes))
				}
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&keyStr, "key", "k", "", "Ed25519 private key (Base64URL, Hex, or @filepath)")
	cmd.Flags().StringVarP(&method, "method", "m", "GET", "HTTP verb (GET, POST, PUT, DELETE)")
	cmd.Flags().StringVarP(&rawUrl, "url", "u", "", "Target URL (e.g. https://api.infinity6.ai/api/data?q=1)")
	cmd.Flags().StringVarP(&rawPath, "path", "p", "/", "Request path (if --url is not used)")
	cmd.Flags().StringVarP(&rawQuery, "query", "q", "", "Request query string (if --url is not used)")
	cmd.Flags().StringVarP(&rawBody, "body", "b", "", "Request body string or @filepath")
	cmd.Flags().StringVarP(&accessToken, "token", "t", "", "Dataset access token string for composite signature")
	cmd.Flags().Int64Var(&expEpoch, "exp", 0, "Expiration Unix epoch timestamp in seconds (overrides --ttl)")
	cmd.Flags().StringVar(&ttlStr, "ttl", "5m", "Expiration TTL duration from now (e.g. 5m, 1h)")
	cmd.Flags().StringVar(&nonceStr, "nonce", "random", "Nonce: 'random', 'zero' (for time-windowed), or 32-char hex")
	cmd.Flags().BoolVar(&unsignedBody, "unsigned-body", false, "Permit any request body (sets Bit 7 of envelope)")
	cmd.Flags().BoolVarP(&quiet, "quiet", "s", false, "Print ONLY the raw X-I6-Signature value")
	cmd.Flags().BoolVar(&showCurl, "curl", true, "Print ready-to-run curl example")

	return cmd
}

func newSignAssertionCmd() *cobra.Command {
	var (
		keyStr       string
		method       string
		path         string
		query        string
		rawBody      string
		expEpoch     int64
		ttlStr       string
		nonceStr     string
		unsignedBody bool
		quiet        bool
		showCurl     bool
	)

	cmd := &cobra.Command{
		Use:   "assertion",
		Short: "sign a standalone refresh assertion for POST /api/userz/refresh-dataset",
		Example: `  # Sign a refresh assertion using a refresh token private key:
  go run ./main.go sign assertion --key @refresh_priv.key

  # Output raw 162-char Base64URL string for scripts:
  ASSERTION=$(go run ./main.go sign assertion -k @refresh_priv.key --quiet)
  curl -X POST https://api.infinity6.ai/api/userz/refresh-dataset \
    -H "X-I6-Signature: $ASSERTION" \
    -H "Content-Type: application/json" \
    -d '{"dataset":"ds1"}'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if keyStr == "" {
				return fmt.Errorf("missing required flag: --key / -k")
			}

			privKey, err := ParsePrivateKey(keyStr)
			if err != nil {
				return fmt.Errorf("invalid private key: %w", err)
			}

			bodyBytes := []byte(rawBody)
			if strings.HasPrefix(rawBody, "@") {
				fBytes, err := os.ReadFile(rawBody[1:])
				if err != nil {
					return fmt.Errorf("failed to read body file: %w", err)
				}
				bodyBytes = fBytes
			}

			var exp uint64
			if expEpoch > 0 {
				exp = uint64(expEpoch)
			} else {
				ttl := 2 * time.Minute
				if ttlStr != "" {
					parsedTtl, err := time.ParseDuration(ttlStr)
					if err != nil {
						return fmt.Errorf("invalid ttl duration: %w", err)
					}
					ttl = parsedTtl
				}
				exp = uint64(time.Now().Add(ttl).Unix())
			}

			var nonce [16]byte
			switch strings.ToLower(strings.TrimSpace(nonceStr)) {
			case "", "zero", "zeros":
				nonce = [16]byte{}
			case "random":
				rnd, err := GenerateRandomNonce()
				if err != nil {
					return err
				}
				nonce = rnd
			default:
				b, err := hex.DecodeString(nonceStr)
				if err != nil || len(b) != 16 {
					return fmt.Errorf("invalid nonce: must be 'random', 'zero', or a 32-character hex string")
				}
				copy(nonce[:], b)
			}

			assertion, err := SignRefreshAssertion(privKey, method, path, query, bodyBytes, exp, nonce, unsignedBody)
			if err != nil {
				return err
			}

			if quiet {
				fmt.Println(assertion.AssertionBase64)
				return nil
			}

			fmt.Println("================================================================================")
			fmt.Println("                  INFINITY6 REFRESH ASSERTION SIGNING RESULT                    ")
			fmt.Println("================================================================================")
			fmt.Printf("HTTP Method:       %s\n", strings.ToUpper(method))
			fmt.Printf("Path:              %s\n", path)
			fmt.Printf("Public Key (JTI):  %s\n", assertion.PublicKeyBase64)
			fmt.Printf("Expires At (Unix): %d (%s)\n", assertion.Exp, time.Unix(int64(assertion.Exp), 0).UTC().Format(time.RFC3339))
			fmt.Printf("Nonce (Hex):       %s\n", hex.EncodeToString(assertion.Nonce[:]))
			fmt.Printf("Unsigned Body:     %t\n", unsignedBody)
			fmt.Println("--------------------------------------------------------------------------------")
			fmt.Println("CANONICAL STRING (7 lines signed by Ed25519):")
			for i, line := range strings.Split(assertion.CanonicalString, "\n") {
				lineName := canonicalLineDesc(i)
				fmt.Printf("  [%d] %-18s: %q\n", i+1, lineName, line)
			}
			fmt.Println("--------------------------------------------------------------------------------")
			fmt.Println("Refresh Assertion (121 bytes binary -> 162 Base64URL chars):")
			fmt.Printf("  %s\n", assertion.AssertionBase64)
			fmt.Println("================================================================================")

			if showCurl {
				fmt.Println("\nExample curl command to refresh dataset access token:")
				fmt.Printf("curl -X POST https://api.infinity6.ai/api/userz/refresh-dataset \\\n")
				fmt.Printf("  -H %q \\\n", fmt.Sprintf("X-I6-Signature: %s", assertion.AssertionBase64))
				fmt.Printf("  -H \"Content-Type: application/json\" \\\n")
				fmt.Printf("  -d '{\"dataset\":\"<your-dataset>\"}'\n")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&keyStr, "key", "k", "", "Ed25519 private key (Base64URL, Hex, or @filepath)")
	cmd.Flags().StringVarP(&method, "method", "m", "POST", "HTTP method (typically POST)")
	cmd.Flags().StringVarP(&path, "path", "p", "/api/userz/refresh-dataset", "Request path")
	cmd.Flags().StringVarP(&query, "query", "q", "", "Query string")
	cmd.Flags().StringVarP(&rawBody, "body", "b", "", "Request body string or @filepath")
	cmd.Flags().Int64Var(&expEpoch, "exp", 0, "Expiration Unix epoch timestamp in seconds")
	cmd.Flags().StringVar(&ttlStr, "ttl", "2m", "Expiration TTL duration from now")
	cmd.Flags().StringVar(&nonceStr, "nonce", "random", "Nonce: 'random', 'zero', or 32-char hex")
	cmd.Flags().BoolVar(&unsignedBody, "unsigned-body", true, "Permit any body (default true for refresh assertion)")
	cmd.Flags().BoolVarP(&quiet, "quiet", "s", false, "Print ONLY the raw 162-char assertion value")
	cmd.Flags().BoolVar(&showCurl, "curl", true, "Print ready-to-run curl example")

	return cmd
}

func newKeygenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "keygen",
		Short: "generate an Ed25519 keypair for test and development",
		RunE: func(cmd *cobra.Command, args []string) error {
			pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
			if err != nil {
				return fmt.Errorf("failed to generate key: %w", err)
			}

			pubB64 := base64.RawURLEncoding.EncodeToString(pubKey)
			privB64 := base64.RawURLEncoding.EncodeToString(privKey)
			seedB64 := base64.RawURLEncoding.EncodeToString(privKey.Seed())

			fmt.Println("================================================================================")
			fmt.Println("                        NEW ED25519 KEYPAIR GENERATED                           ")
			fmt.Println("================================================================================")
			fmt.Printf("Public Key (JTI / Base64URL - 32B):   %s\n", pubB64)
			fmt.Printf("Private Key (Base64URL - 64B):        %s\n", privB64)
			fmt.Printf("Seed Only (Base64URL - 32B):          %s\n", seedB64)
			fmt.Println("--------------------------------------------------------------------------------")
			fmt.Printf("Public Key (Hex):  %s\n", hex.EncodeToString(pubKey))
			fmt.Printf("Private Key (Hex): %s\n", hex.EncodeToString(privKey))
			fmt.Println("================================================================================")
			return nil
		},
	}
}

func canonicalLineDesc(idx int) string {
	names := []string{
		"METHOD",
		"PATH",
		"CANONICAL QUERY",
		"BODY SHA-256",
		"TOKEN SHA-256",
		"EXPIRATION SECONDS",
		"NONCE HEX",
	}
	if idx >= 0 && idx < len(names) {
		return names[idx]
	}
	return "UNKNOWN"
}
