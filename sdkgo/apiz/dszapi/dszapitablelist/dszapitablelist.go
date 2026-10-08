package dszapitablelist

import (
	"github.com/infinity6-ai/gox/fsz/fsz"
	"github.com/infinity6-ai/gox/fsz/fsz/fszsignmethod"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketbpathz"
)

type Req struct {
	Dataset    lakedsz.Dataset          `json:"dataset"`
	Path       laketbpathz.Path         `json:"path"`
	SignMethod fszsignmethod.SignMethod `json:"sign_method"`
	Recursive  bool                     `json:"recursive"`
	Stat       bool                     `json:"stat"`
	FromCursor string                   `json:"from_cursor"`
	Limit      int                      `json:"limit"`
}

type File struct {
	Path      laketbpathz.Path `json:"path"`
	Stat      *fsz.FileStat    `json:"stat"`
	SignedUrl string           `json:"signed_url"`
}

type Resp struct {
	Files      []*File `json:"files"`
	NextCursor string  `json:"next_cursor"`
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
		Path:   "/api/ds/dataset/{dataset}/table-list",
	}
}

func (d *Api) ApiDataRefs() *apiz.DataRefs {
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
						Str:  schemaz.ParseStrGet(&d.Req.Path, laketbpathz.MustParse),
					},
					"sign_method": {
						Desc: &schemaz.Desc{Summary: "Whether to return signed URLs."},
						Str:  schemaz.ParseStrGet(&d.Req.SignMethod, fszsignmethod.MustParseOptional),
					},
					"recursive": {
						Desc: &schemaz.Desc{Summary: "Whether to list recursively."},
						Str:  schemaz.ParseStrBool(&d.Req.Recursive),
					},
					"stat": {
						Desc: &schemaz.Desc{Summary: "Whether to include file stats."},
						Str:  schemaz.ParseStrBool(&d.Req.Stat),
					},
					"from_cursor": {
						Desc: &schemaz.Desc{Summary: "Pagination cursor."},
						Str:  schemaz.ParseStr(&d.Req.FromCursor),
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
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]*File, *File]{
							Slice:      func() *[]*File { return &d.Resp.Files },
							NewElement: func(idx int) *File { return &File{} },
							Get: func(idx int, element **File) *schemaz.Schema {
								return &schemaz.Schema{
									Desc: &schemaz.Desc{Summary: "File entry"},
									Object: func(_ bool) map[string]*schemaz.Schema {
										return map[string]*schemaz.Schema{
											"path": {
												Desc: &schemaz.Desc{Summary: "write it"},
												Raw:  func() any { return &(*element).Path },
											},
											"stat": {
												Desc: &schemaz.Desc{Summary: "write it"},
												Raw:  func() any { return &(*element).Stat },
											},
											"signed_url": {
												Desc: &schemaz.Desc{Summary: "write it"},
												Raw:  func() any { return &(*element).SignedUrl },
											},
										}
									},
								}
							},
						}),
					},
					"next_cursor": {
						Desc: &schemaz.Desc{Summary: "Cursor for the next page of results."},
						Str:  schemaz.ParseStr(&d.Resp.NextCursor),
					},
				}
			},
		},
	}
}
