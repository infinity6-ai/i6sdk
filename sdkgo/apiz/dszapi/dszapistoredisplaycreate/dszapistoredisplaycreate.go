package dszapistoredisplaycreate

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/dszstore"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/storedisplayid"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/storeid"
)

type DszStoreDisplayCreateReq struct {
	Dataset lakedsz.Dataset          `json:"dataset"`
	StoreId storeid.StoreId          `json:"store_id"`
	Display dszstore.DszStoreDisplay `json:"display"`
}

type DszStoreDisplayCreateResp struct {
	Display dszstore.DszStoreDisplay `json:"display"`
}

type Api struct {
	Req  *DszStoreDisplayCreateReq  `json:"req"`
	Resp *DszStoreDisplayCreateResp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-store-display-create",
		Desc:   nil,
		Method: "POST",
		Path:   "/api/ds/dataset/{dataset}/store/display/create",
	}
}

func (d *Api) ApiDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Req == nil {
					return nil
				}
				if d.Req == nil {
					d.Req = &DszStoreDisplayCreateReq{}
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
					d.Req = &DszStoreDisplayCreateReq{}
				}
				return map[string]*schemaz.Schema{
					"store_id": {
						Str: schemaz.ParseStrGet(&d.Req.StoreId, storeid.MustParse),
					},
					"display": {
						Object: func(read bool) map[string]*schemaz.Schema {
							return map[string]*schemaz.Schema{
								"id": {
									Str: schemaz.ParseStrGet(&d.Req.Display.Id, storedisplayid.MustParse),
								},
								"name": {
									Raw: func() any { return &d.Req.Display.Name },
								},
							}
						},
					},
				}
			},
		},
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Resp == nil {
					return nil
				}
				if d.Resp == nil {
					d.Resp = &DszStoreDisplayCreateResp{}
				}
				return map[string]*schemaz.Schema{
					"display": {
						Object: func(read bool) map[string]*schemaz.Schema {
							return map[string]*schemaz.Schema{
								"id": {
									Str: schemaz.ParseStrGet(&d.Resp.Display.Id, storedisplayid.MustParse),
								},
								"name": {
									Raw: func() any { return &d.Resp.Display.Name },
								},
							}
						},
					},
				}
			},
		},
	}
}
