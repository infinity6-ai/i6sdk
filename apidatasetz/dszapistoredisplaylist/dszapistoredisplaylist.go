//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package dszapistoredisplaylist

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/dszstore"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type ReqParams struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type ReqQuery struct {
	// Store ID to filter displays.
	StoreId storeid.StoreId `json:"store_id"`
}

type Req struct {
	Params *ReqParams `json:"params,omitempty"`
	Query  *ReqQuery  `json:"query,omitempty"`
}

type RespPayload struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
	// Store ID queried.
	StoreId storeid.StoreId `json:"store_id"`
	// Displays configured for the store.
	Displays []dszstore.DszStoreDisplay `json:"displays"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// List Store Displays
// Lists displays configured for a store in a dataset.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/store/display/list",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszStoreDisplayListId struct{}

func (DszStoreDisplayListId) DszStoreDisplayList(ctx context.Context, api *Api) (int, error) {
	panic("DszStoreDisplayList not implemented: run 'go generate ./...'")
}
