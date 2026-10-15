package pipezget

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/pipez"
)

type PipezGetReq struct {
	Dataset lakedsz.Dataset `json:"dataset"`
	Id      string          `json:"id"`
}

type PipezGetResp struct {
	Pipeline *pipez.Pipeline `json:"pipeline"`
}

func (*PipezGetResp) JSType() string {
	return jstypez.TypObject
}

type Api struct {
	Req  *PipezGetReq  `json:"req"`
	Resp *PipezGetResp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "pipez-get",
		Desc:   nil,
		Method: "GET",
		Path:   "/api/pipe/dataset/{dataset}/id/{id}",
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
					d.Req = &PipezGetReq{}
				}
				return map[string]*schemaz.Schema{
					"dataset": {Str: lakedsz.ParseStr(&d.Req.Dataset)},
					"id":      {Str: schemaz.ParseStr(&d.Req.Id)},
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
					d.Resp = &PipezGetResp{}
				}
				return map[string]*schemaz.Schema{
					"pipeline": {
						Object: func(read bool) map[string]*schemaz.Schema {
							if read && d.Resp.Pipeline == nil {
								return nil
							}
							if d.Resp.Pipeline == nil {
								d.Resp.Pipeline = &pipez.Pipeline{}
							}
							return map[string]*schemaz.Schema{
								"id": {
									Raw: func() any { return &d.Resp.Pipeline.Id },
								},
								"concurrent_id": {
									Raw: func() any { return &d.Resp.Pipeline.ConcurrentId },
								},
								"solution_domain": {
									Raw: func() any { return &d.Resp.Pipeline.SolutionDomain },
								},
								"name": {
									Raw: func() any { return &d.Resp.Pipeline.Name },
								},
								"dataset": {
									Raw: func() any { return &d.Resp.Pipeline.Dataset },
								},
								"states": {
									Array: schemaz.NewArray(schemaz.ArrayOptions[[]*pipez.PipeState, *pipez.PipeState]{
										Slice:      func() *[]*pipez.PipeState { return &d.Resp.Pipeline.States },
										NewElement: func(idx int) *pipez.PipeState { return &pipez.PipeState{} },
										Get: func(idx int, element **pipez.PipeState) *schemaz.Schema {
											return &schemaz.Schema{
												Object: func(_ bool) map[string]*schemaz.Schema {
													return map[string]*schemaz.Schema{
														"status": {
															Raw: func() any { return &(*element).Status },
														},
														"params": {
															Raw: func() any { return &(*element).Params },
														},
														"created_at": {
															Raw: func() any { return &(*element).CreatedAt },
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
			},
		},
	}
}
