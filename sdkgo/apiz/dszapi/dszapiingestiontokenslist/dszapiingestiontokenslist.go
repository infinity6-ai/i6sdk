package dszapiingestiontokenslist

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/dsztoken"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

type DszListIngestionTokensReq struct {
	Dataset lakedsz.Dataset `json:"dataset"`
}

type DszListIngestionTokensResp []dsztoken.DatasetToken

func (*DszListIngestionTokensResp) JSType() string {
	return jstypez.TypObject
}

type Api struct {
	Req  DszListIngestionTokensReq  `json:"req"`
	Resp DszListIngestionTokensResp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-list-ingestion-tokens",
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/ingestiontoken/list",
	}
}

func (d *Api) ApiDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
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
			Field: func() any { return &d.Resp },
		},
	}
}
