package dszapitableget

import (
	"strings"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakepathz"
)

type Req struct {
	Dataset lakedsz.Dataset `json:"dataset"`
	Path    lakepathz.Path  `json:"path"`
}

type Resp struct {
	Location string `json:"location"`
}

type Api struct {
	Req  Req  `json:"req"`
	Resp Resp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id: "dsz-table-get",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "Get Dataset Table File",
				Summary: "Get table file by redirecting to its signed URL.",
			}
		},
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/table-file/{path...}",
	}
}

func (d *Api) GetDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {
						Desc: &schemaz.Desc{Summary: "Name of the dataset."},
						Str:  lakedsz.ParseStr(&d.Req.Dataset),
					},
					"path": {
						Desc: &schemaz.Desc{Summary: "Path to the file without schema."},
						Str: func() *schemaz.Parser[string] {
							return &schemaz.Parser[string]{
								Parse: func(v string) {
									d.Req.Path = lakepathz.MustParse("i6tb://" + v)
								},
								Format: func() (string, any) {
									var v string
									if d.Req.Path.Optional().IsPresent() {
										v = strings.TrimPrefix(d.Req.Path.Get(), "i6tb://")
									}
									return v, v
								},
							}
						},
					},
				}
			},
		},
		QueryParams: nil,
		ReqHeaders:  nil,
		ReqBody:     nil,
		RespHeaders: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"location": {
						Desc: &schemaz.Desc{Summary: "Redirect URL for the file."},
						Str:  schemaz.ParseStr(&d.Resp.Location),
					},
				}
			},
		},
		RespBody: nil,
	}
}
