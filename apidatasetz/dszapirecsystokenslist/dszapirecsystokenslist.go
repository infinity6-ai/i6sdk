//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package dszapirecsystokenslist

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/dsztoken"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type ReqParams struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type Req struct {
	Params *ReqParams `json:"params,omitempty"`
}

type RespPayload struct {
	// List of recsys tokens.
	Tokens []dsztoken.DatasetToken `json:"tokens"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// List Recsys Tokens
// Lists all recsys tokens provisioned for a dataset.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/recsystoken/list",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszListRecsysTokensId struct{}

func (DszListRecsysTokensId) DszListRecsysTokens(ctx context.Context, api *Api) (int, error) {
	panic("DszListRecsysTokens not implemented: run 'go generate ./...'")
}
