package genjsz

import (
	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstyp"
)

type Api struct {
	Spec ApiSpec   `json:"spec"`
	Req  jstyp.Typ `json:"req"`
	Resp jstyp.Typ `json:"resp"`
}

type ApiSpec struct {
	Id     string        `json:"id,omitempty"`
	Desc   *schemaz.Desc `json:"desc,omitempty"`
	Method string        `json:"method,omitempty"`
	Path   string        `json:"path,omitempty"`
}

func ParseApi(api apiz.Api) *Api {
	refs := api.GetDataRefs()

	parseSpec := func() ApiSpec {
		s := api.ApiSpec()
		ret := ApiSpec{
			Id:     s.Id,
			Method: s.Method,
			Path:   s.Path,
		}
		if s.Desc != nil {
			ret.Desc = s.Desc()
		}
		return ret
	}

	ret := &Api{
		Spec: parseSpec(),
		Req: jstyp.Typ{
			Js:     "object",
			Fields: make(map[string]jstyp.Typ),
		},
		Resp: jstyp.Typ{
			Js:     "object",
			Fields: make(map[string]jstyp.Typ),
		},
	}

	validateDataRef(refs.PathParams, "PathParams")
	validateDataRef(refs.QueryParams, "QueryParams")
	validateDataRef(refs.ReqHeaders, "ReqHeaders")
	validateDataRef(refs.ReqBody, "ReqBody")
	validateDataRef(refs.RespHeaders, "RespHeaders")
	validateDataRef(refs.RespBody, "RespBody")

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
