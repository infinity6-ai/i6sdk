package dszapipendingpaymentlist

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

type Req struct {
	User string
}

type PendingPayment struct {
	Dataset lakedsz.Dataset
}

func (p *PendingPayment) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"dataset": {
					Desc: &schemaz.Desc{Summary: "Id of the dataset that is going to be created after the payment"},
					Raw:  func() any { return &p.Dataset },
				},
			}
		},
	}
}

type Resp struct {
	PendingPayments []*PendingPayment
}

type Api struct {
	Req  Req
	Resp Resp
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id: "dsz-list-pending-payments",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "List Pending Payments",
				Summary: "Lists the datasets of a user that are waiting for a payment.",
			}
		},
		Method: "GET",
		Path:   "/api/dsz/pending-payments",
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
				}
			},
		},
		ReqHeaders:  nil,
		ReqBody:     nil,
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Desc: &schemaz.Desc{Summary: "The pending payments of the user."},
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"pending_payments": {
						Desc: &schemaz.Desc{Summary: "Datasets waiting for a payment"},
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]*PendingPayment, *PendingPayment]{
							NewElement: func(idx int) *PendingPayment { return &PendingPayment{} },
							Slice:      func() *[]*PendingPayment { return &d.Resp.PendingPayments },
							Get:        func(idx int, element **PendingPayment) *schemaz.Schema { return (*element).Schema() },
						}),
					},
				}
			},
		},
	}
}
