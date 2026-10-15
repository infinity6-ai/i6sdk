package ingestzapigeturl

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/ingestztoken"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketbnamez"
)

type GetUrlReq struct {
	Dataset    lakedsz.Dataset
	Table      laketbnamez.TableName
	Partitions map[string]string
	Amount     int
}

type GetUrlResp struct {
	Uploads  []string
	IngestId ingestztoken.IngestToken
}

func (*GetUrlResp) JSType() string {
	return jstypez.TypObject
}

type Api struct {
	Req  GetUrlReq
	Resp GetUrlResp
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id: "ingestz-get-url",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "Get Ingestion Upload URLs",
				Summary: "Generates presigned PUT URLs and an ingestion token for uploading files into a target dataset table and partitions.",
			}
		},
		Method: "POST",
		Path:   "/api/ingest/dataset/{dataset}/get-url",
	}
}

func (d *Api) ApiDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {
						Desc: &schemaz.Desc{Summary: "Name of the dataset that receives the uploaded files."},
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
						Desc: &schemaz.Desc{Summary: "Name of the table, inside the dataset, that receives the uploaded files."},
						Raw:  func() any { return &d.Req.Table },
					},
					"partitions": {
						Desc: &schemaz.Desc{Summary: "Partition key/value pairs that select the table partition the files are ingested into."},
						Raw:  func() any { return &d.Req.Partitions },
					},
					"amount": {
						Desc: &schemaz.Desc{Summary: "Number of presigned upload URLs to generate."},
						Raw:  func() any { return &d.Req.Amount },
					},
				}
			},
		},
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Desc: &schemaz.Desc{Name: "Get URL Response", Summary: "Presigned upload URLs and the ingestion token that identifies this ingestion."},
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"uploads": {
						Desc: &schemaz.Desc{Summary: "Presigned PUT URLs, one per requested file; upload each file to its URL."},
						Raw:  func() any { return &d.Resp.Uploads },
					},
					"ingest_id": {
						Desc: &schemaz.Desc{Name: "Ingest ID", Summary: "Ingestion token that identifies this ingestion and ties the uploaded files together."},
						Raw:  func() any { return &d.Resp.IngestId },
					},
				}
			},
		},
	}
}
