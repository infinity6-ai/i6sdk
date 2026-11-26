//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package dszapidatasetsolutions

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
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
	// Active solutions enabled on the dataset.
	Solutions []solution.Solution `json:"solutions"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// List Active Dataset Solutions
// Lists all active solution domains enabled for a dataset.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/solutions",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszDatasetSolutionsId struct{}

func (DszDatasetSolutionsId) DszDatasetSolutions(ctx context.Context, api *Api) (int, error) {
	panic("DszDatasetSolutions not implemented: run 'go generate ./...'")
}
