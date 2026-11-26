//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package pipezlist

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solutiondomain"
	"github.com/infinity6-ai/i6sdk/sdkcorez/pipez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type ReqParams struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type ReqQuery struct {
	// Filter by solution domain.
	SolutionDomain solutiondomain.SolutionDomain `json:"solution_domain,omitempty"`
	// Filter by pipeline name.
	Name string `json:"name,omitempty"`
	// Filter by pipeline status.
	Status string `json:"status,omitempty"`
	// Cursor to fetch the next page.
	Cursor string `json:"cursor,omitempty"`
	// Maximum number of results.
	Limit int `json:"limit,omitempty"`
}

type Req struct {
	Params *ReqParams `json:"params,omitempty"`
	Query  *ReqQuery  `json:"query,omitempty"`
}

// Response containing the listed pipelines.
type RespPayload struct {
	// List of pipelines.
	Pipelines []*pipez.Pipeline `json:"pipelines"`
	// Cursor for the next page.
	Cursor string `json:"cursor,omitempty"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// List Pipelines
// Lists the pipelines of a dataset, optionally filtered, with cursor pagination.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/pipe/dataset/{dataset}/list",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type PipezListId struct{}

func (PipezListId) PipezList(ctx context.Context, api *Api) (int, error) {
	panic("PipezList not implemented: run 'go generate ./...'")
}
