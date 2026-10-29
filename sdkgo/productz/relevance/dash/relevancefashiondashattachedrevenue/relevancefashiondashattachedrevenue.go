package relevancefashiondashattachedrevenue

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
)

type Params struct {
	MinYearMonth string
	MaxYearMonth string
}

type Partitions struct {
	StoreId string
}

type Row struct {
	AttachedRevenue   float64
	AttachedPurchases float64
	AttachedSkus      float64
	AvgTicket         float64
}

type Result struct {
	Rows []*Row
}

type Req struct {
	Dataset        lakedsz.Dataset
	DashPartitions Partitions
	Params         Params
}

type Resp struct {
	Results Result
}

type Func struct {
	Req  Req
	Resp Resp
}

func (d *Func) SdkFuncSpec() sdkfuncz.Spec {
	return sdkfuncz.Spec{
		Id:      "product-relevance-fashion-dash-attached-revenue",
		Version: "v1",
		Typ:     sdkfuncz.TypClientQuery,
		Desc:    nil,
	}
}

func (d *Func) SdkFuncDataRefs() *sdkfuncz.DataRefs {
	return &sdkfuncz.DataRefs{
		Input: &schemaz.Schema{
			Object: func(_ bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {Field: func() any { return &d.Req.Dataset }},
					"params": {
						Object: func(_ bool) map[string]*schemaz.Schema {
							return map[string]*schemaz.Schema{
								"min_year_month": {Field: func() any { return &d.Req.Params.MinYearMonth }},
								"max_year_month": {Field: func() any { return &d.Req.Params.MaxYearMonth }},
							}
						},
					},
				}
			},
		},
		Partitions: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"store_id": {Field: func() any { return &d.Req.DashPartitions.StoreId }},
				}
			},
		},
		Output: &schemaz.Schema{
			Object: func(_ bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"rows": {
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]*Row, *Row]{
							Slice:      func() *[]*Row { return &d.Resp.Results.Rows },
							NewElement: func(idx int) *Row { return &Row{} },
							Get: func(idx int, element **Row) *schemaz.Schema {
								return &schemaz.Schema{
									Object: func(_ bool) map[string]*schemaz.Schema {
										return map[string]*schemaz.Schema{
											"attached_revenue":   {Field: func() any { return &(*element).AttachedRevenue }},
											"attached_purchases": {Field: func() any { return &(*element).AttachedPurchases }},
											"attached_skus":      {Field: func() any { return &(*element).AttachedSkus }},
											"avg_ticket":         {Field: func() any { return &(*element).AvgTicket }},
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
