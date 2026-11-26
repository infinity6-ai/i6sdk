//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package dszapistoredisplaycreate

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

type ReqPayload struct {
	// Store ID.
	StoreId storeid.StoreId `json:"store_id"`
	// Store display details.
	Display dszstore.DszStoreDisplay `json:"display"`
}

type Req struct {
	Params  *ReqParams  `json:"params,omitempty"`
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// Created store display details.
	Display dszstore.DszStoreDisplay `json:"display"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Create Store Display
// Creates or updates a store display configuration.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/ds/dataset/{dataset}/store/display/create",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszStoreDisplayCreateId struct{}

func (DszStoreDisplayCreateId) DszStoreDisplayCreate(ctx context.Context, api *Api) (int, error) {
	panic("DszStoreDisplayCreate not implemented: run 'go generate ./...'")
}
