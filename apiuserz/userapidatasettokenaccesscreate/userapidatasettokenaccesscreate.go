//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package userapidatasettokenaccesscreate

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
	// Target solution domain restriction (empty for all).
	Solution solution.Solution `json:"solution"`
	// Compact uint8 grant bitmask.
	Grants grants.Grant `json:"grants,omitempty"`
	// Contextual metadata labels.
	Labels map[string]string `json:"labels,omitempty"`
	// Base64URL client public key for PoP binding (empty for unbound bearer).
	ClientPubKey string `json:"client_pub_key,omitempty"`
	// Requested lifetime in seconds.
	ExpSeconds int64 `json:"exp_seconds,omitempty"`
}

type Req struct {
	Params  *ReqParams  `json:"params,omitempty"`
	Headers *ReqHeaders `json:"headers,omitempty"`
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// Binary-packed Base64URL dataset access token.
	AccessToken string `json:"access_token"`
	// Target dataset ID.
	Dataset lakedsz.Dataset `json:"dataset"`
	// Solution domain restriction.
	Solution solution.Solution `json:"solution"`
	// Effective uint8 grant bitmask.
	Grants grants.Grant `json:"grants"`
	// Contextual metadata labels.
	Labels map[string]string `json:"labels,omitempty"`
	// Expiration epoch milliseconds.
	ExpiresAt int64 `json:"expires_at"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Create Dataset Access Token
// Mints a short-lived dataset access token for an authenticated user.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/userz/dataset-access-token/{dataset}",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserDatasetTokenAccessCreateId struct{}

func (UserDatasetTokenAccessCreateId) UserDatasetTokenAccessCreate(ctx context.Context, api *Api) (int, error) {
	panic("UserDatasetTokenAccessCreate not implemented: run 'go generate ./...'")
}
