package relevancefashiondashimpactedmetricsbyshowcasecategory

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
	ShowcaseCategory *string
	Recomendations   float64
	Clicks           float64
	CartAdds         float64
	Purchases        float64
	Revenue          float64
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
		Id:      "product-relevance-fashion-dash-impacted-metrics-by-showcase-category",
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
											"showcase_category": {Field: func() any { return &(*element).ShowcaseCategory }},
											"recomendations":    {Field: func() any { return &(*element).Recomendations }},
											"clicks":            {Field: func() any { return &(*element).Clicks }},
											"cart_adds":         {Field: func() any { return &(*element).CartAdds }},
											"purchases":         {Field: func() any { return &(*element).Purchases }},
											"revenue":           {Field: func() any { return &(*element).Revenue }},
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
