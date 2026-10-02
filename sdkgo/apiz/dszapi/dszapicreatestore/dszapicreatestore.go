package dszapicreatestore

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

type DszStore struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (*DszStore) JSTypedef() (*jstypez.Typedef, *jstypez.Typedefs) {
	td := &jstypez.Typedef{
		Head:       &jstypez.Tuple{Typ: jstypez.TypObject, Name: ""},
		Properties: &jstypez.Tuples{},
	}
	td.Properties.Add(
		jstypez.Tuple{Typ: jstypez.TypString, Name: "id"},
		jstypez.Tuple{Typ: jstypez.TypString, Name: "name"},
	)
	return td, nil
}

type DszCreateStoreReq struct {
	Dataset lakedsz.Dataset `json:"dataset"`
	Store   DszStore        `json:"store"`
}

type DszCreateStoreApi struct {
	Req *DszCreateStoreReq `json:"req"`
}

func (d *DszCreateStoreApi) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-create-store",
		Desc:   nil,
		Method: "POST",
		Path:   "/api/ds/dataset/{dataset}/store/create",
	}
}

func (d *DszCreateStoreApi) GetDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Req == nil {
					return nil
				}
				if d.Req == nil {
					d.Req = &DszCreateStoreReq{}
				}
				return map[string]*schemaz.Schema{
					"dataset": {Str: lakedsz.ParseStr(&d.Req.Dataset)},
				}
			},
		},
		QueryParams: nil,
		ReqHeaders:  nil,
		ReqBody: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Req == nil {
					return nil
				}
				if d.Req == nil {
					d.Req = &DszCreateStoreReq{}
				}
				return map[string]*schemaz.Schema{
					"store": {
						Raw: func() any { return &d.Req.Store },
					},
				}
			},
		},
		RespHeaders: nil,
		RespBody:    nil,
	}
}
