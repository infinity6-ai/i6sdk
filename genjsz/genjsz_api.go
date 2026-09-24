package genjsz

import (
	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstyp"
)

type Api struct {
	Req  jstyp.Typ `json:"req"`
	Resp jstyp.Typ `json:"resp"`
}

func ParseApi(api apiz.Api) *Api {
	// return ParseSchema(api.GetDataRefs().PathParams)
	refs := api.GetDataRefs()

	validateDataRef(refs.PathParams, "PathParams")
	validateDataRef(refs.QueryParams, "QueryParams")
	validateDataRef(refs.ReqHeaders, "ReqHeaders")
	validateDataRef(refs.ReqBody, "ReqBody")
	validateDataRef(refs.RespHeaders, "RespHeaders")
	validateDataRef(refs.RespBody, "RespBody")

	ret := &Api{
		Req: jstyp.Typ{
			Js:     "object",
			Fields: make(map[string]jstyp.Typ),
		},
		Resp: jstyp.Typ{
			Js:     "object",
			Fields: make(map[string]jstyp.Typ),
		},
	}

	parseRefs := func(out jstyp.Typ, schemas ...*schemaz.Schema) {
		for _, s := range schemas {
			p := ParseSchema(s)
			checker.Equal("object", p.Js, "must be object")
			for k, v := range p.Fields {
				old, ok := out.Fields[k]
				checker.False(ok, "field colision %s, old: %s, new: %s", k, old, v)
				out.Fields[k] = v
			}
		}
	}

	parseRefs(ret.Req, refs.PathParams)
	parseRefs(ret.Req, refs.QueryParams)
	parseRefs(ret.Req, refs.ReqHeaders)
	parseRefs(ret.Req, refs.ReqBody)
	parseRefs(ret.Resp, refs.RespHeaders)
	parseRefs(ret.Resp, refs.RespBody)

	return ret
}

func validateDataRef(schema *schemaz.Schema, n string) {
	if schema != nil {
		checker.NotNil(schema.Object, "%s", n)
	}
}
