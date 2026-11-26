//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package dszapicustomerportalcreate

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type ReqParams struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type ReqPayload struct {
	// Return URL after customer portal session.
	ReturnUrl string `json:"return_url"`
	// Preferred locale for customer portal.
	Locale string `json:"locale,omitempty"`
}

type Req struct {
	Params  *ReqParams  `json:"params,omitempty"`
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// Stripe customer portal URL.
	Url string `json:"url"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Create Customer Portal
// Creates a customer portal session for managing subscription and billing.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/ds/dataset/{dataset}/portal",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszCreateCustomerPortalId struct{}

func (DszCreateCustomerPortalId) DszCreateCustomerPortal(ctx context.Context, api *Api) (int, error) {
	panic("DszCreateCustomerPortal not implemented: run 'go generate ./...'")
}
