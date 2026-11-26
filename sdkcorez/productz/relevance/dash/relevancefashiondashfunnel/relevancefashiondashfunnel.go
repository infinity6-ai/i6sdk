//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package relevancefashiondashfunnel

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
	// Event type name.
	EventType string `json:"event_type"`
	// Whether event originated from recommendation.
	FromRec string `json:"from_rec"`
	// Event count.
	QtEvents string `json:"qt_events"`
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

// Get Funnel
// Gets the event counts of the recommendation funnel for a store and period.
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

type ProductRelevanceFashionDashFunnelV1Id struct{}

func (ProductRelevanceFashionDashFunnelV1Id) ProductRelevanceFashionDashFunnelV1(ctx context.Context, fn *Func) (int, error) {
	panic("ProductRelevanceFashionDashFunnelV1 not implemented: run 'go generate ./...'")
}
