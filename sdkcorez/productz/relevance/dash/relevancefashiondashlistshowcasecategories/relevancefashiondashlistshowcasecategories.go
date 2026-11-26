//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package relevancefashiondashlistshowcasecategories

import (
	"context"

	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
)

type Row struct {
	// Showcase category name.
	ShowcaseCategory *string `json:"showcase_category"`
}

type Result struct {
	// Result rows.
	Rows []*Row `json:"rows"`
}

type ReqInput struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
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

// List Showcase Categories
// Lists the showcase categories of a store.
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

type ProductRelevanceFashionDashListShowcaseCategoriesV1Id struct{}

func (ProductRelevanceFashionDashListShowcaseCategoriesV1Id) ProductRelevanceFashionDashListShowcaseCategoriesV1(ctx context.Context, fn *Func) (int, error) {
	panic("ProductRelevanceFashionDashListShowcaseCategoriesV1 not implemented: run 'go generate ./...'")
}
