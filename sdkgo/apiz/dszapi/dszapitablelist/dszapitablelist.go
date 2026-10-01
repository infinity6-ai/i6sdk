package dszapitablelist

import (
	"github.com/infinity6-ai/gox/fsz/fsz"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakepathz"
)

type Req struct {
	Dataset   lakedsz.Dataset `json:"dataset"`
	Path      lakepathz.Path  `json:"path"`      // query
	Sign      bool            `json:"sign"`      // query
	Recursive bool            `json:"recursive"` // query
	Stat      bool            `json:"stat"`      // query
	Cursor    string          `json:"cursor"`    // query
	Limit     int             `json:"limist"`    // query
}

type File struct {
	Path      lakepathz.Path `json:"path"`
	Stat      fsz.FileStat   `json:"stat"`
	SignedUrl string         `json:"signed_url"`
}

func (*File) JSType() string {
	return jstypez.TypeObject
}

type Resp struct {
	Files []File `json:"files"`
}

type Api struct {
	Req  Req  `json:"req"`
	Resp Resp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id: "dsz-table-list",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "List Dataset Tables",
				Summary: "Lists files and tables in the specified dataset.",
			}
		},
		Method: "GET",
		Path:   "/api/dsz/dataset/{dataset}/table-list",
	}
}

func (d *Api) GetDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {
						Desc: &schemaz.Desc{Summary: "Name of the dataset to list tables from."},
						Str:  lakedsz.ParseStr(&d.Req.Dataset),
					},
				}
			},
		},
		QueryParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"path": {
						Desc: &schemaz.Desc{Summary: "Path to list."},
						Raw:  func() any { return &d.Req.Path },
					},
					"sign": {
						Desc: &schemaz.Desc{Summary: "Whether to return signed URLs."},
						Raw:  func() any { return &d.Req.Sign },
					},
					"recursive": {
						Desc: &schemaz.Desc{Summary: "Whether to list recursively."},
						Raw:  func() any { return &d.Req.Recursive },
					},
					"stat": {
						Desc: &schemaz.Desc{Summary: "Whether to include file stats."},
						Raw:  func() any { return &d.Req.Stat },
					},
					"cursor": {
						Desc: &schemaz.Desc{Summary: "Pagination cursor."},
						Str:  schemaz.ParseStr(&d.Req.Cursor),
					},
					"limist": {
						Desc: &schemaz.Desc{Summary: "Limit the number of results."},
						Str:  schemaz.ParseStrNumber(&d.Req.Limit),
					},
				}
			},
		},
		ReqHeaders:  nil,
		ReqBody:     nil,
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Desc: &schemaz.Desc{Summary: "List of files in the dataset."},
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"files": {
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]File, File]{
							Slice: func() *[]File { return &d.Resp.Files },
							Get: func(idx int, element *File) *schemaz.Schema {
								return &schemaz.Schema{
									Desc: &schemaz.Desc{Summary: "File entry"},
									Raw: func() any {
										return element
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
