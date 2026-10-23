package dszapitablefileload

import (
	"encoding/json"
	"reflect"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/resultz"
	"go.code.infinity6.ai/sdkgo/sdkgo/resultz/relevancefashionresultz"
)

func AllApis() []apiz.Api {
	return []apiz.Api{
		&Api[*relevancefashionresultz.CatalogSimple]{},
	}
}

type Req[T resultz.Result] struct {
	Dataset lakedsz.Dataset `json:"dataset"`
	Query   any             `json:"query"`
}

func (r *Req[T]) UnmarshalJSON(data []byte) error {
	type Alias Req[T]
	var temp Alias
	var x T
	temp.Query = x.NewQuery()
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	*r = Req[T](temp)
	return nil
}

type Resp[T resultz.Result] struct {
	Data []T `json:"data"`
}

type Api[T resultz.Result] struct {
	Req  Req[T]  `json:"req"`
	Resp Resp[T] `json:"resp"`
}

func (d *Api[T]) ApiSpec() apiz.ApiSpec {
	var x T
	return apiz.ApiSpec{
		Id: "dsz-table-file-load-" + x.ResultName(),
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "Load Table File as JSON",
				Summary: "Load a table file from the datalake as JSON and return it.",
			}
		},
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/table-file-load/" + x.ResultName(),
	}
}

func (d *Api[T]) ApiDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {
						Desc: &schemaz.Desc{Summary: "Name of the dataset."},
						UrlValue:  schemaz.ParseStrGet(&d.Req.Dataset, lakedsz.MustParse),
					},
				}
			},
		},
		QueryParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Req.Query == nil {
					return nil
				}
				if d.Req.Query == nil {
					var x T
					d.Req.Query = x.NewQuery()
				}
				var x T
				return x.QuerySchema(d.Req.Query).Object(read)
			},
		},
		ReqHeaders:  nil,
		ReqBody:     nil,
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"data": {
						Desc: &schemaz.Desc{Summary: "Generic table data loaded from the datalake file."},
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]T, T]{
							Slice: func() *[]T { return &d.Resp.Data },
							NewElement: func(idx int) T {
								typ := reflect.TypeFor[T]()
								if typ.Kind() == reflect.Pointer {
									return reflect.New(typ.Elem()).Interface().(T)
								}
								var x T
								return x
							},
							Get: func(idx int, element *T) *schemaz.Schema {
								return (*element).Schema()
							},
						}),
					},
				}
			},
		},
	}
}