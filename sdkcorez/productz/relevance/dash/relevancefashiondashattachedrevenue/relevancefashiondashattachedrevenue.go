//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package relevancefashiondashattachedrevenue

import (
	"context"

	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
)

type Params struct {
	// First year-month of period (e.g. 2026-07).
	MinYearMonth string `json:"min_year_month"`
	// Last year-month of period (e.g. 2026-11).
	MaxYearMonth string `json:"max_year_month"`
}

type Row struct {
	// Revenue attached to recommendations.
	AttachedRevenue float64 `json:"attached_revenue"`
	// Number of purchases attached to recommendations.
	AttachedPurchases float64 `json:"attached_purchases"`
	// Number of distinct SKUs attached to recommendations.
	AttachedSkus float64 `json:"attached_skus"`
	// Average ticket value for attached purchases.
	AvgTicket float64 `json:"avg_ticket"`
}

type Result struct {
	// Result rows.
	Rows []*Row `json:"rows"`
}

type ReqInput struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
	// Query params.
	Params *Params `json:"params,omitempty"`
}

type ReqPartitions struct {
	// Store identifier.
	StoreId string `json:"store_id"`
}

type Req struct {
	Input      *ReqInput      `json:"input,omitempty"`
	Partitions *ReqPartitions `json:"partitions,omitempty"`
}

type RespOutput struct {
	// Query results.
	Results Result `json:"results"`
}

type Resp struct {
	Output *RespOutput `json:"output,omitempty"`
}

// Get Attached Revenue
// Gets the revenue attached to recommendations for a store and period.
type Func struct {
	sdkfuncz.UnimplementedFunc `json:"-"`
	Req                        *Req  `json:"req,omitempty"`
	Resp                       *Resp `json:"resp,omitempty"`
}

func (d *Func) SdkFuncRoute() sdkfuncz.Spec {
	return sdkfuncz.Spec{
		Version: "v1",
		Typ:     sdkfuncz.TypClientQuery,
	}
}

type ProductRelevanceFashionDashAttachedRevenueV1Id struct{}

func (ProductRelevanceFashionDashAttachedRevenueV1Id) ProductRelevanceFashionDashAttachedRevenueV1(ctx context.Context, fn *Func) (int, error) {
	panic("ProductRelevanceFashionDashAttachedRevenueV1 not implemented: run 'go generate ./...'")
}
