package relevancefashiondashdisplaymetrics

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
)

// Package relevancefashiondashdisplaymetrics implements the client query
// product-relevance-fashion-dash-display-metrics: recommendation funnel metrics
// grouped by display (showcase) and by recommendation origin, for one store and
// a year-month period.
//
// The SDK only declares the input and output schemas. The filters below are
// applied by the query itself, which receives the params exactly as sent.

// Special values accepted by Params.DisplayId. Any other value is a display id.
const (
	// DisplayNoDisplay returns only rows whose display id is null, which are
	// the displays that are not registered.
	DisplayNoDisplay = "NO_DISPLAY"
	// DisplayAllDisplays returns rows of every display, registered or not.
	DisplayAllDisplays = "ALL_DISPLAYS"
)

// Params are the query params. Optional ones are pointers so that absent and
// null (nil) can be told apart from an explicit value.
type Params struct {
	// MinYearMonth is the first year-month of the period, for example 2026-07.
	MinYearMonth string
	// MaxYearMonth is the last year-month of the period, for example 2026-09.
	MaxYearMonth string
	// FromRec is an optional tri-state filter on the recommendation origin:
	//   nil (absent or null): returns all rows, recommended or not.
	//   true: returns only rows that come from recommendations.
	//   false: returns only rows that do not come from recommendations.
	FromRec *bool
	// DisplayId is an optional display filter, passed to the query as sent:
	//   nil (absent or null): no display filter is requested.
	//   DisplayNoDisplay: only rows whose display id is null.
	//   DisplayAllDisplays: rows of all displays.
	//   any other value: only rows of that display id.
	DisplayId *string
}

// Partitions are the dash partitions.
type Partitions struct {
	// StoreId is the store the metrics belong to.
	StoreId string
}

// Row is the metrics of one display and recommendation origin.
type Row struct {
	// SrcDisplayId is the source display id. It is nil for displays that are
	// not registered.
	SrcDisplayId *string
	// FromRec is 1 when the row counts recommendation traffic and 0 otherwise.
	FromRec int
	// RecRaws is the number of raw recommendations. It is 0 when FromRec is 0.
	RecRaws float64
	// RecShows is the number of recommendations shown. It is 0 when FromRec is 0.
	RecShows float64
	// Clicks is the number of clicks.
	Clicks float64
	// CartAdds is the number of cart adds.
	CartAdds float64
	// Purchases is the number of purchases.
	Purchases float64
	// ConversionRate is purchases divided by clicks.
	ConversionRate float64
}

// Result is the query result.
type Result struct {
	// Rows has one row per display and recommendation origin.
	Rows []*Row
}

// Req is the request: the dataset, the store partition and the params.
type Req struct {
	Dataset        lakedsz.Dataset
	DashPartitions Partitions
	Params         Params
}

// Resp is the response.
type Resp struct {
	Results Result
}

// Func is the client query product-relevance-fashion-dash-display-metrics.
type Func struct {
	Req  Req
	Resp Resp
}

func (d *Func) SdkFuncSpec() sdkfuncz.Spec {
	return sdkfuncz.Spec{
		Id:      "product-relevance-fashion-dash-display-metrics",
		Version: "v1",
		Typ:     sdkfuncz.TypClientQuery,
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "Get Display Metrics",
				Summary: "Gets recommendation metrics by display for a store and period, optionally filtered by display and recommendation origin.",
			}
		},
	}
}

func (d *Func) SdkFuncDataRefs() *sdkfuncz.DataRefs {
	return &sdkfuncz.DataRefs{
		Input: &schemaz.Schema{
			Object: func(_ bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {
						Desc:  &schemaz.Desc{Summary: "Dataset identifier."},
						Field: func() any { return &d.Req.Dataset },
					},
					"params": {
						Desc: &schemaz.Desc{Summary: "Query params."},
						Object: func(_ bool) map[string]*schemaz.Schema {
							return map[string]*schemaz.Schema{
								"min_year_month": {
									Desc:  &schemaz.Desc{Summary: "First year-month of the period."},
									Field: func() any { return &d.Req.Params.MinYearMonth },
								},
								"max_year_month": {
									Desc:  &schemaz.Desc{Summary: "Last year-month of the period."},
									Field: func() any { return &d.Req.Params.MaxYearMonth },
								},
								"from_rec": {
									Desc:  &schemaz.Desc{Summary: "Recommendation filter: true, false, or null for all."},
									Field: func() any { return &d.Req.Params.FromRec },
								},
								"display_id": {
									Desc:  &schemaz.Desc{Summary: "Display filter: NO_DISPLAY, ALL_DISPLAYS or a display id."},
									Field: func() any { return &d.Req.Params.DisplayId },
								},
							}
						},
					},
				}
			},
		},
		Partitions: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"store_id": {
						Desc:  &schemaz.Desc{Summary: "Store identifier."},
						Field: func() any { return &d.Req.DashPartitions.StoreId },
					},
				}
			},
		},
		Output: &schemaz.Schema{
			Object: func(_ bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"rows": {
						Desc: &schemaz.Desc{Summary: "Result rows."},
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]*Row, *Row]{
							Slice:      func() *[]*Row { return &d.Resp.Results.Rows },
							NewElement: func(idx int) *Row { return &Row{} },
							Get: func(idx int, element **Row) *schemaz.Schema {
								return &schemaz.Schema{
									Object: func(_ bool) map[string]*schemaz.Schema {
										return map[string]*schemaz.Schema{
											"src_display_id": {
												Desc:  &schemaz.Desc{Summary: "Display id, null if not registered."},
												Field: func() any { return &(*element).SrcDisplayId },
											},
											"from_rec": {
												Desc:  &schemaz.Desc{Summary: "1 if from recommendations, else 0."},
												Field: func() any { return &(*element).FromRec },
											},
											"recraws": {
												Desc:  &schemaz.Desc{Summary: "Raw recommendations count."},
												Field: func() any { return &(*element).RecRaws },
											},
											"recshows": {
												Desc:  &schemaz.Desc{Summary: "Recommendations shown count."},
												Field: func() any { return &(*element).RecShows },
											},
											"clicks": {
												Desc:  &schemaz.Desc{Summary: "Number of clicks."},
												Field: func() any { return &(*element).Clicks },
											},
											"cart_adds": {
												Desc:  &schemaz.Desc{Summary: "Number of cart adds."},
												Field: func() any { return &(*element).CartAdds },
											},
											"purchases": {
												Desc:  &schemaz.Desc{Summary: "Number of purchases."},
												Field: func() any { return &(*element).Purchases },
											},
											"conversion_rate": {
												Desc:  &schemaz.Desc{Summary: "Purchases divided by clicks."},
												Field: func() any { return &(*element).ConversionRate },
											},
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
