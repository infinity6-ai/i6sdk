//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package ingestzapistreamevent

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz/relevancefashiondataz"
)

type ReqParams struct {
	// Name of the dataset that receives the stream data.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type ReqPayload struct {
	// Telemetry events to ingest.
	Data []*relevancefashiondataz.Event `json:"data"`
}

type Req struct {
	Params  *ReqParams  `json:"params,omitempty"`
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// Number of ingested telemetry events.
	Count int `json:"count"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Ingest Relevance Fashion Event Stream
// Ingests a batch of telemetry events into the dataset event table.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/ingest/dataset/{dataset}/stream/relevance-fashion-event",
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

func (d *Api) PayloadDataSlice() []tabledataz.TableData {
	if d.Req == nil || d.Req.Payload == nil {
		return nil
	}
	res := make([]tabledataz.TableData, len(d.Req.Payload.Data))
	for i, item := range d.Req.Payload.Data {
		res[i] = item
	}
	return res
}

func (d *Api) SetRespCount(count int) {
	if d.Resp == nil {
		d.Resp = &Resp{}
	}
	if d.Resp.Payload == nil {
		d.Resp.Payload = &RespPayload{}
	}
	d.Resp.Payload.Count = count
}

type IngestzStreamRelevanceFashionEventId struct{}

func (IngestzStreamRelevanceFashionEventId) IngestzStreamRelevanceFashionEvent(ctx context.Context, api *Api) (int, error) {
	panic("IngestzStreamRelevanceFashionEvent not implemented: run 'go generate ./...'")
}
