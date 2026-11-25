//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package pipezstartbatchingest

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbnamez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/pipez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/pipez/relevancefashionpipez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

// Params contains path parameters embedded in the route URL.
type Params struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
}

// Payload contains pipeline parameters for batch ingestion.
type Payload struct {
	// Batch ingestion id.
	IngestId string `json:"ingest_id"`
	// Table that received the batch.
	Table laketbnamez.TableName `json:"table"`
}

// RespPayload contains the execution result of starting the pipeline.
type RespPayload struct {
	// Pipeline run identifier.
	Id string `json:"id"`
}

// Req encapsulates all incoming HTTP data for the endpoint.
type Req struct {
	Params  *Params  `json:"params,omitempty"`
	Payload *Payload `json:"payload,omitempty"`
}

// Resp encapsulates all outgoing HTTP data for the endpoint.
type Resp struct {
	Payload *RespPayload `json:"payload"`
}

// Start Batch Ingestion Pipeline
// Starts a batch ingestion pipeline on a dataset with the given params and returns its id.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req"`
	Resp                  *Resp `json:"resp"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "pipez-start-relevance_fashion-relevance-fashion-batch-ingest",
		Method: "POST",
		Path:   "/api/pipe/dataset/{dataset}/start/relevance_fashion/relevance-fashion-batch-ingest",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

func (d *Api) Dataset() lakedsz.Dataset {
	if d.Req == nil || d.Req.Params == nil {
		return lakedsz.Dataset{}
	}
	return d.Req.Params.Dataset
}

func (d *Api) Pipe() pipez.Pipe {
	p := &relevancefashionpipez.BatchIngest{}
	if d.Req != nil && d.Req.Payload != nil {
		p.IngestId = d.Req.Payload.IngestId
		p.Table = d.Req.Payload.Table
	}
	return p
}

func (d *Api) SetRespId(id string) {
	if d.Resp == nil {
		d.Resp = &Resp{}
	}
	if d.Resp.Payload == nil {
		d.Resp.Payload = &RespPayload{}
	}
	d.Resp.Payload.Id = id
}
