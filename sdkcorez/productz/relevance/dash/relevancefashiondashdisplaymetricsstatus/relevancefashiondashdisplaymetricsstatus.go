//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package relevancefashiondashdisplaymetricsstatus

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
	// Source display identifier.
	SrcDisplayId string `json:"src_display_id"`
	// Number of raw recommendations.
	RecRaws float64 `json:"rec_raws"`
	// Number of recommendations shown.
	RecShows float64 `json:"rec_shows"`
	// Number of clicks.
	Clicks float64 `json:"clicks"`
	// Number of cart adds.
	CartAdds float64 `json:"cart_adds"`
	// Number of purchases.
	Purchases float64 `json:"purchases"`
	// Total revenue generated.
	Revenue float64 `json:"revenue"`
	// Unique sessions count.
	UniqueSessions float64 `json:"unique_sessions"`
	// ISO-8601 timestamp of last display activity.
	LastDisplayActivity string `json:"last_display_activity"`
	// Click-through rate.
	Ctr float64 `json:"ctr"`
	// Conversion rate.
	ConversionRate float64 `json:"conversion_rate"`
	// Performance uplift value.
	Uplift float64 `json:"uplift"`
	// Status indicator color for uplift.
	UpliftColor string `json:"uplift_color"`
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

// Get Display Metrics Status
// Gets recommendation metrics, revenue, session counts, and uplift status by display for a store and period.
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

type ProductRelevanceFashionDashDisplayMetricsStatusV1Id struct{}

func (ProductRelevanceFashionDashDisplayMetricsStatusV1Id) ProductRelevanceFashionDashDisplayMetricsStatusV1(ctx context.Context, fn *Func) (int, error) {
	panic("ProductRelevanceFashionDashDisplayMetricsStatusV1 not implemented: run 'go generate ./...'")
}
