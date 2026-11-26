//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package dszapiloadresultcatalog

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
	"github.com/infinity6-ai/i6sdk/sdkcorez/resultz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/resultz/relevancefashionresultz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type ReqParams struct {
	// Name of the dataset.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type ReqQuery struct {
	// Store ID to query.
	StoreId storeid.StoreId `json:"store_id"`
}

type Req struct {
	Params *ReqParams `json:"params,omitempty"`
	Query  *ReqQuery  `json:"query,omitempty"`
}

type RespPayload struct {
	// Loaded result data
	Data *relevancefashionresultz.CatalogSimple `json:"data"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Load Relevance Fashion Catalog Result
// Load a table file from the datalake as JSON and return it.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) Dataset() lakedsz.Dataset {
	if d.Req == nil || d.Req.Params == nil {
		return lakedsz.Dataset{}
	}
	return d.Req.Params.Dataset
}

func (d *Api) ReqRespData() resultz.Result {
	if d.Resp == nil {
		d.Resp = &Resp{Payload: &RespPayload{Data: &relevancefashionresultz.CatalogSimple{}}}
	} else if d.Resp.Payload == nil {
		d.Resp.Payload = &RespPayload{Data: &relevancefashionresultz.CatalogSimple{}}
	} else if d.Resp.Payload.Data == nil {
		d.Resp.Payload.Data = &relevancefashionresultz.CatalogSimple{}
	}
	if d.Req != nil && d.Req.Query != nil {
		d.Resp.Payload.Data.Req.StoreId = d.Req.Query.StoreId
	}
	return d.Resp.Payload.Data
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/results/relevance-fashion-catalog",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszResultRelevanceFashionCatalogId struct{}

func (DszResultRelevanceFashionCatalogId) DszResultRelevanceFashionCatalog(ctx context.Context, api *Api) (int, error) {
	panic("DszResultRelevanceFashionCatalog not implemented: run 'go generate ./...'")
}
