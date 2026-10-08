package inferzrelevancefashionfbt

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
)

type Sku struct {
	SkuId string `json:"sku_id"`
}

type Result struct {
	SkuId  string `json:"sku_id"`
	Rating string `json:"rating"`
}

type Req struct {
	Skus []*Sku `json:"skus"`
}

type Resp struct {
	Results []*Result `json:"results"`
}

type Api struct {
	Req  *Req  `json:"req"`
	Resp *Resp `json:"resp"`
}

func (d *Api) InferSpec() sdkfuncz.InferSpec {
	return sdkfuncz.InferSpec{
		Id:       "product-relevance-fashion-fbt",
		Typ:      sdkfuncz.TypClientInfer,
		FuncName: "fashionFbt",
		Desc:     nil,
	}
}

func (d *Api) InferDataRefs() *sdkfuncz.DataRefs {
	return &sdkfuncz.DataRefs{
		Input: &schemaz.Schema{
			Object: func(_ bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"skus": {
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]*Sku, *Sku]{
							Get: func(idx int, element **Sku) *schemaz.Schema {
								return &schemaz.Schema{
									Object: func(_ bool) map[string]*schemaz.Schema {
										return map[string]*schemaz.Schema{
											"sku_id": {Str: schemaz.ParseStr(&(*element).SkuId)},
										}
									},
								}
							},
						}),
					},
				}
			},
		},
		Output: &schemaz.Schema{
			Object: func(_ bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"results": {
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]*Result, *Result]{
							Get: func(idx int, element **Result) *schemaz.Schema {
								return &schemaz.Schema{
									Object: func(_ bool) map[string]*schemaz.Schema {
										return map[string]*schemaz.Schema{
											"sku_id": {Str: schemaz.ParseStr(&(*element).SkuId)},
											"rating": {Str: schemaz.ParseStr(&(*element).Rating)},
										}
									},
								}
							},
						}),
					},
				}
			},
		},
	}
}
