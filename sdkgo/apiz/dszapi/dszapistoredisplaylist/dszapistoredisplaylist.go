package dszapistoredisplaylist

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/dszstore"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/storedisplayid"
)

type DszStoreDisplayListReq struct {
	Dataset lakedsz.Dataset `json:"dataset"`
	StoreId string          `json:"store_id"`
}

type DszStoreDisplayListResp struct {
	Dataset  lakedsz.Dataset            `json:"dataset"`
	StoreId  string                     `json:"store_id"`
	Displays []dszstore.DszStoreDisplay `json:"displays"`
}

type Api struct {
	Req  *DszStoreDisplayListReq  `json:"req"`
	Resp *DszStoreDisplayListResp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-store-display-list",
		Desc:   nil,
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/store/display/list",
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
					d.Req = &DszStoreDisplayListReq{}
				}
				return map[string]*schemaz.Schema{
					"dataset": {Str: lakedsz.ParseStr(&d.Req.Dataset)},
				}
			},
		},
		QueryParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Req == nil {
					return nil
				}
				if d.Req == nil {
					d.Req = &DszStoreDisplayListReq{}
				}
				return map[string]*schemaz.Schema{
					"store_id": {Str: schemaz.ParseStr(&d.Req.StoreId)},
				}
			},
		},
		ReqHeaders:  nil,
		ReqBody:     nil,
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Resp == nil {
					return nil
				}
				if d.Resp == nil {
					d.Resp = &DszStoreDisplayListResp{}
				}
				if d.Resp.Displays == nil {
					d.Resp.Displays = []dszstore.DszStoreDisplay{}
				}
				return map[string]*schemaz.Schema{
					"dataset": {
						Raw: func() any { return &d.Resp.Dataset },
					},
					"store_id": {
						Raw: func() any { return &d.Resp.StoreId },
					},
					"displays": {
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]dszstore.DszStoreDisplay, dszstore.DszStoreDisplay]{
							Slice: func() *[]dszstore.DszStoreDisplay { return &d.Resp.Displays },
							Get: func(idx int, element *dszstore.DszStoreDisplay) *schemaz.Schema {
								return &schemaz.Schema{
									Object: func(_ bool) map[string]*schemaz.Schema {
										return map[string]*schemaz.Schema{
											"id": {
												Str: schemaz.ParseStrGet(&element.Id, storedisplayid.MustParse),
											},
											"name": {
												Raw: func() any { return &element.Name },
											},
										}
									},
								}
							},
						}),
					},
				}
			},
		},
	}
}
