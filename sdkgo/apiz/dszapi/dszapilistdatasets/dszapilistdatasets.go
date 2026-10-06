package dszapilistdatasets

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type Solution struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Dataset struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Grants []string `json:"grants,omitempty"`
}

type DszListDatasetsReq struct {
	User string `json:"user"`
}

type DszListDatasetsResp struct {
	Datasets []*Dataset `json:"datasets"`
}

type DszListDatasetsApi struct {
	Req  *DszListDatasetsReq  `json:"req"`
	Resp *DszListDatasetsResp `json:"resp"`
}

func (d *DszListDatasetsApi) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-list-datasets",
		Desc:   nil,
		Method: "GET",
		Path:   "/api/dsz/datasets",
	}
}

func (d *DszListDatasetsApi) GetDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: nil,
		QueryParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Req == nil {
					return nil
				}
				if d.Req == nil {
					d.Req = &DszListDatasetsReq{}
				}
				return map[string]*schemaz.Schema{
					"user": {Str: schemaz.ParseStr(&d.Req.User)},
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
					d.Resp = &DszListDatasetsResp{}
				}
				return map[string]*schemaz.Schema{
					"datasets": {
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]*Dataset, *Dataset]{
							Slice:      func() *[]*Dataset { return &d.Resp.Datasets },
							NewElement: func(idx int) *Dataset { return &Dataset{} },
							Get: func(idx int, element **Dataset) *schemaz.Schema {
								return &schemaz.Schema{
									Object: func(_ bool) map[string]*schemaz.Schema {
										return map[string]*schemaz.Schema{
											"id": {
												Raw: func() any { return &(*element).Id },
											},
											"name": {
												Raw: func() any { return &(*element).Name },
											},
											"grants": {
												Array: schemaz.NewArray(schemaz.ArrayOptions[[]string, string]{
													Slice: func() *[]string { return &(*element).Grants },
													Get: func(idx int, grant *string) *schemaz.Schema {
														return &schemaz.Schema{
															Raw: func() any { return grant },
														}
													},
												}),
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
