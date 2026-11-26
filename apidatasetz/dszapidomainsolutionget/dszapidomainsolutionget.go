//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package dszapidomainsolutionget

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/domain"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type SolutionDomain struct {
	Solution solution.Solution `json:"solution"`
	Domain   domain.Domain     `json:"domain"`
}

type ReqParams struct {
	// Dataset identifier.
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
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
	// Store ID queried.
	StoreId storeid.StoreId `json:"store_id"`
	// Solution domain mappings.
	Solutions []SolutionDomain `json:"solutions"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Get Store Solution Domain
// Retrieves the solution domain for a store within a dataset.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/store/solution-domain",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszDomainGetId struct{}

func (DszDomainGetId) DszDomainGet(ctx context.Context, api *Api) (int, error) {
	panic("DszDomainGet not implemented: run 'go generate ./...'")
}
