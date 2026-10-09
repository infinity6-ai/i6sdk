package inferzrelevancefashionfbt

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
)

type Sku struct {
	SkuId string
}

type Partitions struct {
	StoreId string
}

type Result struct {
	Skus   []*Sku
	Rating string
}

type Req struct {
	Dataset         lakedsz.Dataset
	Skus            []*Sku
	ModelPartitions Partitions
}

type Resp struct {
	Results []*Result
}

type Func struct {
	Req  Req
	Resp Resp
}

func (d *Func) SdkFuncSpec() sdkfuncz.Spec {
	return sdkfuncz.Spec{
		Id:      "product-relevance-fashion-fbt",
		Version: "v1",
		Typ:     sdkfuncz.TypClientInfer,
		Desc:    nil,
	}
}

func (d *Func) SdkFuncDataRefs() *sdkfuncz.DataRefs {
	return &sdkfuncz.DataRefs{
		Input: &schemaz.Schema{
			Object: func(_ bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {Str: lakedsz.ParseStr(&d.Req.Dataset)},
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
					"store_id": {Str: schemaz.ParseStr(&d.Req.ModelPartitions.StoreId)},
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
											"skus": {
												Array: schemaz.NewArray(schemaz.ArrayOptions[[]*Sku, *Sku]{
													Slice:      func() *[]*Sku { return &(*element).Skus },
													NewElement: func(idx int) *Sku { return &Sku{} },
													Get: func(idx int, ele **Sku) *schemaz.Schema {
														return &schemaz.Schema{
															Object: func(_ bool) map[string]*schemaz.Schema {
																return map[string]*schemaz.Schema{
																	"sku_id": {Str: schemaz.ParseStr(&(*ele).SkuId)},
																}
															},
														}
													},
												}),
											},
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
