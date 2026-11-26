//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package schedulerdeactivatemain

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/pipez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/pipez/relevancefashionpipez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type ReqParams struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type Req struct {
	Params  *ReqParams                   `json:"params,omitempty"`
	Payload *relevancefashionpipez.Main `json:"payload,omitempty"`
}

// Deactivate Relevance Fashion Main Pipeline Scheduler
// Deactivates a scheduler for the relevance fashion main pipeline on a dataset.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req `json:"req,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/pipe/dataset/{dataset}/scheduler/deactivate/relevance-fashion/relevance-fashion-main",
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
	if d.Req == nil || d.Req.Payload == nil {
		return &relevancefashionpipez.Main{}
	}
	return d.Req.Payload
}

type PipezSchedulerDeactivateRelevanceFashionRelevanceFashionMainId struct{}

func (PipezSchedulerDeactivateRelevanceFashionRelevanceFashionMainId) PipezSchedulerDeactivateRelevanceFashionRelevanceFashionMain(ctx context.Context, api *Api) (int, error) {
	panic("PipezSchedulerDeactivateRelevanceFashionRelevanceFashionMain not implemented: run 'go generate ./...'")
}
