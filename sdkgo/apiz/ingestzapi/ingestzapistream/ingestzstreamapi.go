package ingestzapistream

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketbnamez"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkobjectz"
)

type StreamReq[T sdkobjectz.SdkObject] struct {
	Dataset lakedsz.Dataset
	Table   laketbnamez.TableName
	Data    []T
}

type Api[T sdkobjectz.SdkObject] struct {
	Req StreamReq[T]
}

func (d *Api[T]) ApiSpec() apiz.ApiSpec {
	var x T
	return apiz.ApiSpec{
		Id: "ingestz-stream",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "Stream Ingestion",
				Summary: "Ingests stream data into a target dataset table.",
			}
		},
		Method: "POST",
		Path:   "/api/ingest/dataset/{dataset}/stream/" + x.Id(),
	}
}

func (d *Api[T]) ApiDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {
						Desc: &schemaz.Desc{Summary: "Name of the dataset that receives the stream data."},
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
					"table": {
						Desc: &schemaz.Desc{Summary: "Name of the table that receives the stream data."},
						Raw:  func() any { return &d.Req.Table },
					},
					"data": {
						Desc: &schemaz.Desc{Summary: "Generic stream data to ingest."},
						Array: func() *schemaz.Array {
							return &schemaz.Array{
								// implement it
							}
						},
					},
				}
			},
		},
		RespHeaders: nil,
		RespBody:    nil,
	}
}
