package dszapistripecheckout

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

type Req struct {
	DatasetName   string
	Plan          string
	SuccessUrl    string
	CancelUrl     string
	CampaignState string
	PaymentPhone  string
	PaymentEmail  string
	Document      string
}

type Resp struct {
	Id      string
	Url     string
	Dataset lakedsz.Dataset
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
					"dataset_name": {
						Desc: &schemaz.Desc{Summary: "Dataset name"},
						Field:  func() any { return &d.Req.DatasetName },
					},
					"plan": {
						Desc: &schemaz.Desc{Summary: "Lookup key of the Stripe price of the plan to buy; the service gets the price id from the Stripe prices api"},
						Field:  func() any { return &d.Req.Plan },
					},
					"success_url": {
						Desc: &schemaz.Desc{Summary: "Url to go after a successful checkout"},
						Field:  func() any { return &d.Req.SuccessUrl },
					},
					"cancel_url": {
						Desc: &schemaz.Desc{Summary: "Url to go after a canceled checkout"},
						Field:  func() any { return &d.Req.CancelUrl },
					},
					"campaign_state": {
						Desc: &schemaz.Desc{Summary: "Base64 state token with the trial days and the metadata of the checkout; the service parses it"},
						Field:  func() any { return &d.Req.CampaignState },
					},
					"payment_phone": {
						Desc: &schemaz.Desc{Summary: "Phone of the buyer"},
						Field:  func() any { return &d.Req.PaymentPhone },
					},
					"payment_email": {
						Desc: &schemaz.Desc{Summary: "Email of the buyer used in Stripe"},
						Field:  func() any { return &d.Req.PaymentEmail },
					},
					"document": {
						Desc: &schemaz.Desc{Summary: "Document of the buyer"},
						Field:  func() any { return &d.Req.Document },
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
						Field:  func() any { return &d.Resp.Id },
					},
					"url": {
						Desc: &schemaz.Desc{Summary: "Checkout session url"},
						Field:  func() any { return &d.Resp.Url },
					},
					"dataset": {
						Desc: &schemaz.Desc{Summary: "Id of the dataset that is going to be created after the payment (starts with ds)"},
						Field:  func() any { return &d.Resp.Dataset },
					},
				}
			},
		},
	}
}
