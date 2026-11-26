//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package dszapipaymentautorenewupdate

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
	// Solution domain.
	Solution solution.Solution `json:"solution"`
}

type ReqPayload struct {
	// Auto-renew preference flag.
	AutoRenew bool `json:"auto_renew"`
}

type Req struct {
	Params  *ReqParams  `json:"params,omitempty"`
	Payload *ReqPayload `json:"payload,omitempty"`
}

// Update Payment Auto-Renew
// Updates the auto-renewal status for a dataset subscription solution.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req `json:"req,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/ds/dataset/{dataset}/solution/{solution}/payment/autorenew",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszUpdatePaymentAutoRenewId struct{}

func (DszUpdatePaymentAutoRenewId) DszUpdatePaymentAutoRenew(ctx context.Context, api *Api) (int, error) {
	panic("DszUpdatePaymentAutoRenew not implemented: run 'go generate ./...'")
}
