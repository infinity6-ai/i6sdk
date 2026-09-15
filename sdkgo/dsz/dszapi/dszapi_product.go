package dszapi

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

type Product string

const RaasProduct Product = "raas"
const ForecastProduct Product = "forecast"

var AllProducts = []Product{
	RaasProduct,
	ForecastProduct,
}

type DszProductSelectReq struct {
	Dataset lakedsz.Dataset `json:"dataset"`
	Product string          `json:"product"`
}

type DszProductSelectResp struct {
	Dataset  lakedsz.Dataset `json:"dataset"`
	Products []string        `json:"products"`
}

type DszProductSelectApi struct {
	Req  *DszProductSelectReq  `json:"req"`
	Resp *DszProductSelectResp `json:"resp"`
}

func (d *DszProductSelectApi) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-product-select",
		Desc:   nil,
		Method: "POST",
		Path:   "/api/dsz/product-select/dataset/{dataset}",
	}
}

func (d *DszProductSelectApi) GetDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Req == nil {
					return nil
				}
				if d.Req == nil {
					d.Req = &DszProductSelectReq{}
				}
				return map[string]*schemaz.Schema{
					"dataset": {Str: lakedsz.ParseStr(&d.Req.Dataset)},
				}
			},
		},
		QueryParams: nil,
		ReqHeaders:  nil,
		ReqBody: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Req == nil {
					return nil
				}
				if d.Req == nil {
					d.Req = &DszProductSelectReq{}
				}
				return map[string]*schemaz.Schema{
					"product": {
						Raw: func() any { return &d.Req.Product },
					},
				}
			},
		},
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Resp == nil {
					return nil
				}
				if d.Resp == nil {
					d.Resp = &DszProductSelectResp{}
				}
				return map[string]*schemaz.Schema{
					"dataset": {
						Raw: func() any { return &d.Resp.Dataset },
					},
					"products": {
						Raw: func() any { return &d.Resp.Products },
					},
				}
			},
		},
	}
}
