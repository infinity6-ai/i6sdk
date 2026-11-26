//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package userapidatasettokenrefreshlist

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
)

type RefreshTokenSummary struct {
	// Token JTI.
	Jti string `json:"jti"`
	// Target dataset ID.
	Dataset lakedsz.Dataset `json:"dataset"`
	// Solution domain restriction.
	Solution solution.Solution `json:"solution"`
	// Grants bitmask.
	Grants grants.Grant `json:"grants"`
	// Metadata labels.
	Labels map[string]string `json:"labels,omitempty"`
	// User ID who created the token.
	CreatedBy string `json:"created_by"`
	// Creation timestamp in millis.
	CreatedAt int64 `json:"created_at"`
	// Epoch milliseconds when token was last refreshed (0 = never used).
	LastUsedAt int64 `json:"last_used_at"`
	// Expiration timestamp in millis.
	ExpiresAt int64 `json:"expires_at"`
	// Revocation status flag.
	Revoked bool `json:"revoked"`
	// Description note.
	Annotation string `json:"annotation,omitempty"`
}

type ReqParams struct {
	// Target dataset ID.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type ReqHeaders struct {
	// User Bearer JWT access token.
	Authorization string `json:"authorization"`
}

type Req struct {
	Params  *ReqParams  `json:"params,omitempty"`
	Headers *ReqHeaders `json:"headers,omitempty"`
}

type RespPayload struct {
	// List of refresh token records.
	Tokens []RefreshTokenSummary `json:"tokens"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// List Dataset Refresh Tokens
// Lists all active machine refresh tokens provisioned for a dataset.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/userz/dataset-refresh-token/{dataset}",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserDatasetTokenRefreshListId struct{}

func (UserDatasetTokenRefreshListId) UserDatasetTokenRefreshList(ctx context.Context, api *Api) (int, error) {
	panic("UserDatasetTokenRefreshList not implemented: run 'go generate ./...'")
}
