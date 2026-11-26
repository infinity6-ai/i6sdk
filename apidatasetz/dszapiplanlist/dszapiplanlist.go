//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package dszapiplanlist

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type Plan struct {
	Key                       string `json:"key"`
	Tier                      string `json:"tier"`
	Frequency                 string `json:"frequency"`
	Name                      string `json:"name"`
	Description               string `json:"description"`
	Currency                  string `json:"currency"`
	Amount                    int64  `json:"amount"`
	LicAmount                 int64  `json:"lic_amount"`
	ReqAmount                 int64  `json:"req_amount"`
	ReqAmountDecimal          string `json:"req_amount_decimal,omitempty"`
	IncludedRequests          *int64 `json:"included_requests,omitempty"`
	IncludedRequestsUnlimited bool   `json:"included_requests_unlimited"`
}

type ReqParams struct {
	// Solution product identifier.
	Solution solution.Solution `json:"solution"`
}

type Req struct {
	Params *ReqParams `json:"params,omitempty"`
}

type RespPayload struct {
	// Available subscription plans.
	Plans []Plan `json:"plans"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// List Subscription Plans
// Lists all available subscription pricing plans for a solution product.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/ds/solution/{solution}/plans",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszListPlansId struct{}

func (DszListPlansId) DszListPlans(ctx context.Context, api *Api) (int, error) {
	panic("DszListPlans not implemented: run 'go generate ./...'")
}
