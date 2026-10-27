package ingestzapievent

import (
	"reflect"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/tabledataz"
	"go.code.infinity6.ai/sdkgo/sdkgo/tabledataz/relevancefashiondataz"
)

func AllApis() []apiz.Api {
	return []apiz.Api{
		&Api[*relevancefashiondataz.Catalog]{},
		&Api[*relevancefashiondataz.Event]{},
	}
}

type Req[T tabledataz.TableData] struct {
	Dataset lakedsz.Dataset
	Suffix  string
	Payload T
}

type Resp struct {
	Count int
}

type Api[T tabledataz.TableData] struct {
	Req  Req[T]
	Resp Resp
}

func (d *Api[T]) ApiSpec() apiz.ApiSpec {
	var x T
	return apiz.ApiSpec{
		Id: "ingestz-event-" + x.Table().String(),
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "Event Ingestion",
				Summary: "Ingests events data into a target dataset table.",
			}
		},
		Method: "GET",
		Path:   "/api/ingest/dataset/{dataset}/event/" + x.Table().String() + "/text.txt",
	}
}

func (d *Api[T]) ApiDataRefs() *apiz.DataRefs {
	val := reflect.ValueOf(&d.Req.Payload).Elem()
	if val.Kind() == reflect.Pointer && val.IsNil() {
		val.Set(reflect.New(val.Type().Elem()))
	}
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {
						Desc:     &schemaz.Desc{Summary: "Name of the dataset that receives the stream data."},
						UrlValue: schemaz.ParseStrGet(&d.Req.Dataset, lakedsz.MustParse),
					},
				}
			},
		},
		QueryParams: d.Req.Payload.Schema(),
		ReqHeaders:  nil,
		ReqBody:     nil,
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"count": {Field: func() any { return &d.Resp.Count }},
				}
			},
		},
	}
}
