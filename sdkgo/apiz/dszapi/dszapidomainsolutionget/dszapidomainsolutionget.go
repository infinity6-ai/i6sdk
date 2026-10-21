package dszapidomainsolutionget

import (
	"reflect"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/domain"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solution"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/storeid"
)

type DszSolutionDomainGetReq struct {
	Dataset lakedsz.Dataset
	StoreId storeid.StoreId
}

type SolutionDomain struct {
	Solution solution.Solution `json:"solution"`
	Domain   domain.Domain     `json:"domain"`
}

func (s *SolutionDomain) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"solution": {
					Field: func() any { return &s.Solution },
				},
				"domain": {
					Field: func() any { return &s.Domain },
				},
			}
		},
	}
}

func (s *SolutionDomain) JSTypedef() (*jstypez.Typedef, *jstypez.Typedefs) {
	t := reflect.TypeOf(*s).String()
	td := &jstypez.Typedef{
		Head:       &jstypez.Tuple{Typ: jstypez.TypObject, GoTyp: t, Name: ""},
		Properties: &jstypez.Tuples{},
	}
	td.Properties.Add(
		jstypez.Tuple{Typ: jstypez.TypString, GoTyp: reflect.TypeOf(s.Solution).String(), Name: "solution"},
		jstypez.Tuple{Typ: jstypez.TypString, GoTyp: reflect.TypeOf(s.Domain).String(), Name: "domain"},
	)
	return td, nil
}

type DszSolutionDomainGetResp struct {
	Dataset   lakedsz.Dataset  `json:"dataset"`
	Solutions []SolutionDomain `json:"solutions"`
	StoreId   storeid.StoreId  `json:"store_id"`
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

func (d *Api) ApiDataRefs() *apiz.DataRefs {
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
					"dataset": {UrlValue: schemaz.ParseStrGet(&d.Req.Dataset, lakedsz.MustParse)},
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
					"store_id": {UrlValue: schemaz.ParseStrGet(&d.Req.StoreId, storeid.MustParse)},
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
					"dataset": {Field: func() any { return &d.Resp.Dataset }},
					"solutions": {
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]SolutionDomain, SolutionDomain]{
							Slice: func() *[]SolutionDomain { return &d.Resp.Solutions },
							Get: func(idx int, element *SolutionDomain) *schemaz.Schema {
								return element.Schema()
							},
						}),
					},
					"store_id": {
						Field: func() any { return &d.Resp.StoreId },
					},
				}
			},
		},
	}
}
