package jsparser

import (
	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/camelstring"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

type Api struct {
	Spec     ApiSpec          `json:"spec"`
	Req      jstypez.Typedef  `json:"req"`
	Resp     jstypez.Typedef  `json:"resp"`
	Typedefs jstypez.Typedefs `json:"typedefs"`
}

func (a *Api) ReqResp() *jstypez.Typedefs {
	ret := &jstypez.Typedefs{}
	ret.Add(&a.Req, &a.Resp)
	return ret
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

	apiId := camelstring.KebabToCamel(api.ApiSpec().Id)

	ret := &Api{
		Spec: parseSpec(),
		Req:  jstypez.Typedef{Properties: &jstypez.Tuples{}, Head: &jstypez.Tuple{Typ: "object", Name: apiId + "Req", Summary: apiId + " Request"}},
		Resp: jstypez.Typedef{Properties: &jstypez.Tuples{}, Head: &jstypez.Tuple{Typ: "object", Name: apiId + "Resp", Summary: apiId + " Response"}},
	}

	validateDataRef(refs.PathParams, "PathParams")
	validateDataRef(refs.QueryParams, "QueryParams")
	validateDataRef(refs.ReqHeaders, "ReqHeaders")
	validateDataRef(refs.RespHeaders, "RespHeaders")

	parseRefs := func(out *jstypez.Typedef, s *schemaz.Schema) {
		if s == nil {
			return
		}
		p, tdfs := ParseSchema(s)
		checker.Equal("object", p.Head.Typ, "must be object")
		out.Properties.Add(p.Properties.All()...)
		ret.Typedefs.Add(tdfs.All()...)
	}

	parseRefs(&ret.Req, refs.PathParams)
	parseRefs(&ret.Req, refs.QueryParams)
	parseRefs(&ret.Req, refs.ReqHeaders)
	if refs.ReqBody != nil {
		p, tdfs := ParseSchema(refs.ReqBody)
		p.Head.Name = ret.Req.Head.Name + "Payload"
		ret.Req.Properties.Add(jstypez.Tuple{
			Typ:     p.Head.Name,
			Name:    "payload",
			Summary: p.Head.Summary,
		})
		ret.Typedefs.Add(tdfs.All()...)
		ret.Typedefs.Add(p)
	}
	parseRefs(&ret.Resp, refs.RespHeaders)
	if refs.RespBody != nil {
		p, tdfs := ParseSchema(refs.RespBody)
		p.Head.Name = ret.Resp.Head.Name + "Payload"
		ret.Resp.Properties.Add(jstypez.Tuple{
			Typ:     p.Head.Name,
			Name:    "payload",
			Summary: p.Head.Summary,
		})
		ret.Typedefs.ExtractFrom(tdfs)
		ret.Typedefs.Add(p)
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
