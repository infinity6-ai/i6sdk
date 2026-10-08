package ingestzapistream

import (
	"encoding/json"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketbnamez"
)

type Payload struct {
	Value any
}

func (*Payload) JSType() string {
	return jstypez.TypAny
}

func (p *Payload) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Value)
}

func (p *Payload) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p.Value)
}

type StreamReq struct {
	Dataset  lakedsz.Dataset       `json:"dataset"`
	Table    laketbnamez.TableName `json:"table"`
	Payloads []Payload             `json:"payloads"`
}

type Api struct {
	Req StreamReq `json:"req"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id: "ingestz-stream",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "Stream Ingestion",
				Summary: "Ingests stream payloads into a target dataset table.",
			}
		},
		Method: "POST",
		Path:   "/api/ingest/stream/dataset/{dataset}",
	}
}

func (d *Api) ApiDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {
						Desc: &schemaz.Desc{Summary: "Name of the dataset that receives the stream payloads."},
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
						Desc: &schemaz.Desc{Summary: "Name of the table that receives the stream payloads."},
						Raw:  func() any { return &d.Req.Table },
					},
					"payloads": {
						Desc: &schemaz.Desc{Summary: "Generic stream payloads to ingest."},
						Raw:  func() any { return &d.Req.Payloads },
					},
				}
			},
		},
		RespHeaders: nil,
		RespBody:    nil,
	}
}
