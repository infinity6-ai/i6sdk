package dszapistripecheckout

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type LineItem struct {
	Price    string `json:"price"`
	Quantity int    `json:"quantity"`
}

func (l *LineItem) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"price": {
					Desc: &schemaz.Desc{Summary: "Price id"},
					Raw:  func() any { return &l.Price },
				},
				"quantity": {
					Desc: &schemaz.Desc{Summary: "Quantity"},
					Raw:  func() any { return &l.Quantity },
				},
			}
		},
	}
}

type Req struct {
	LineItems                       []*LineItem
	DatasetName                     string
	CustomerEmail                   string
	SuccessUrl                      string
	CancelUrl                       string
	SubscriptionDataTrialPeriodDays int
	CnpjHub                         string
	ResponsibleNameHub              string
	ResponsibleEmailHub             string
	IdSolutionHub                   string
	IdCampanhaHub                   string
	TermsVersion                    string
	ContractVersion                 string
	TermsAcceptedAt                 string
}

type Resp struct {
	Id  string
	Url string
}

type Api struct {
	Req  Req
	Resp Resp
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id: "dsz-stripe-checkout",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "Create Stripe Checkout",
				Summary: "Creates a Stripe checkout session for a subscription and returns its id and url.",
			}
		},
		Method: "POST",
		Path:   "/api/stripe/checkout",
	}
}

func (d *Api) ApiDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams:  nil,
		QueryParams: nil,
		ReqHeaders:  nil,
		ReqBody: &schemaz.Schema{
			Desc: &schemaz.Desc{Summary: "Checkout session parameters."},
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"line_items": {
						Desc: &schemaz.Desc{Summary: "Items being bought"},
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]*LineItem, *LineItem]{
							NewElement: func(idx int) *LineItem { return &LineItem{} },
							Slice:      func() *[]*LineItem { return &d.Req.LineItems },
							Get:        func(idx int, element **LineItem) *schemaz.Schema { return (*element).Schema() },
						}),
					},
					"dataset_name": {
						Desc: &schemaz.Desc{Summary: "Dataset name"},
						Raw:  func() any { return &d.Req.DatasetName },
					},
					"customer_email": {
						Desc: &schemaz.Desc{Summary: "Customer email"},
						Raw:  func() any { return &d.Req.CustomerEmail },
					},
					"success_url": {
						Desc: &schemaz.Desc{Summary: "Url to go after a successful checkout"},
						Raw:  func() any { return &d.Req.SuccessUrl },
					},
					"cancel_url": {
						Desc: &schemaz.Desc{Summary: "Url to go after a canceled checkout"},
						Raw:  func() any { return &d.Req.CancelUrl },
					},
					"subscription_data_trial_period_days": {
						Desc: &schemaz.Desc{Summary: "Trial period in days"},
						Raw:  func() any { return &d.Req.SubscriptionDataTrialPeriodDays },
					},
					"cnpj_hub": {
						Raw: func() any { return &d.Req.CnpjHub },
					},
					"responsible_name_hub": {
						Raw: func() any { return &d.Req.ResponsibleNameHub },
					},
					"responsible_email_hub": {
						Raw: func() any { return &d.Req.ResponsibleEmailHub },
					},
					"id_solution_hub": {
						Raw: func() any { return &d.Req.IdSolutionHub },
					},
					"id_campanha_hub": {
						Raw: func() any { return &d.Req.IdCampanhaHub },
					},
					"terms_version": {
						Raw: func() any { return &d.Req.TermsVersion },
					},
					"contract_version": {
						Raw: func() any { return &d.Req.ContractVersion },
					},
					"terms_accepted_at": {
						Raw: func() any { return &d.Req.TermsAcceptedAt },
					},
				}
			},
		},
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Desc: &schemaz.Desc{Summary: "The created checkout session."},
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"id": {
						Desc: &schemaz.Desc{Summary: "Checkout session id"},
						Raw:  func() any { return &d.Resp.Id },
					},
					"url": {
						Desc: &schemaz.Desc{Summary: "Checkout session url"},
						Raw:  func() any { return &d.Resp.Url },
					},
				}
			},
		},
	}
}
