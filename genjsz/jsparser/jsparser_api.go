package jsparser

import (
	"strings"

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

	parseIntoObject := func(input *schemaz.Schema, out *jstypez.Typedef, fieldName string) {
		if input == nil {
			return
		}
		p, tdfs := ParseSchema(input)
		upperFiledName := strings.ToUpper(fieldName[0:1]) + fieldName[1:]
		p.Head.Name = out.Head.Name + upperFiledName
		out.Properties.Add(jstypez.Tuple{
			Typ:     p.Head.Name,
			Name:    fieldName,
			Summary: p.Head.Summary,
		})
		ret.Typedefs.ExtractFrom(tdfs)
		ret.Typedefs.Add(p)
	}

	parseRefs(&ret.Req, refs.PathParams)
	parseRefs(&ret.Req, refs.QueryParams)
	parseRefs(&ret.Req, refs.ReqHeaders)
	parseIntoObject(refs.ReqBody, &ret.Req, "payload")
	parseRefs(&ret.Resp, refs.RespHeaders)
	parseIntoObject(refs.RespBody, &ret.Resp, "payload")

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
