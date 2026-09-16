package dszapi

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

type Solution string

const RaasProduct Solution = "raas"
const ForecastSolution Solution = "forecast"

var AllSolutions = []Solution{
	RaasProduct,
	ForecastSolution,
}

type DszSolutionSelectReq struct {
	Dataset  lakedsz.Dataset `json:"dataset"`
	Solution string          `json:"solution"`
}

type DszSolutionSelectResp struct {
	Dataset   lakedsz.Dataset `json:"dataset"`
	Solutions []string        `json:"solutions"`
}

type DszSolutionSelectApi struct {
	Req  *DszSolutionSelectReq  `json:"req"`
	Resp *DszSolutionSelectResp `json:"resp"`
}

func (d *DszSolutionSelectApi) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-solution-select",
		Desc:   nil,
		Method: "POST",
		Path:   "/api/dsz/solution-select/dataset/{dataset}",
	}
}

func (d *DszSolutionSelectApi) GetDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Req == nil {
					return nil
				}
				if d.Req == nil {
					d.Req = &DszSolutionSelectReq{}
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
					d.Req = &DszSolutionSelectReq{}
				}
				return map[string]*schemaz.Schema{
					"solution": {
						Raw: func() any { return &d.Req.Solution },
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
					d.Resp = &DszSolutionSelectResp{}
				}
				return map[string]*schemaz.Schema{
					"dataset": {
						Raw: func() any { return &d.Resp.Dataset },
					},
					"solutions": {
						Raw: func() any { return &d.Resp.Solutions },
					},
				}
			},
		},
	}
}
