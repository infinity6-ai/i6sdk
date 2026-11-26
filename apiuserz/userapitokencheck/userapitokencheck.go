//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package userapitokencheck

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type ReqHeaders struct {
	// Standalone 121-byte RefreshAssertion Ed25519 signature.
	SignatureHeader string `json:"x_i6_signature,omitempty"`
	// User Bearer JWT access token or dataset token.
	AuthHeader string `json:"authorization,omitempty"`
}

type ReqPayload struct {
	// Bearer token in payload mode (User JWT or Dataset Access Token).
	Token string `json:"token"`
	// Request signature (envelope + token) in payload mode.
	Signature string `json:"signature"`
	// Target HTTP method.
	Method string `json:"method"`
	// Target URL path.
	Path string `json:"path"`
	// Target URL query string.
	Query string `json:"query"`
	// Target request body.
	Body string `json:"body"`
}

type Req struct {
	Headers *ReqHeaders `json:"headers,omitempty"`
	Payload *ReqPayload `json:"payload,omitempty"`
}

type Breakdown struct {
	// HTTP method.
	Method string `json:"method"`
	// Normalized URL path.
	Path string `json:"path"`
	// Canonical sorted query string.
	CanonicalQuery string `json:"canonical_query"`
	// SHA-256 digest of request body.
	BodySha256 string `json:"body_sha256"`
	// SHA-256 digest of access token.
	TokenSha256 string `json:"token_sha256"`
	// Expiration timestamp in epoch seconds.
	Exp uint64 `json:"exp"`
	// Nonce in lowercase hex.
	Nonce string `json:"nonce"`
	// Reports whether unsigned body flag was enabled.
	UnsignedBody bool `json:"unsigned_body"`
}

type Envelope struct {
	// Wire token type identifier.
	TokenType uint8 `json:"token_type"`
	// Expiration epoch seconds.
	Exp uint64 `json:"exp"`
	// Seconds remaining until expiration.
	TimeRemainingSeconds int64 `json:"time_remaining_seconds"`
	// Reports whether envelope has expired.
	IsExpired bool `json:"is_expired"`
	// Nonce in lowercase hex.
	NonceHex string `json:"nonce_hex"`
	// Reports whether time-windowed replay mode is active.
	IsTimeWindowed bool `json:"is_time_windowed"`
}

type EmbeddedToken struct {
	// Target dataset ID.
	Dataset string `json:"dataset"`
	// Client Ed25519 public key in lowercase hex.
	ClientPubKeyHex string `json:"client_pub_key_hex"`
	// Permission bitmask.
	Grants uint8 `json:"grants"`
	// Solutions bitmask.
	Solutions uint8 `json:"solutions"`
	// Expiration in epoch milliseconds.
	ExpiresAt int64 `json:"expires_at"`
}

type CryptoCheck struct {
	// Client public key used to verify signature.
	VerifiedAgainstPubKey string `json:"verified_against_pubkey"`
	// Reports whether Ed25519 verification succeeded.
	Match bool `json:"match"`
}

type Diagnostics struct {
	// Server-computed canonical string for signatures.
	ServerCanonicalString string `json:"server_canonical_string,omitempty"`
	// Component fields of canonical string.
	Breakdown *Breakdown `json:"breakdown,omitempty"`
	// Envelope metadata.
	Envelope *Envelope `json:"envelope,omitempty"`
	// Embedded dataset access token details.
	EmbeddedToken *EmbeddedToken `json:"embedded_token,omitempty"`
	// Cryptographic signature verification results.
	CryptoCheck *CryptoCheck `json:"crypto_check,omitempty"`

	// User ID extracted from User JWT.
	UserId string `json:"user_id,omitempty"`
	// Audience client app extracted from User JWT.
	Audience string `json:"audience,omitempty"`
	// Expiration timestamp.
	ExpiresAt int64 `json:"expires_at,omitempty"`

	// Dataset ID from dataset access token.
	Dataset string `json:"dataset,omitempty"`
	// Dataset token mode: unbound_bearer or sender_constrained.
	Mode string `json:"mode,omitempty"`
	// Client Ed25519 public key in lowercase hex.
	ClientPubKeyHex string `json:"client_pub_key_hex,omitempty"`
	// Permission bitmask.
	Grants uint8 `json:"grants,omitempty"`
	// Solutions bitmask.
	Solutions uint8 `json:"solutions,omitempty"`

	// Seconds remaining until expiration.
	TimeRemainingSeconds int64 `json:"time_remaining_seconds,omitempty"`
	// Reports whether credential has expired.
	IsExpired bool `json:"is_expired,omitempty"`

	// Actionable troubleshooting guidance.
	Hint string `json:"hint,omitempty"`
}

type RespPayload struct {
	// Indicates whether credential was mathematically verified.
	Valid bool `json:"valid"`
	// Detected credential type: request_signature, user_jwt, or dataset_access_token.
	TokenKind string `json:"token_kind"`
	// Standard error code if invalid.
	ErrorCode string `json:"error_code,omitempty"`
	// High-level summary of verification status.
	Summary string `json:"summary"`
	// Detailed diagnostic breakdown and server-computed values.
	Diagnostics Diagnostics `json:"diagnostics"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Universal Token & Signature Diagnostics
// Auto-detects, validates, and diagnoses any credential supported by authz.Filter without state side-effects.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/userz/token/check",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserTokenCheckId struct{}

func (UserTokenCheckId) UserTokenCheck(ctx context.Context, api *Api) (int, error) {
	panic("UserTokenCheck not implemented: run 'go generate ./...'")
}
