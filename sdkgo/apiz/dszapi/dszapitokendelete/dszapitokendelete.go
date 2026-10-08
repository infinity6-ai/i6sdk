package dszapitokendelete

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

type DszDeleteTokenReq struct {
	Dataset lakedsz.Dataset `json:"dataset"`
	Token   string          `json:"token"`
}

type Api struct {
	Req DszDeleteTokenReq `json:"req"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-delete-token",
		Method: "DELETE",
		Path:   "/api/ds/dataset/{dataset}/token/delete",
	}
}

func (d *Api) ApiDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {Str: lakedsz.ParseStr(&d.Req.Dataset)},
				}
			},
		},
		QueryParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"token": {Str: schemaz.ParseStr(&d.Req.Token)},
				}
			},
		},
		ReqHeaders:  nil,
		ReqBody:     nil,
		RespHeaders: nil,
		RespBody:    nil,
	}
}
