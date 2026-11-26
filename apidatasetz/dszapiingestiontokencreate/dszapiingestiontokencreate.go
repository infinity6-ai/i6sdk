//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package dszapiingestiontokencreate

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

type ReqPayload struct {
	// Optional description note.
	Annotation string `json:"annotation,omitempty"`
}

type Req struct {
	Params  *ReqParams  `json:"params,omitempty"`
	Payload *ReqPayload `json:"payload,omitempty"`
}

// Create Ingestion Token
// Creates a legacy ingestion token for dataset ingestion.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req `json:"req,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/ds/dataset/{dataset}/ingestiontoken/create",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszCreateIngestionTokenId struct{}

func (DszCreateIngestionTokenId) DszCreateIngestionToken(ctx context.Context, api *Api) (int, error) {
	panic("DszCreateIngestionToken not implemented: run 'go generate ./...'")
}
