//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package userapidatasettokenrefreshcreate

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
)

type ReqParams struct {
	// Target dataset ID.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type ReqHeaders struct {
	// User Bearer JWT access token.
	Authorization string `json:"authorization"`
}

type ReqPayload struct {
	// Solution domain restriction.
	Solution solution.Solution `json:"solution"`
	// Grants bitmask.
	Grants grants.Grant `json:"grants,omitempty"`
	// Metadata labels.
	Labels map[string]string `json:"labels,omitempty"`
	// Requested expiration duration in seconds.
	ExpSeconds int64 `json:"exp_seconds,omitempty"`
	// Optional description for token management.
	Annotation string `json:"annotation,omitempty"`
}

type Req struct {
	Params  *ReqParams  `json:"params,omitempty"`
	Headers *ReqHeaders `json:"headers,omitempty"`
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// Token JTI public identifier.
	Jti string `json:"jti"`
	// Target dataset ID.
	Dataset lakedsz.Dataset `json:"dataset"`
	// Base64URL-encoded Ed25519 private key (returned once only).
	PrivKey string `json:"priv_key"`
	// Solution domain restriction.
	Solution solution.Solution `json:"solution"`
	// Granted permissions.
	Grants grants.Grant `json:"grants"`
	// Metadata labels.
	Labels map[string]string `json:"labels,omitempty"`
	// Expiration epoch milliseconds.
	ExpiresAt int64 `json:"expires_at"`
	// Description note.
	Annotation string `json:"annotation,omitempty"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Provision Dataset Refresh Token
// Provisions a long-lived machine dataset refresh token.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/userz/dataset-refresh-token/{dataset}",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserDatasetTokenRefreshCreateId struct{}

func (UserDatasetTokenRefreshCreateId) UserDatasetTokenRefreshCreate(ctx context.Context, api *Api) (int, error) {
	panic("UserDatasetTokenRefreshCreate not implemented: run 'go generate ./...'")
}
