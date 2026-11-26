//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package pipezget

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/pipez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type ReqParams struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
	// Pipeline id.
	Id string `json:"id"`
}

type Req struct {
	Params *ReqParams `json:"params,omitempty"`
}

// Response containing the pipeline.
type RespPayload struct {
	// Pipeline and its states.
	Pipeline *pipez.Pipeline `json:"pipeline"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Get Pipeline
// Gets a pipeline of a dataset by its id.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/pipe/dataset/{dataset}/id/{id}",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type PipezGetId struct{}

func (PipezGetId) PipezGet(ctx context.Context, api *Api) (int, error) {
	panic("PipezGet not implemented: run 'go generate ./...'")
}
