//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package userapidatasettokenrefresh

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
)

type ReqHeaders struct {
	// Standalone 121-byte RefreshAssertion Ed25519 signature.
	SignatureHeader string `json:"x_i6_signature,omitempty"`
}

type ReqPayload struct {
	// Target dataset ID.
	Dataset lakedsz.Dataset `json:"dataset"`
	// Optional sub-solution restriction.
	Subsolution solution.Solution `json:"subsolution,omitempty"`
	// Optional reduced grant bitmask.
	Subgrants grants.Grant `json:"subgrants,omitempty"`
	// Additional restrictive label constraints.
	Newlabels map[string]string `json:"newlabels,omitempty"`
	// Requested lifetime in seconds.
	ExpSeconds int64 `json:"exp_seconds,omitempty"`
	// If true, mints an unbound bearer token usable with Authorization: Bearer.
	Unbound bool `json:"unbound,omitempty"`
}

type Req struct {
	Headers *ReqHeaders `json:"headers,omitempty"`
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// Binary-packed Base64URL dataset access token.
	AccessToken string `json:"access_token"`
	// Target dataset ID.
	Dataset lakedsz.Dataset `json:"dataset"`
	// Effective solution domain.
	Solution solution.Solution `json:"solution"`
	// Effective uint8 grant bitmask.
	Grants grants.Grant `json:"grants"`
	// Effective metadata labels.
	Labels map[string]string `json:"labels,omitempty"`
	// Expiration epoch milliseconds.
	ExpiresAt int64 `json:"expires_at"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Refresh Dataset Access Token
// Refreshes and attenuates a dataset access token using an Ed25519 signature assertion.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/userz/refresh-dataset",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserDatasetTokenRefreshId struct{}

func (UserDatasetTokenRefreshId) UserDatasetTokenRefresh(ctx context.Context, api *Api) (int, error) {
	panic("UserDatasetTokenRefresh not implemented: run 'go generate ./...'")
}
