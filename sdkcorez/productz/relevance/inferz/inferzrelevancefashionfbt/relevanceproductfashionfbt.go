//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package inferzrelevancefashionfbt

import (
	"context"

	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
)

type Sku struct {
	// SKU identifier.
	SkuId string `json:"sku_id"`
}

type Result struct {
	// Result rating.
	Rating float64 `json:"rating"`
	// Matching SKUs.
	Skus []*Sku `json:"skus"`
}

type ReqInput struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
	// Session identifier.
	SessionId string `json:"session_id,omitempty"`
	// SKUs to find matches for.
	Skus []*Sku `json:"skus"`
	// User identifier.
	UserId string `json:"user_id,omitempty"`
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
	// Inference results.
	Results []*Result `json:"results"`
}

type Resp struct {
	Output *RespOutput `json:"output,omitempty"`
}

// Get Bought Together
// Infers the SKUs frequently bought together with the given SKUs for a store.
type Func struct {
	sdkfuncz.UnimplementedFunc `json:"-"`
	Req                        *Req  `json:"req,omitempty"`
	Resp                       *Resp `json:"resp,omitempty"`
}

func (d *Func) SdkFuncRoute() sdkfuncz.Spec {
	return sdkfuncz.Spec{
		Version: "v1",
		Typ:     sdkfuncz.TypClientInfer,
	}
}

type ProductRelevanceFashionFbtV1Id struct{}

func (ProductRelevanceFashionFbtV1Id) ProductRelevanceFashionFbtV1(ctx context.Context, fn *Func) (int, error) {
	panic("ProductRelevanceFashionFbtV1 not implemented: run 'go generate ./...'")
}
