package dszapi

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

type RaasSolution string

const FashionRaasSolution RaasSolution = "fashion"

var AllRaasSolutions = []RaasSolution{
	FashionRaasSolution,
}

type DszRaasSolutionSelectReq struct {
	Dataset  lakedsz.Dataset `json:"dataset"`
	Solution string          `json:"solution"`
	Store    string          `json:"store"`
}

type DszRaasSolutionSelectResp struct {
	Dataset  lakedsz.Dataset `json:"dataset"`
	Solution string          `json:"solution"`
	Store    string          `json:"store"`
}

type DszRaasSolutionSelectApi struct {
	Req  *DszRaasSolutionSelectReq  `json:"req"`
	Resp *DszRaasSolutionSelectResp `json:"resp"`
}

func (d *DszRaasSolutionSelectApi) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-raas-solution-select",
		Desc:   nil,
		Method: "POST",
		Path:   "/api/dsz/raas-solution-select/dataset/{dataset}",
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
					d.Req = &DszRaasSolutionSelectReq{}
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
					d.Req = &DszRaasSolutionSelectReq{}
				}
				return map[string]*schemaz.Schema{
					"solution": {
						Raw: func() any { return &d.Req.Solution },
					},
					"store": {
						Raw: func() any { return &d.Req.Store },
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
					d.Resp = &DszRaasSolutionSelectResp{}
				}
				return map[string]*schemaz.Schema{
					"dataset": {
						Raw: func() any { return &d.Resp.Dataset },
					},
					"solution": {
						Raw: func() any { return &d.Resp.Solution },
					},
					"store": {
						Raw: func() any { return &d.Resp.Store },
					},
				}
			},
		},
	}
}
