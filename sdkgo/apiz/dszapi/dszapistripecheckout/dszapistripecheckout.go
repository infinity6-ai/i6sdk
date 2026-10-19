package dszapistripecheckout

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type Req struct {
	User        string
	DatasetName string
	Plan        string
	SuccessUrl  string
	CancelUrl   string
	State       string
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
					"user": {
						Desc: &schemaz.Desc{Summary: "User that is buying"},
						Raw:  func() any { return &d.Req.User },
					},
					"dataset_name": {
						Desc: &schemaz.Desc{Summary: "Dataset name"},
						Raw:  func() any { return &d.Req.DatasetName },
					},
					"plan": {
						Desc: &schemaz.Desc{Summary: "Lookup key of the Stripe price of the plan to buy; the service gets the price id from the Stripe prices api"},
						Raw:  func() any { return &d.Req.Plan },
					},
					"success_url": {
						Desc: &schemaz.Desc{Summary: "Url to go after a successful checkout"},
						Raw:  func() any { return &d.Req.SuccessUrl },
					},
					"cancel_url": {
						Desc: &schemaz.Desc{Summary: "Url to go after a canceled checkout"},
						Raw:  func() any { return &d.Req.CancelUrl },
					},
					"state": {
						Desc: &schemaz.Desc{Summary: "Base64 state token with the trial days and the metadata of the checkout; the service parses it"},
						Raw:  func() any { return &d.Req.State },
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
