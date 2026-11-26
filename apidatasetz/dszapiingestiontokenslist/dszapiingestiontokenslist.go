//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package dszapiingestiontokenslist

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
	// List of ingestion tokens.
	Tokens []dsztoken.DatasetToken `json:"tokens"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// List Ingestion Tokens
// Lists all ingestion tokens provisioned for a dataset.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/ingestiontoken/list",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszListIngestionTokensId struct{}

func (DszListIngestionTokensId) DszListIngestionTokens(ctx context.Context, api *Api) (int, error) {
	panic("DszListIngestionTokens not implemented: run 'go generate ./...'")
}
