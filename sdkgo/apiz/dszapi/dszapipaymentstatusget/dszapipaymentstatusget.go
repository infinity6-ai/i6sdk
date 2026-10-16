package dszapipaymentstatusget

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solution"
)

type Req struct {
	User    string
	Dataset lakedsz.Dataset
}

type Resp struct {
	UserId        string
	Dataset       string
	DatasetName   string
	PaymentId     string
	PaymentType   string
	PaymentStatus string
	Solution      solution.Solution
}

type Api struct {
	Req  Req
	Resp Resp
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id: "dsz-get-payment-status",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "Get Payment Status",
				Summary: "Checks the payment of a dataset. When it is paid the dataset is created and granted to the user.",
			}
		},
		Method: "GET",
		Path:   "/api/ds/payment-status",
	}
}

func (d *Api) ApiDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: nil,
		QueryParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"user": {
						Desc: &schemaz.Desc{Summary: "User that is buying"},
						Str:  schemaz.ParseStr(&d.Req.User),
					},
					"dataset": {
						Desc: &schemaz.Desc{Summary: "Id of the dataset that is waiting for the payment"},
						Str:  lakedsz.ParseStr(&d.Req.Dataset),
					},
				}
			},
		},
		ReqHeaders:  nil,
		ReqBody:     nil,
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Desc: &schemaz.Desc{Summary: "The payment of the dataset."},
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"user_id": {
						Desc: &schemaz.Desc{Summary: "User that is buying"},
						Raw:  func() any { return &d.Resp.UserId },
					},
					"dataset": {
						Desc: &schemaz.Desc{Summary: "Id of the dataset"},
						Raw:  func() any { return &d.Resp.Dataset },
					},
					"dataset_name": {
						Desc: &schemaz.Desc{Summary: "Name of the dataset"},
						Raw:  func() any { return &d.Resp.DatasetName },
					},
					"payment_id": {
						Desc: &schemaz.Desc{Summary: "Id of the payment on the payment provider"},
						Raw:  func() any { return &d.Resp.PaymentId },
					},
					"payment_type": {
						Desc: &schemaz.Desc{Summary: "Payment provider"},
						Raw:  func() any { return &d.Resp.PaymentType },
					},
					"payment_status": {
						Desc: &schemaz.Desc{Summary: "paid, expired, open or unpaid. Only open and unpaid payments are still pending"},
						Raw:  func() any { return &d.Resp.PaymentStatus },
					},
					"solution": {
						Desc: &schemaz.Desc{Summary: "puchased solution"},
						Raw:  func() any { return &d.Resp.Solution },
					},
				}
			},
		},
	}
}
