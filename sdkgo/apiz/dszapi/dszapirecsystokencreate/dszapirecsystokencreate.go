package dszapirecsystokencreate

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/dsztoken"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

type DszCreateRecsysTokenReq struct {
	Dataset    lakedsz.Dataset `json:"dataset"`
	Annotation string          `json:"annotation"`
	Channel    string          `json:"channel"`
}

type Api struct {
	Req  DszCreateRecsysTokenReq `json:"req"`
	Resp dsztoken.DatasetToken   `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-create-recsys-token",
		Method: "POST",
		Path:   "/api/ds/dataset/{dataset}/recsystoken/create",
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
		ReqBody: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"annotation": {Field: func() any { return &d.Req.Annotation }},
					"channel":    {Field: func() any { return &d.Req.Channel }},
				}
			},
		},
		RespHeaders: nil,
		RespBody:    d.Resp.Schema(),
	}
}
