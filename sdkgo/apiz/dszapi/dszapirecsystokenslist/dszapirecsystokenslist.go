package dszapirecsystokenslist

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/dsztoken"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

type DszListRecsysTokensReq struct {
	Dataset lakedsz.Dataset `json:"dataset"`
}

type DszListRecsysTokensResp []dsztoken.DatasetToken

func (*DszListRecsysTokensResp) JSType() string {
	return jstypez.TypObject
}

type DszListRecsysTokensApi struct {
	Req  DszListRecsysTokensReq  `json:"req"`
	Resp DszListRecsysTokensResp `json:"resp"`
}

func (d *DszListRecsysTokensApi) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-list-recsys-tokens",
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/recsystoken/list",
	}
}

func (d *DszListRecsysTokensApi) GetDataRefs() *apiz.DataRefs {
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
		ReqBody:     nil,
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Raw: func() any { return &d.Resp },
		},
	}
}
