package pipezlist

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solutiondomain"
	"go.code.infinity6.ai/sdkgo/sdkgo/pipez"
)

type PipezListReq struct {
	Dataset        lakedsz.Dataset
	SolutionDomain solutiondomain.SolutionDomain
	Name           string
	Status         string
	Cursor         string
	Limit          int
}

type PipezListResp struct {
	Pipelines []*pipez.Pipeline
	Cursor    string
}

func (*PipezListResp) JSType() string {
	return jstypez.TypObject
}

type Api struct {
	Req  *PipezListReq
	Resp *PipezListResp
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "pipez-list",
		Desc:   nil,
		Method: "GET",
		Path:   "/api/pipe/dataset/{dataset}/list",
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
					"dataset": {Str: schemaz.ParseStrGet(&d.Req.Dataset, lakedsz.MustParse)},
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
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Resp == nil {
					return nil
				}
				if d.Resp == nil {
					d.Resp = &PipezListResp{}
				}
				return map[string]*schemaz.Schema{
					"cursor": {
						Raw: func() any { return &d.Resp.Cursor },
					},
					"pipelines": {
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]*pipez.Pipeline, *pipez.Pipeline]{
							Slice:      func() *[]*pipez.Pipeline { return &d.Resp.Pipelines },
							NewElement: func(idx int) *pipez.Pipeline { return &pipez.Pipeline{} },
							Get: func(idx int, element **pipez.Pipeline) *schemaz.Schema {
								return &schemaz.Schema{
									Object: func(read bool) map[string]*schemaz.Schema {
										if read && *element == nil {
											return nil
										}
										if *element == nil {
											*element = &pipez.Pipeline{}
										}
										return (*element).Schema().Object(read)
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
