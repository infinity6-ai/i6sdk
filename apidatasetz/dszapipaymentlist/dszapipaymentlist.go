//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package dszapipaymentlist

import (
	"context"
	"time"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type PaymentLine struct {
	Id          string            `json:"id"`
	Solution    solution.Solution `json:"solution"`
	Amount      int64             `json:"amount"`
	Description string            `json:"description"`
}

type Payment struct {
	Id            string        `json:"id"`
	Number        string        `json:"number"`
	Status        string        `json:"status"`
	Currency      string        `json:"currency"`
	AmountDue     int64         `json:"amount_due"`
	AmountPaid    int64         `json:"amount_paid"`
	Created       time.Time     `json:"created"`
	PeriodStart   time.Time     `json:"period_start"`
	PeriodEnd     time.Time     `json:"period_end"`
	InvoicePdfUrl string        `json:"invoice_pdf_url"`
	ReceiptPdfUrl string        `json:"receipt_pdf_url"`
	Lines         []PaymentLine `json:"lines"`
}

type ReqParams struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type Req struct {
	Params *ReqParams `json:"params,omitempty"`
}

type RespPayload struct {
	// Payment invoice records.
	Payments []Payment `json:"payments"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// List Payments
// Lists all Stripe payment records for a dataset.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/payments",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszListPaymentsId struct{}

func (DszListPaymentsId) DszListPayments(ctx context.Context, api *Api) (int, error) {
	panic("DszListPayments not implemented: run 'go generate ./...'")
}
