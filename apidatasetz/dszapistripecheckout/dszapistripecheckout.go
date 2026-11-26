//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package dszapistripecheckout

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
)

type ReqPayload struct {
	// Dataset name.
	DatasetName string `json:"dataset_name"`
	// Plan key to buy.
	Plan string `json:"plan"`
	// URL to redirect after successful checkout.
	SuccessUrl string `json:"success_url"`
	// URL to redirect after canceled checkout.
	CancelUrl string `json:"cancel_url"`
	// Base64 state token with checkout metadata.
	CampaignState string `json:"campaign_state,omitempty"`
	// Phone of buyer.
	PaymentPhone string `json:"payment_phone,omitempty"`
	// Email of buyer in Stripe.
	PaymentEmail string `json:"payment_email,omitempty"`
	// Document/tax ID of buyer.
	Document string `json:"document,omitempty"`
	// Optional callback URL for redeeming checkout artifact ticket.
	UserExchangeCallbackUrl string `json:"user_exchange_callback_url,omitempty"`
}

type Req struct {
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// Checkout session ID.
	Id string `json:"id"`
	// Checkout session URL.
	Url string `json:"url"`
	// Dataset identifier created after payment.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Create Stripe Checkout
// Creates a Stripe checkout session for a subscription and returns its ID and URL.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/stripe/checkout",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{}
		},
	}
}

type DszStripeCheckoutId struct{}

func (DszStripeCheckoutId) DszStripeCheckout(ctx context.Context, api *Api) (int, error) {
	panic("DszStripeCheckout not implemented: run 'go generate ./...'")
}
