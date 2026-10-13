package pipezlist

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solutiondomain"
)

type PipeState struct {
	Status    string         `json:"status"`
	Params    map[string]any `json:"params"`
	CreatedAt string         `json:"created_at"`
}

type Pipeline struct {
	Id           string       `json:"id"`
	ConcurrentId string       `json:"concurrent_id"`
	Angle        string       `json:"angle"`
	Name         string       `json:"name"`
	Dataset      string       `json:"dataset"`
	States       []*PipeState `json:"states"`
}

type PipezListReq struct {
	Dataset        lakedsz.Dataset               `json:"dataset"`
	SolutionDomain solutiondomain.SolutionDomain `json:"solution_domain"`
	Name           string                        `json:"name"`
	Status         string                        `json:"status"`
	Cursor         string                        `json:"cursor"`
	Limit          int                           `json:"limit"`
}

type PipezListResp struct {
	Pipelines []*Pipeline `json:"pipelines"`
	Cursor    string      `json:"cursor"`
}

func (*PipezListResp) JSType() string {
	return jstypez.TypObject
}

type Api struct {
	Req  *PipezListReq  `json:"req"`
	Resp *PipezListResp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "pipez-list",
		Desc:   nil,
		Method: "GET",
		Path:   "/api/pipe/list/dataset/{dataset}",
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
					d.Req = &PipezListReq{}
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
					d.Req = &PipezListReq{}
				}
				return map[string]*schemaz.Schema{
					"solution_domain": {Str: solutiondomain.ParseStr(&d.Req.SolutionDomain)},
					"name":            {Str: schemaz.ParseStr(&d.Req.Name)},
					"status":          {Str: schemaz.ParseStr(&d.Req.Status)},
					"cursor":          {Str: schemaz.ParseStr(&d.Req.Cursor)},
					"limit":           {Str: schemaz.ParseStrNumber(&d.Req.Limit)},
				}
			},
		},
		ReqHeaders:  nil,
		ReqBody:     nil,
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Raw: func() any { return &d.Resp },
		},
	}
}
