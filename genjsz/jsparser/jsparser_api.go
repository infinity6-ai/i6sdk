package jsparser

import (
	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

type Api struct {
	Spec     ApiSpec            `json:"spec"`
	Req      jstypez.Typedef    `json:"req"`
	Resp     jstypez.Typedef    `json:"resp"`
	Typedefs jstypez.TypedefsV2 `json:"typedefs"`
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
		Req:  jstypez.Typedef{Typ: "object"},
		Resp: jstypez.Typedef{Typ: "object"},
	}

	validateDataRef(refs.PathParams, "PathParams")
	validateDataRef(refs.QueryParams, "QueryParams")
	validateDataRef(refs.ReqHeaders, "ReqHeaders")
	validateDataRef(refs.RespHeaders, "RespHeaders")

	parseRefs := func(out *jstypez.Typedef, s *schemaz.Schema) {
		if s == nil {
			return
		}
		p, _ := ParseSchema(s)
		checker.Equal("object", p.Typ, "must be object")
		for _, v := range p.Properties {
			out.SetProperty(v)
		}
	}

	parseRefs(&ret.Req, refs.PathParams)
	parseRefs(&ret.Req, refs.QueryParams)
	parseRefs(&ret.Req, refs.ReqHeaders)
	if refs.ReqBody != nil {
		// p := ParseSchema(refs.ReqBody)
		ret.Req.SetProperty(jstypez.Property{
			Typ:  "object",
			Name: "payload",
		})
	}
	parseRefs(&ret.Resp, refs.RespHeaders)
	if refs.RespBody != nil {
		// p := ParseSchema(refs.RespBody)
		ret.Resp.SetProperty(jstypez.Property{
			Typ:  "object",
			Name: "payload",
		})
	}

	return ret
}

func validateDataRef(schema *schemaz.Schema, n string) {
	if schema != nil {
		checker.NotNil(schema.Object, "%s", n)
	}
}

func descSummary(d *schemaz.Desc) string {
	if d == nil {
		return ""
	}
	return d.Summary
}
