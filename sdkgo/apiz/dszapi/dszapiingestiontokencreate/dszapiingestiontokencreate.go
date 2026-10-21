package dszapiingestiontokencreate

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/dsztoken"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

type DszCreateIngestionTokenReq struct {
	Dataset    lakedsz.Dataset `json:"dataset"`
	Annotation string          `json:"annotation"`
}

type Api struct {
	Req  DszCreateIngestionTokenReq `json:"req"`
	Resp dsztoken.DatasetToken      `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-create-ingestion-token",
		Method: "POST",
		Path:   "/api/ds/dataset/{dataset}/ingestiontoken/create",
	}
}

func (d *Api) ApiDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {Str: schemaz.ParseStrGet(&d.Req.Dataset, lakedsz.MustParse)},
				}
			},
		},
		QueryParams: nil,
		ReqHeaders:  nil,
		ReqBody: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"annotation": {Str: schemaz.ParseStr(&d.Req.Annotation)},
				}
			},
		},
		RespHeaders: nil,
		RespBody:    d.Resp.Schema(),
	}
}
