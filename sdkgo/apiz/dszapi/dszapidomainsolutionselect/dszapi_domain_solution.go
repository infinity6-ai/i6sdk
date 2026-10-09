package dszapidomainsolutionselect

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/domain"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solution"
)

type DszSolutionDomainSelectReq struct {
	Dataset  lakedsz.Dataset   `json:"dataset"`
	Solution solution.Solution `json:"solution"`
	Domain   domain.Domain     `json:"domain"`
	StoreId  string            `json:"store_id"`
}

type Api struct {
	Req *DszSolutionDomainSelectReq `json:"req"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-domain-select",
		Desc:   nil,
		Method: "POST",
		Path:   "/api/ds/dataset/{dataset}/store/solution-domain",
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
					d.Req = &DszSolutionDomainSelectReq{}
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
					d.Req = &DszSolutionDomainSelectReq{}
				}
				return map[string]*schemaz.Schema{
					"solution": {
						Raw: func() any { return &d.Req.Solution },
					},
					"domain": {
						Raw: func() any { return &d.Req.Domain },
					},
					"store_id": {
						Raw: func() any { return &d.Req.StoreId },
					},
				}
			},
		},
		RespHeaders: nil,
		RespBody:    nil,
	}
}
