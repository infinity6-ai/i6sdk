package dszapi

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/domain"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solution"
)

type DszDomainSelectReq struct {
	Dataset  lakedsz.Dataset   `json:"dataset"`
	Solution solution.Solution `json:"solution"`
	Domain   domain.Domain     `json:"domain"`
	Store    string            `json:"store"`
}

type DszRaasSolutionSelectApi struct {
	Req  *DszDomainSelectReq `json:"req"`
	Resp *DszDomainSelectReq `json:"resp"`
}

func (d *DszRaasSolutionSelectApi) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-domain-select",
		Desc:   nil,
		Method: "POST",
		Path:   "/api/dsz/domain-select/dataset/{dataset}",
	}
}

func (d *DszRaasSolutionSelectApi) GetDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Req == nil {
					return nil
				}
				if d.Req == nil {
					d.Req = &DszDomainSelectReq{}
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
					d.Req = &DszDomainSelectReq{}
				}
				return map[string]*schemaz.Schema{
					"solution": {
						Raw: func() any { return &d.Req.Solution },
					},
					"domain": {
						Raw: func() any { return &d.Req.Domain },
					},
					"store": {
						Raw: func() any { return &d.Req.Store },
					},
				}
			},
		},
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Raw: func() any { return &d.Resp },
		},
	}
}
