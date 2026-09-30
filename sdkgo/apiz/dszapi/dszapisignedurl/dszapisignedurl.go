package dszapisignedurl

import (
	"github.com/infinity6-ai/gox/commonz/slicez"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakepathz"
)

type Path struct {
	Method string         `json:"method"`
	Path   lakepathz.Path `json:"path"`
}

type Req struct {
	Dataset lakedsz.Dataset `json:"dataset"`
	Paths   []Path          `json:"paths"`
}

type Urls []string

type Resp struct {
	Urls Urls `json:"urls"`
}

func (*Resp) JSType() string {
	return jstypez.TypeObject
}

type Api struct {
	Req  Req  `json:"req"`
	Resp Resp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id: "dsz-signed-url",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "write it",
				Summary: "write it",
			}
		},
		Method: "POST",
		Path:   "/api/dsz/signed-url/dataset/{dataset}",
	}
}

func (d *Api) GetDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {
						Desc: &schemaz.Desc{Summary: "write it1"},
						Str:  lakedsz.ParseStr(&d.Req.Dataset),
					},
				}
			},
		},
		QueryParams: nil,
		ReqHeaders:  nil,
		ReqBody: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"paths": {
						Desc: &schemaz.Desc{Summary: "write it2"},
						Array: func() *schemaz.Array {
							return &schemaz.Array{
								Get: func(idx int, read bool) *schemaz.Schema {
									if !read {
										d.Req.Paths = slicez.GrowLenTo(d.Req.Paths, idx+1)
									}
									return &schemaz.Schema{
										Desc: &schemaz.Desc{Summary: "write it3"},
										Object: func(read bool) map[string]*schemaz.Schema {
											return map[string]*schemaz.Schema{
												"method": {
													Desc: &schemaz.Desc{Summary: "write it4"},
													Raw:  func() any { return &d.Req.Paths[idx].Method },
												},
												"path": {
													Desc: &schemaz.Desc{Summary: "write it5"},
													Raw:  func() any { return &d.Req.Paths[idx].Path },
												},
											}
										},
									}
								},
							}
						},
					},
				}
			},
		},
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Desc: &schemaz.Desc{Name: "Get URL Response", Summary: "write it6"},
			Raw:  func() any { return &d.Resp },
		},
	}
}
