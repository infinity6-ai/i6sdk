//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package dszapitokendelete

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type ReqParams struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type ReqQuery struct {
	// Token identifier to revoke.
	Token string `json:"token"`
}

type Req struct {
	Params *ReqParams `json:"params,omitempty"`
	Query  *ReqQuery  `json:"query,omitempty"`
}

// Delete Dataset Token
// Revokes a legacy dataset token.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req `json:"req,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "DELETE",
		Path:   "/api/ds/dataset/{dataset}/token/delete",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszDeleteTokenId struct{}

func (DszDeleteTokenId) DszDeleteToken(ctx context.Context, api *Api) (int, error) {
	panic("DszDeleteToken not implemented: run 'go generate ./...'")
}
