//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package dszapiloadresultbatchsummary

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbnamez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
	"github.com/infinity6-ai/i6sdk/sdkcorez/resultz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/resultz/platformresultz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type ReqParams struct {
	// Name of the dataset.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type ReqQuery struct {
	// Table of the batch in question.
	Table laketbnamez.TableName `json:"table"`
	// Store ID to query.
	StoreId storeid.StoreId `json:"store_id"`
}

type Req struct {
	Params *ReqParams `json:"params,omitempty"`
	Query  *ReqQuery  `json:"query,omitempty"`
}

type RespPayload struct {
	// Loaded result data
	Data *platformresultz.BatchSummary `json:"data"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Load Last Batch Summary Result
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
		d.Resp = &Resp{Payload: &RespPayload{Data: &platformresultz.BatchSummary{}}}
	} else if d.Resp.Payload == nil {
		d.Resp.Payload = &RespPayload{Data: &platformresultz.BatchSummary{}}
	} else if d.Resp.Payload.Data == nil {
		d.Resp.Payload.Data = &platformresultz.BatchSummary{}
	}
	if d.Req != nil && d.Req.Query != nil {
		if d.Resp.Payload.Data.Req == nil {
			d.Resp.Payload.Data.Req = &platformresultz.BatchSummaryReq{}
		}
		d.Resp.Payload.Data.Req.Table = d.Req.Query.Table
		d.Resp.Payload.Data.Req.StoreId = d.Req.Query.StoreId
	}
	return d.Resp.Payload.Data
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/results/last-batch-summary",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszResultLastBatchSummaryId struct{}

func (DszResultLastBatchSummaryId) DszResultLastBatchSummary(ctx context.Context, api *Api) (int, error) {
	panic("DszResultLastBatchSummary not implemented: run 'go generate ./...'")
}
