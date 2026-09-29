package ingestzgeturlapi

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/ingestztoken"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketbnamez"
)

type GetUrlReq struct {
	Dataset    lakedsz.Dataset       `json:"dataset"`
	Table      laketbnamez.TableName `json:"table"`
	Partitions map[string]string     `json:"partitions"`
	Amount     int                   `json:"amount"`
}

type GetUrlResp struct {
	Uploads    []string                 `json:"uploads"`
	IngestedAt ingestztoken.IngestToken `json:"ingest_id"`
}

func (*GetUrlResp) JSType() string {
	return jstypez.TypeObject
}

type GetUrlApi struct {
	Req  GetUrlReq  `json:"req"`
	Resp GetUrlResp `json:"resp"`
}

func (d *GetUrlApi) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id: "ingestz-get-url",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "Ingestz get url",
				Summary: "Get the url that will be used to ingest files",
			}
		},
		Method: "POST",
		Path:   "/api/ingest/get-url/dataset/{dataset}",
	}
}

func (d *GetUrlApi) GetDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {Str: lakedsz.ParseStr(&d.Req.Dataset)},
				}
			},
		},
		QueryParams: nil,
		ReqHeaders:  nil,
		ReqBody: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"table":      {Raw: func() any { return &d.Req.Table }},
					"partitions": {Raw: func() any { return &d.Req.Partitions }},
					"amount":     {Raw: func() any { return &d.Req.Amount }},
				}
			},
		},
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Raw: func() any { return &d.Resp },
		},
	}
}
