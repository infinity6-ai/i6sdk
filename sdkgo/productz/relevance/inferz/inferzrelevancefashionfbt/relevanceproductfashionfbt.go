package inferzrelevancefashionfbt

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
)

type Sku struct {
	SkuId string `json:"sku_id"`
}

type Partitions struct {
	Version string `json:"version"`
	StoreId string `json:"store_id"`
}

type Result struct {
	SkuId  string `json:"sku_id"`
	Rating string `json:"rating"`
}

type Req struct {
	Skus       []*Sku     `json:"skus"`
	Partitions Partitions `json:"partitions"`
}

type Resp struct {
	Results []*Result `json:"results"`
}

type Func struct {
	Req  Req  `json:"req"`
	Resp Resp `json:"resp"`
}

func (d *Func) SdkFuncSpec() sdkfuncz.Spec {
	return sdkfuncz.Spec{
		Id:       "product-relevance-fashion-fbt",
		Typ:      sdkfuncz.TypClientInfer,
		FuncName: "fashionFbt",
		Desc:     nil,
	}
}

func (d *Func) SdkFuncDataRefs() *sdkfuncz.DataRefs {
	return &sdkfuncz.DataRefs{
		Input: &schemaz.Schema{
			Object: func(_ bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"skus": {
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]*Sku, *Sku]{
							Slice:      func() *[]*Sku { return &d.Req.Skus },
							NewElement: func(idx int) *Sku { return &Sku{} },
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
		Partitions: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"version":  {Str: schemaz.ParseStr(&d.Req.Partitions.Version)},
					"store_id": {Str: schemaz.ParseStr(&d.Req.Partitions.StoreId)},
				}
			},
		},
		Output: &schemaz.Schema{
			Object: func(_ bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"results": {
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]*Result, *Result]{
							Slice:      func() *[]*Result { return &d.Resp.Results },
							NewElement: func(idx int) *Result { return &Result{} },
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
