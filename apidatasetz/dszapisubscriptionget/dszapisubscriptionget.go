//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package dszapisubscriptionget

import (
	"context"
	"time"

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

type Req struct {
	Params *ReqParams `json:"params,omitempty"`
}

type RespPayload struct {
	// Lookup key of base price in Stripe.
	PlanKey string `json:"plan_key"`
	// Subscription status in Stripe.
	Status string `json:"status"`
	// Cancellation timestamp in UTC.
	CancelAt *time.Time `json:"cancel_at,omitempty"`
	// Subscription creation timestamp in UTC.
	CreatedAt time.Time `json:"created_at"`
	// Whether subscription will be canceled at period end.
	CancelAtPeriodEnd bool `json:"cancel_at_period_end"`
	// Trial start timestamp in UTC.
	TrialStartAt *time.Time `json:"trial_start_at,omitempty"`
	// Trial end timestamp in UTC.
	TrialEndAt *time.Time `json:"trial_end_at,omitempty"`
	// Current billing period start timestamp in UTC.
	CurrentPeriodStartAt time.Time `json:"current_period_start_at"`
	// Current billing period end timestamp in UTC.
	CurrentPeriodEndAt time.Time `json:"current_period_end_at"`
	// Next invoice timestamp in UTC.
	NextInvoiceAt *time.Time `json:"next_invoice_at,omitempty"`
}

// Custom tells that the dataset has no stripe subscription
func (r *RespPayload) Custom() bool {
	return r.Status == ""
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

func (r *Resp) Custom() bool {
	return r == nil || r.Payload == nil || r.Payload.Custom()
}

// Get Subscription
// Gets the plan and status of the Stripe subscription for a dataset solution.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/solution/{solution}/subscription",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszGetSubscriptionId struct{}

func (DszGetSubscriptionId) DszGetSubscription(ctx context.Context, api *Api) (int, error) {
	panic("DszGetSubscription not implemented: run 'go generate ./...'")
}
