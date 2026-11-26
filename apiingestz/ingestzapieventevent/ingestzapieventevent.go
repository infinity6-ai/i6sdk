//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package ingestzapieventevent

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
	// Name of the dataset that receives the event data.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type Req struct {
	Params  *ReqParams                  `json:"params,omitempty"`
	Payload *relevancefashiondataz.Event `json:"payload,omitempty"`
}

// Ingest Relevance Fashion Event
// Ingests telemetry event data into the dataset event table.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req `json:"req,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/ingest/dataset/{dataset}/event/relevance-fashion-event/text.txt",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
					"sinkable":             "true",
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

func (d *Api) PayloadData() tabledataz.TableData {
	if d.Req == nil {
		return nil
	}
	return d.Req.Payload
}

type IngestzEventRelevanceFashionEventId struct{}

func (IngestzEventRelevanceFashionEventId) IngestzEventRelevanceFashionEvent(ctx context.Context, api *Api) (int, error) {
	panic("IngestzEventRelevanceFashionEvent not implemented: run 'go generate ./...'")
}
