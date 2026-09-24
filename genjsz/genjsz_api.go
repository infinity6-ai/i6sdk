package genjsz

import (
	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstyp"
)

func ParseApi(api apiz.Api) jstyp.Typ {
	// return ParseSchema(api.GetDataRefs().PathParams)
	refs := api.GetDataRefs()

	validateDataRef(refs.PathParams, "PathParams")
	validateDataRef(refs.QueryParams, "QueryParams")
	validateDataRef(refs.ReqHeaders, "ReqHeaders")
	validateDataRef(refs.ReqBody, "ReqBody")
	validateDataRef(refs.RespHeaders, "RespHeaders")
	validateDataRef(refs.RespBody, "RespBody")

	ret := jstyp.Typ{
		Js:     "object",
		Fields: make(map[string]jstyp.Typ),
	}

	parseApi := func(s *schemaz.Schema) {
		p := ParseSchema(s)
		checker.Equal("object", p.Js, "must be object")
		for k, v := range p.Fields {
			old, ok := ret.Fields[k]
			checker.False(ok, "field colision %s, old: %s, new: %s", k, old, v)
			ret.Fields[k] = v
		}
	}

	parseApi(refs.PathParams)
	parseApi(refs.QueryParams)
	parseApi(refs.ReqHeaders)
	parseApi(refs.ReqBody)
	parseApi(refs.RespHeaders)
	parseApi(refs.RespBody)

	return ret
}

func validateDataRef(schema *schemaz.Schema, n string) {
	if schema != nil {
		checker.NotNil(schema.Object, "%s", n)
	}
}
