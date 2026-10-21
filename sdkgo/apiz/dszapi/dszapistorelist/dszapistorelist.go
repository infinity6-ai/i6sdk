package dszapistorelist

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/dszstore"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

type DszListStoresReq struct {
	Dataset lakedsz.Dataset `json:"dataset"`
}

type DszListStoresResp struct {
	Dataset lakedsz.Dataset     `json:"dataset"`
	Stores  []dszstore.DszStore `json:"stores"`
}

type Api struct {
	Req  *DszListStoresReq  `json:"req"`
	Resp *DszListStoresResp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-list-stores",
		Desc:   nil,
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/store/list",
	}
}

func (d *Api) ApiDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Req == nil {
					return nil
				}
				if d.Req == nil {
					d.Req = &DszListStoresReq{}
				}
				return map[string]*schemaz.Schema{
					"dataset": {UrlValue: schemaz.ParseStrGet(&d.Req.Dataset, lakedsz.MustParse)},
				}
			},
		},
		QueryParams: nil,
		ReqHeaders:  nil,
		ReqBody:     nil,
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Resp == nil {
					return nil
				}
				if d.Resp == nil {
					d.Resp = &DszListStoresResp{}
				}
				if d.Resp.Stores == nil {
					d.Resp.Stores = []dszstore.DszStore{}
				}
				return map[string]*schemaz.Schema{
					"dataset": {Field: func() any { return &d.Resp.Dataset }},
					"stores": {
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]dszstore.DszStore, dszstore.DszStore]{
							Slice: func() *[]dszstore.DszStore { return &d.Resp.Stores },
							Get: func(idx int, element *dszstore.DszStore) *schemaz.Schema {
								return &schemaz.Schema{
									Object: func(_ bool) map[string]*schemaz.Schema {
										return map[string]*schemaz.Schema{
											"id": {
												Field: func() any { return &element.Id },
											},
											"name": {
												Field: func() any { return &element.Name },
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
