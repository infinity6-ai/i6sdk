package userapitokencheck

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type Req struct {
	SignatureHeader string `json:"x_i6_signature"`
	AuthHeader      string `json:"authorization"`
	Token           string `json:"token"`
	Signature       string `json:"signature"`
	Method          string `json:"method"`
	Path            string `json:"path"`
	Query           string `json:"query"`
	Body            string `json:"body"`
}

type Breakdown struct {
	Method         string `json:"method"`
	Path           string `json:"path"`
	CanonicalQuery string `json:"canonical_query"`
	BodySha256     string `json:"body_sha256"`
	TokenSha256    string `json:"token_sha256"`
	Exp            uint64 `json:"exp"`
	Nonce          string `json:"nonce"`
	UnsignedBody   bool   `json:"unsigned_body"`
}

type Envelope struct {
	TokenType            uint8  `json:"token_type"`
	Exp                  uint64 `json:"exp"`
	TimeRemainingSeconds int64  `json:"time_remaining_seconds"`
	IsExpired            bool   `json:"is_expired"`
	NonceHex             string `json:"nonce_hex"`
	IsTimeWindowed       bool   `json:"is_time_windowed"`
}

type EmbeddedToken struct {
	Dataset         string `json:"dataset"`
	ClientPubKeyHex string `json:"client_pub_key_hex"`
	Grants          uint8  `json:"grants"`
	Solutions       uint8  `json:"solutions"`
	ExpiresAt       int64  `json:"expires_at"`
}

type CryptoCheck struct {
	VerifiedAgainstPubKey string `json:"verified_against_pubkey"`
	Match                 bool   `json:"match"`
}

type Diagnostics struct {
	// Request signature fields
	ServerCanonicalString string         `json:"server_canonical_string,omitempty"`
	Breakdown             *Breakdown     `json:"breakdown,omitempty"`
	Envelope              *Envelope      `json:"envelope,omitempty"`
	EmbeddedToken         *EmbeddedToken `json:"embedded_token,omitempty"`
	CryptoCheck           *CryptoCheck   `json:"crypto_check,omitempty"`

	// User JWT fields
	UserId    string `json:"user_id,omitempty"`
	Audience  string `json:"audience,omitempty"`
	ExpiresAt int64  `json:"expires_at,omitempty"`

	// Dataset Access Token fields
	Dataset         string `json:"dataset,omitempty"`
	Mode            string `json:"mode,omitempty"`
	ClientPubKeyHex string `json:"client_pub_key_hex,omitempty"`
	Grants          uint8  `json:"grants,omitempty"`
	Solutions       uint8  `json:"solutions,omitempty"`

	// Common lifetime fields
	TimeRemainingSeconds int64 `json:"time_remaining_seconds,omitempty"`
	IsExpired            bool  `json:"is_expired,omitempty"`

	// Troubleshooting hint
	Hint string `json:"hint,omitempty"`
}

type Resp struct {
	Valid       bool        `json:"valid"`
	TokenKind   string      `json:"token_kind"` // "request_signature", "user_jwt", "dataset_access_token"
	ErrorCode   string      `json:"error_code,omitempty"`
	Summary     string      `json:"summary"`
	Diagnostics Diagnostics `json:"diagnostics"`
}

type Api struct {
	Req  Req  `json:"req"`
	Resp Resp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id: "user-token-check",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "Universal Token & Signature Diagnostics",
				Summary: "Auto-detects, validates, and diagnoses any credential supported by authz.Filter without state side-effects.",
				Labels:  map[string]string{"servname": "userz"},
			}
		},
		Method: "POST",
		Path:   "/api/userz/token/check",
	}
}

func (d *Api) ApiDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams:  nil,
		QueryParams: nil,
		ReqHeaders: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"x_i6_signature": {
						Desc:     &schemaz.Desc{Summary: "Optional request signature header (X-I6-Signature)."},
						UrlValue: schemaz.ParseStr(&d.Req.SignatureHeader),
					},
					"authorization": {
						Desc:     &schemaz.Desc{Summary: "Optional Bearer authorization header."},
						UrlValue: schemaz.ParseStr(&d.Req.AuthHeader),
					},
				}
			},
		},
		ReqBody: &schemaz.Schema{
			Desc: &schemaz.Desc{Summary: "Diagnostic request payload."},
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"token": {
						Desc:  &schemaz.Desc{Summary: "Bearer token in payload mode (User JWT or Dataset Access Token)."},
						Field: func() any { return &d.Req.Token },
					},
					"signature": {
						Desc:  &schemaz.Desc{Summary: "Request signature (envelope + token) in payload mode."},
						Field: func() any { return &d.Req.Signature },
					},
					"method": {
						Desc:  &schemaz.Desc{Summary: "Target HTTP method."},
						Field: func() any { return &d.Req.Method },
					},
					"path": {
						Desc:  &schemaz.Desc{Summary: "Target URL path."},
						Field: func() any { return &d.Req.Path },
					},
					"query": {
						Desc:  &schemaz.Desc{Summary: "Target URL query string."},
						Field: func() any { return &d.Req.Query },
					},
					"body": {
						Desc:  &schemaz.Desc{Summary: "Target request body."},
						Field: func() any { return &d.Req.Body },
					},
				}
			},
		},
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Desc: &schemaz.Desc{Summary: "Diagnostic verification results."},
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"valid": {
						Desc:  &schemaz.Desc{Summary: "Indicates whether credential was mathematically verified."},
						Field: func() any { return &d.Resp.Valid },
					},
					"token_kind": {
						Desc:  &schemaz.Desc{Summary: "Detected credential type: request_signature, user_jwt, or dataset_access_token."},
						Field: func() any { return &d.Resp.TokenKind },
					},
					"error_code": {
						Desc:  &schemaz.Desc{Summary: "Standard error code if invalid."},
						Field: func() any { return &d.Resp.ErrorCode },
					},
					"summary": {
						Desc:  &schemaz.Desc{Summary: "High-level summary of verification status."},
						Field: func() any { return &d.Resp.Summary },
					},
					"diagnostics": {
						Desc: &schemaz.Desc{Summary: "Detailed diagnostic breakdown and server-computed values."},
						Object: func(read bool) map[string]*schemaz.Schema {
							return map[string]*schemaz.Schema{
								"server_canonical_string": {
									Desc:  &schemaz.Desc{Summary: "Server-computed canonical string for signatures."},
									Field: func() any { return &d.Resp.Diagnostics.ServerCanonicalString },
								},
								"user_id": {
									Desc:  &schemaz.Desc{Summary: "User ID extracted from User JWT."},
									Field: func() any { return &d.Resp.Diagnostics.UserId },
								},
								"audience": {
									Desc:  &schemaz.Desc{Summary: "Audience client app extracted from User JWT."},
									Field: func() any { return &d.Resp.Diagnostics.Audience },
								},
								"dataset": {
									Desc:  &schemaz.Desc{Summary: "Dataset ID from dataset access token."},
									Field: func() any { return &d.Resp.Diagnostics.Dataset },
								},
								"mode": {
									Desc:  &schemaz.Desc{Summary: "Dataset token mode: unbound_bearer or sender_constrained."},
									Field: func() any { return &d.Resp.Diagnostics.Mode },
								},
								"client_pub_key_hex": {
									Desc:  &schemaz.Desc{Summary: "Client Ed25519 public key in lowercase hex."},
									Field: func() any { return &d.Resp.Diagnostics.ClientPubKeyHex },
								},
								"grants": {
									Desc:  &schemaz.Desc{Summary: "Permission bitmask."},
									Field: func() any { return &d.Resp.Diagnostics.Grants },
								},
								"solutions": {
									Desc:  &schemaz.Desc{Summary: "Solutions bitmask."},
									Field: func() any { return &d.Resp.Diagnostics.Solutions },
								},
								"expires_at": {
									Desc:  &schemaz.Desc{Summary: "Expiration timestamp."},
									Field: func() any { return &d.Resp.Diagnostics.ExpiresAt },
								},
								"time_remaining_seconds": {
									Desc:  &schemaz.Desc{Summary: "Seconds remaining until expiration."},
									Field: func() any { return &d.Resp.Diagnostics.TimeRemainingSeconds },
								},
								"is_expired": {
									Desc:  &schemaz.Desc{Summary: "Reports whether credential has expired."},
									Field: func() any { return &d.Resp.Diagnostics.IsExpired },
								},
								"hint": {
									Desc:  &schemaz.Desc{Summary: "Actionable troubleshooting guidance."},
									Field: func() any { return &d.Resp.Diagnostics.Hint },
								},
								"breakdown": {
									Desc: &schemaz.Desc{Summary: "Component fields of canonical string."},
									Object: func(read bool) map[string]*schemaz.Schema {
										if read && d.Resp.Diagnostics.Breakdown == nil {
											return nil
										}
										if d.Resp.Diagnostics.Breakdown == nil {
											d.Resp.Diagnostics.Breakdown = &Breakdown{}
										}
										return map[string]*schemaz.Schema{
											"method": {
												Desc:  &schemaz.Desc{Summary: "HTTP method."},
												Field: func() any { return &d.Resp.Diagnostics.Breakdown.Method },
											},
											"path": {
												Desc:  &schemaz.Desc{Summary: "Normalized URL path."},
												Field: func() any { return &d.Resp.Diagnostics.Breakdown.Path },
											},
											"canonical_query": {
												Desc:  &schemaz.Desc{Summary: "Canonical sorted query string."},
												Field: func() any { return &d.Resp.Diagnostics.Breakdown.CanonicalQuery },
											},
											"body_sha256": {
												Desc:  &schemaz.Desc{Summary: "SHA-256 digest of request body."},
												Field: func() any { return &d.Resp.Diagnostics.Breakdown.BodySha256 },
											},
											"token_sha256": {
												Desc:  &schemaz.Desc{Summary: "SHA-256 digest of access token."},
												Field: func() any { return &d.Resp.Diagnostics.Breakdown.TokenSha256 },
											},
											"exp": {
												Desc:  &schemaz.Desc{Summary: "Expiration timestamp in epoch seconds."},
												Field: func() any { return &d.Resp.Diagnostics.Breakdown.Exp },
											},
											"nonce": {
												Desc:  &schemaz.Desc{Summary: "Nonce in lowercase hex."},
												Field: func() any { return &d.Resp.Diagnostics.Breakdown.Nonce },
											},
											"unsigned_body": {
												Desc:  &schemaz.Desc{Summary: "Reports whether unsigned body flag was enabled."},
												Field: func() any { return &d.Resp.Diagnostics.Breakdown.UnsignedBody },
											},
										}
									},
								},
								"envelope": {
									Desc: &schemaz.Desc{Summary: "Envelope metadata."},
									Object: func(read bool) map[string]*schemaz.Schema {
										if read && d.Resp.Diagnostics.Envelope == nil {
											return nil
										}
										if d.Resp.Diagnostics.Envelope == nil {
											d.Resp.Diagnostics.Envelope = &Envelope{}
										}
										return map[string]*schemaz.Schema{
											"token_type": {
												Desc:  &schemaz.Desc{Summary: "Wire token type identifier."},
												Field: func() any { return &d.Resp.Diagnostics.Envelope.TokenType },
											},
											"exp": {
												Desc:  &schemaz.Desc{Summary: "Expiration epoch seconds."},
												Field: func() any { return &d.Resp.Diagnostics.Envelope.Exp },
											},
											"time_remaining_seconds": {
												Desc:  &schemaz.Desc{Summary: "Seconds remaining until expiration."},
												Field: func() any { return &d.Resp.Diagnostics.Envelope.TimeRemainingSeconds },
											},
											"is_expired": {
												Desc:  &schemaz.Desc{Summary: "Reports whether envelope has expired."},
												Field: func() any { return &d.Resp.Diagnostics.Envelope.IsExpired },
											},
											"nonce_hex": {
												Desc:  &schemaz.Desc{Summary: "Nonce in lowercase hex."},
												Field: func() any { return &d.Resp.Diagnostics.Envelope.NonceHex },
											},
											"is_time_windowed": {
												Desc:  &schemaz.Desc{Summary: "Reports whether time-windowed replay mode is active."},
												Field: func() any { return &d.Resp.Diagnostics.Envelope.IsTimeWindowed },
											},
										}
									},
								},
								"embedded_token": {
									Desc: &schemaz.Desc{Summary: "Embedded dataset access token details."},
									Object: func(read bool) map[string]*schemaz.Schema {
										if read && d.Resp.Diagnostics.EmbeddedToken == nil {
											return nil
										}
										if d.Resp.Diagnostics.EmbeddedToken == nil {
											d.Resp.Diagnostics.EmbeddedToken = &EmbeddedToken{}
										}
										return map[string]*schemaz.Schema{
											"dataset": {
												Desc:  &schemaz.Desc{Summary: "Target dataset ID."},
												Field: func() any { return &d.Resp.Diagnostics.EmbeddedToken.Dataset },
											},
											"client_pub_key_hex": {
												Desc:  &schemaz.Desc{Summary: "Client Ed25519 public key in lowercase hex."},
												Field: func() any { return &d.Resp.Diagnostics.EmbeddedToken.ClientPubKeyHex },
											},
											"grants": {
												Desc:  &schemaz.Desc{Summary: "Permission bitmask."},
												Field: func() any { return &d.Resp.Diagnostics.EmbeddedToken.Grants },
											},
											"solutions": {
												Desc:  &schemaz.Desc{Summary: "Solutions bitmask."},
												Field: func() any { return &d.Resp.Diagnostics.EmbeddedToken.Solutions },
											},
											"expires_at": {
												Desc:  &schemaz.Desc{Summary: "Expiration in epoch milliseconds."},
												Field: func() any { return &d.Resp.Diagnostics.EmbeddedToken.ExpiresAt },
											},
										}
									},
								},
								"crypto_check": {
									Desc: &schemaz.Desc{Summary: "Cryptographic signature verification results."},
									Object: func(read bool) map[string]*schemaz.Schema {
										if read && d.Resp.Diagnostics.CryptoCheck == nil {
											return nil
										}
										if d.Resp.Diagnostics.CryptoCheck == nil {
											d.Resp.Diagnostics.CryptoCheck = &CryptoCheck{}
										}
										return map[string]*schemaz.Schema{
											"verified_against_pubkey": {
												Desc:  &schemaz.Desc{Summary: "Client public key used to verify signature."},
												Field: func() any { return &d.Resp.Diagnostics.CryptoCheck.VerifiedAgainstPubKey },
											},
											"match": {
												Desc:  &schemaz.Desc{Summary: "Reports whether Ed25519 verification succeeded."},
												Field: func() any { return &d.Resp.Diagnostics.CryptoCheck.Match },
											},
										}
									},
								},
							}
						},
					},
				}
			},
		},
	}
}
