package relevancefashiondashdisplaymetricsstatus

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
)

// Package relevancefashiondashdisplaymetricsstatus implements the client query
// product-relevance-fashion-dash-display-metrics-status: recommendation funnel
// metrics and status indicators grouped by display (showcase) for one store and
// a year-month period.
//
// The SDK only declares the input and output schemas. The filters below are
// applied by the query itself, which receives the params exactly as sent.

// Special values used in Row.SrcDisplayId. Any other value is a display id.
const (
	// DisplayNoDisplay identifies displays that are not registered.
	DisplayNoDisplay = "NO_DISPLAY"
)

// Uplift indicator color values returned in Row.UpliftColor.
const (
	UpliftColorGreen  = "GREEN"
	UpliftColorYellow = "YELLOW"
	UpliftColorRed    = "RED"
)

// Params are the query params.
type Params struct {
	// MinYearMonth is the first year-month of the period, for example 2026-07.
	MinYearMonth string
	// MaxYearMonth is the last year-month of the period, for example 2026-11.
	MaxYearMonth string
}

// Partitions are the dash partitions.
type Partitions struct {
	// StoreId is the store the metrics belong to.
	StoreId string
}

// Row is the metrics and status of one display.
type Row struct {
	// SrcDisplayId is the source display id, or DisplayNoDisplay ("NO_DISPLAY")
	// for displays that are not registered.
	SrcDisplayId string
	// RecRaws is the number of raw recommendations.
	RecRaws float64
	// RecShows is the number of recommendations shown.
	RecShows float64
	// Clicks is the number of clicks.
	Clicks float64
	// CartAdds is the number of cart adds.
	CartAdds float64
	// Purchases is the number of purchases.
	Purchases float64
	// Revenue is the total revenue generated.
	Revenue float64
	// UniqueSessions is the count of unique sessions.
	UniqueSessions float64
	// LastDisplayActivity is the ISO-8601 timestamp of the last display activity.
	LastDisplayActivity string
	// Ctr is the click-through rate.
	Ctr float64
	// ConversionRate is purchases divided by clicks.
	ConversionRate float64
	// Uplift is the performance uplift value.
	Uplift float64
	// UpliftColor is the status indicator color for the uplift (e.g., GREEN, YELLOW, RED).
	UpliftColor string
}

// Result is the query result.
type Result struct {
	// Rows has one row per display.
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

// Func is the client query product-relevance-fashion-dash-display-metrics-status.
type Func struct {
	Req  Req
	Resp Resp
}

func (d *Func) SdkFuncSpec() sdkfuncz.Spec {
	return sdkfuncz.Spec{
		Id:      "product-relevance-fashion-dash-display-metrics-status",
		Version: "v1",
		Typ:     sdkfuncz.TypClientQuery,
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "Get Display Metrics Status",
				Summary: "Gets recommendation metrics, revenue, session counts, and uplift status by display for a store and period.",
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
												Desc:  &schemaz.Desc{Summary: "Display id, or NO_DISPLAY if not registered."},
												Field: func() any { return &(*element).SrcDisplayId },
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
											"revenue": {
												Desc:  &schemaz.Desc{Summary: "Total revenue generated."},
												Field: func() any { return &(*element).Revenue },
											},
											"unique_sessions": {
												Desc:  &schemaz.Desc{Summary: "Number of unique sessions."},
												Field: func() any { return &(*element).UniqueSessions },
											},
											"last_display_activity": {
												Desc:  &schemaz.Desc{Summary: "ISO-8601 timestamp of the last display activity."},
												Field: func() any { return &(*element).LastDisplayActivity },
											},
											"ctr": {
												Desc:  &schemaz.Desc{Summary: "Click-through rate."},
												Field: func() any { return &(*element).Ctr },
											},
											"conversion_rate": {
												Desc:  &schemaz.Desc{Summary: "Purchases divided by clicks."},
												Field: func() any { return &(*element).ConversionRate },
											},
											"uplift": {
												Desc:  &schemaz.Desc{Summary: "Performance uplift metric."},
												Field: func() any { return &(*element).Uplift },
											},
											"uplift_color": {
												Desc:  &schemaz.Desc{Summary: "Uplift status color indicator."},
												Field: func() any { return &(*element).UpliftColor },
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
