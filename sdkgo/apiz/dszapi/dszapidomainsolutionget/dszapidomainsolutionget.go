package dszapidomainsolutionget

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/domain"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solution"
)

type DszSolutionDomainGetReq struct {
	Dataset lakedsz.Dataset `json:"dataset"`
	Store   string          `json:"store"`
}

type SolutionDomain struct {
	Solution solution.Solution `json:"solution"`
	Domain   domain.Domain     `json:"domain"`
}

func (*SolutionDomain) JSTypedef() (*jstypez.Typedef, *jstypez.Typedefs) {
	td := &jstypez.Typedef{
		Head:       &jstypez.Tuple{Typ: jstypez.TypObject, Name: ""},
		Properties: &jstypez.Tuples{},
	}
	td.Properties.Add(
		jstypez.Tuple{Typ: jstypez.TypString, Name: "solution"},
		jstypez.Tuple{Typ: jstypez.TypString, Name: "domain"},
	)
	return td, nil
}

type DszSolutionDomainGetResp struct {
	Dataset   lakedsz.Dataset  `json:"dataset"`
	Solutions []SolutionDomain `json:"solutions"`
	Store     string           `json:"store"`
}

type Api struct {
	Req  *DszSolutionDomainGetReq  `json:"req"`
	Resp *DszSolutionDomainGetResp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-domain-get",
		Desc:   nil,
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/store/solution-domain",
	}
}

func (d *Api) GetDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Req == nil {
					return nil
				}
				if d.Req == nil {
					d.Req = &DszSolutionDomainGetReq{}
				}
				return map[string]*schemaz.Schema{
					"dataset": {Str: lakedsz.ParseStr(&d.Req.Dataset)},
				}
			},
		},
		QueryParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Req == nil {
					return nil
				}
				if d.Req == nil {
					d.Req = &DszSolutionDomainGetReq{}
				}
				return map[string]*schemaz.Schema{
					"store": {Str: schemaz.ParseStr(&d.Req.Store)},
				}
			},
		},
		ReqHeaders:  nil,
		ReqBody:     nil,
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Resp == nil {
					return nil
				}
				if d.Resp == nil {
					d.Resp = &DszSolutionDomainGetResp{}
				}
				return map[string]*schemaz.Schema{
					"dataset": {Str: lakedsz.ParseStr(&d.Resp.Dataset)},
					"solutions": {
						Raw: func() any { return &d.Resp.Solutions },
					},
					"store": {
						Raw: func() any { return &d.Resp.Store },
					},
				}
			},
		},
	}
}
