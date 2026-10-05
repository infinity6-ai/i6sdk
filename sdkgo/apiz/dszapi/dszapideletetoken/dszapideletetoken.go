package dszapideletetoken

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type DszDeleteTokenReq struct {
	Dataset string `json:"dataset"`
	Token   string `json:"token"`
}

type DszDeleteTokenApi struct {
	Req  *DszDeleteTokenReq `json:"req"`
	Resp *any               `json:"resp"`
}

func (d *DszDeleteTokenApi) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-delete-token",
		Method: "DELETE",
		Path:   "/api/datasetmanager/dataset/token/delete",
	}
}

func (d *DszDeleteTokenApi) GetDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: nil,
		QueryParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Req == nil {
					return nil
				}
				if d.Req == nil {
					d.Req = &DszDeleteTokenReq{}
				}
				return map[string]*schemaz.Schema{
					"dataset": {Str: schemaz.ParseStr(&d.Req.Dataset)},
					"token":   {Str: schemaz.ParseStr(&d.Req.Token)},
				}
			},
		},
		ReqHeaders:  nil,
		ReqBody:     nil,
		RespHeaders: nil,
		RespBody:    nil,
	}
}
