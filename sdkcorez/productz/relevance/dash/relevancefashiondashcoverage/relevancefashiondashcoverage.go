//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package relevancefashiondashcoverage

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
	// Catalog recommendation coverage percentage.
	RecomendationCoverage float64 `json:"recomendation_coverage"`
	// Recommended items count.
	RecommendedItems float64 `json:"recommended_items"`
	// Recommendable items count.
	RecommendableItems float64 `json:"recommendable_items"`
	// Long-tail coverage percentage.
	LongTailCoverage float64 `json:"long_tail_coverage"`
	// Top 1 coverage percentage.
	Top1Coverage float64 `json:"top_1_coverage"`
	// Top 1 to 5 coverage percentage.
	Top1To5Coverage float64 `json:"top_1_to_5_coverage"`
	// Top 5 to 10 coverage percentage.
	Top5To10Coverage float64 `json:"top_5_to_10_coverage"`
	// Bottom 90 coverage percentage.
	Bottom90Coverage float64 `json:"bottom_90_coverage"`
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

// Get Coverage
// Gets the recommendation coverage of the catalog for a store and period.
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

type ProductRelevanceFashionDashCoverageV1Id struct{}

func (ProductRelevanceFashionDashCoverageV1Id) ProductRelevanceFashionDashCoverageV1(ctx context.Context, fn *Func) (int, error) {
	panic("ProductRelevanceFashionDashCoverageV1 not implemented: run 'go generate ./...'")
}
