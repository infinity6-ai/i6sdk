//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package relevancefashiondashimpactedmetricsbyshowcasecategory

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
	// Showcase category name.
	ShowcaseCategory *string `json:"showcase_category"`
	// Number of recommendations.
	Recomendations float64 `json:"recomendations"`
	// Number of clicks.
	Clicks float64 `json:"clicks"`
	// Number of cart adds.
	CartAdds float64 `json:"cart_adds"`
	// Number of purchases.
	Purchases float64 `json:"purchases"`
	// Total revenue generated.
	Revenue float64 `json:"revenue"`
	// Participation rate.
	Participation float64 `json:"participation"`
	// Conversion rate.
	ConversionRate float64 `json:"conversion_rate"`
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

// Get Metrics by Category
// Gets the metrics impacted by recommendations by showcase category for a store and period.
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

type ProductRelevanceFashionDashImpactedMetricsByShowcaseCategoryV1Id struct{}

func (ProductRelevanceFashionDashImpactedMetricsByShowcaseCategoryV1Id) ProductRelevanceFashionDashImpactedMetricsByShowcaseCategoryV1(ctx context.Context, fn *Func) (int, error) {
	panic("ProductRelevanceFashionDashImpactedMetricsByShowcaseCategoryV1 not implemented: run 'go generate ./...'")
}
