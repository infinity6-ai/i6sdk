package pipezstart

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solutiondomain"
)

type PipeParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (*PipeParam) JSType() string {
	return jstypez.TypObject
}

type PipezStartReq struct {
	Dataset        lakedsz.Dataset               `json:"dataset"`
	SolutionDomain solutiondomain.SolutionDomain `json:"solution_domain"`
	Pipename       string                        `json:"pipename"`
	Params         []*PipeParam                  `json:"params"`
}

type PipezStartResp struct {
	Id string `json:"id"`
}

func (*PipezStartResp) JSType() string {
	return jstypez.TypObject
}

type Api struct {
	Req  *PipezStartReq  `json:"req"`
	Resp *PipezStartResp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "pipez-start",
		Desc:   nil,
		Method: "POST",
		Path:   "/api/pipe/start/dataset/{dataset}",
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
					d.Req = &PipezStartReq{}
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
				return map[string]*schemaz.Schema{
					"solution_domain": {Raw: func() any { return &d.Req.SolutionDomain }},
					"pipename":        {Raw: func() any { return &d.Req.Pipename }},
					"params":          {Raw: func() any { return &d.Req.Params }},
				}
			},
		},
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Raw: func() any { return &d.Resp },
		},
	}
}
