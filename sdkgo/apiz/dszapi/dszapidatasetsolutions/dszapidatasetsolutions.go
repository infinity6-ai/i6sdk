package dszapidatasetsolutions

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solution"
)

type DszListActiveSolutionsReq struct {
	Dataset lakedsz.Dataset
}

type DszListActiveSolutionsResp struct {
	Dataset   lakedsz.Dataset
	Solutions []solution.Solution
}

type Api struct {
	Req  *DszListActiveSolutionsReq
	Resp *DszListActiveSolutionsResp
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-dataset-solutions",
		Desc:   nil,
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/solutions",
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
					d.Req = &DszListActiveSolutionsReq{}
				}
				return map[string]*schemaz.Schema{
					"dataset": {UrlValue: schemaz.ParseStrGet(&d.Req.Dataset, lakedsz.MustParse)},
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
					d.Resp = &DszListActiveSolutionsResp{}
				}
				return map[string]*schemaz.Schema{
					"dataset": {Field: func() any { return &d.Resp.Dataset }},
					"solutions": {
						Field: func() any { return &d.Resp.Solutions },
					},
				}
			},
		},
	}
}
