//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package userapidatasettokenrefreshdelete

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
)

type ReqParams struct {
	// Target dataset ID.
	Dataset lakedsz.Dataset `json:"dataset"`
	// Token JTI (public key identifier) to revoke.
	Jti string `json:"jti"`
}

type ReqHeaders struct {
	// User Bearer JWT access token.
	Authorization string `json:"authorization"`
}

type Req struct {
	Params  *ReqParams  `json:"params,omitempty"`
	Headers *ReqHeaders `json:"headers,omitempty"`
}

// Revoke Dataset Refresh Token
// Instantly revokes a machine dataset refresh token by JTI.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req `json:"req,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "DELETE",
		Path:   "/api/userz/dataset-refresh-token/{dataset}/jti/{jti}",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserDatasetTokenRefreshDeleteId struct{}

func (UserDatasetTokenRefreshDeleteId) UserDatasetTokenRefreshDelete(ctx context.Context, api *Api) (int, error) {
	panic("UserDatasetTokenRefreshDelete not implemented: run 'go generate ./...'")
}
