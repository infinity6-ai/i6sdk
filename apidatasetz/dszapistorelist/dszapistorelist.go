//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package dszapistorelist

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/dszstore"
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
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
	// Stores configured for the dataset.
	Stores []dszstore.DszStore `json:"stores"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// List Stores
// Lists all stores configured for a dataset.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/store/list",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszListStoresId struct{}

func (DszListStoresId) DszListStores(ctx context.Context, api *Api) (int, error) {
	panic("DszListStores not implemented: run 'go generate ./...'")
}
