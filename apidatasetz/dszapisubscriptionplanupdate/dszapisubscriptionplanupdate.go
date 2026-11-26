//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package dszapisubscriptionplanupdate

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
	// Solution product.
	Solution solution.Solution `json:"solution"`
}

type ReqPayload struct {
	// Lookup key of new Stripe price.
	Plan string `json:"plan"`
}

type Req struct {
	Params  *ReqParams  `json:"params,omitempty"`
	Payload *ReqPayload `json:"payload,omitempty"`
}

// Update Subscription Plan
// Updates the plan (price) of the Stripe subscription for a dataset product/solution.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req `json:"req,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/ds/dataset/{dataset}/solution/{solution}/subscription/plan",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszUpdateSubscriptionPlanId struct{}

func (DszUpdateSubscriptionPlanId) DszUpdateSubscriptionPlan(ctx context.Context, api *Api) (int, error) {
	panic("DszUpdateSubscriptionPlan not implemented: run 'go generate ./...'")
}
