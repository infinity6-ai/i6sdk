package jsparser

import (
	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/camelstring"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

type Req struct {
	Params  jstypez.Typedef `json:"params"`
	Query   jstypez.Typedef `json:"query"`
	Headers jstypez.Typedef `json:"headers"`
	Payload jstypez.Typedef `json:"payload"`
}

type Resp struct {
	Headers jstypez.Typedef `json:"headers"`
	Payload jstypez.Typedef `json:"payload"`
}

type Api struct {
	Spec     ApiSpec          `json:"spec"`
	Req      jstypez.Typedef  `json:"-"`
	Resp     jstypez.Typedef  `json:"-"`
	Typedefs jstypez.Typedefs `json:"typedefs"`

	ReqV2  Req  `json:"req_v2"`
	RespV2 Resp `json:"resp_v2"`
}

func (a *Api) ReqResp() *jstypez.Typedefs {
	ret := &jstypez.Typedefs{}
	ret.Add(&a.Req, &a.Resp)
	return ret
}

type ApiSpec struct {
	Id      string        `json:"id,omitempty"`
	CamelId string        `json:"camel_id"`
	Desc    *schemaz.Desc `json:"desc,omitempty"`
	Method  string        `json:"method,omitempty"`
	Path    string        `json:"path,omitempty"`
}

func ParseApi(api apiz.Api) *Api {
	refs := api.GetDataRefs()

	apiId := camelstring.KebabToCamel(api.ApiSpec().Id)

	parseSpec := func() ApiSpec {
		s := api.ApiSpec()
		ret := ApiSpec{
			Id:      s.Id,
			CamelId: apiId,
			Method:  s.Method,
			Path:    s.Path,
		}
		if s.Desc != nil {
			ret.Desc = s.Desc()
		}
		return ret
	}

	ret := &Api{
		Spec: parseSpec(),
		// Req:  jstypez.Typedef{Properties: &jstypez.Tuples{}, Head: &jstypez.Tuple{Typ: "object", Name: apiId + "Req", Summary: apiId + " Request"}},
		// Resp: jstypez.Typedef{Properties: &jstypez.Tuples{}, Head: &jstypez.Tuple{Typ: "object", Name: apiId + "Resp", Summary: apiId + " Response"}},
	}

	validateDataRef(refs.PathParams, "PathParams")
	validateDataRef(refs.QueryParams, "QueryParams")
	validateDataRef(refs.ReqHeaders, "ReqHeaders")
	validateDataRef(refs.RespHeaders, "RespHeaders")

	parseRefsV2 := func(out *jstypez.Typedef, suffix string, s *schemaz.Schema) {
		if s == nil {
			return
		}
		upperFiledName := ret.Spec.CamelId + suffix
		out.Head = &jstypez.Tuple{Typ: "object", Name: upperFiledName}
		out.Properties = &jstypez.Tuples{}
		p, tdfs := ParseSchema(s)
		checker.Equal("object", p.Head.Typ, "must be object")
		out.Properties.Add(p.Properties.All()...)
		ret.Typedefs.ExtractFrom(tdfs)
	}

	// parseRefs := func(out *jstypez.Typedef, s *schemaz.Schema) {
	// 	if s == nil {
	// 		return
	// 	}
	// 	p, tdfs := ParseSchema(s)
	// 	checker.Equal("object", p.Head.Typ, "must be object")
	// 	out.Properties.Add(p.Properties.All()...)
	// 	ret.Typedefs.Add(tdfs.All()...)
	// }

	// parseIntoObject := func(input *schemaz.Schema, out *jstypez.Typedef, fieldName string) {
	// 	if input == nil {
	// 		return
	// 	}
	// 	p, tdfs := ParseSchema(input)
	// 	upperFiledName := strings.ToUpper(fieldName[0:1]) + fieldName[1:]
	// 	p.Head.Name = out.Head.Name + upperFiledName
	// 	out.Properties.Add(jstypez.Tuple{
	// 		Typ:     p.Head.Name,
	// 		Name:    fieldName,
	// 		Summary: p.Head.Summary,
	// 	})
	// 	ret.Typedefs.ExtractFrom(tdfs)
	// 	ret.Typedefs.Add(p)
	// }

	parseRefsV2(&ret.ReqV2.Params, "Params", refs.PathParams)
	parseRefsV2(&ret.ReqV2.Query, "Query", refs.QueryParams)
	parseRefsV2(&ret.ReqV2.Headers, "Headers", refs.ReqHeaders)
	parseRefsV2(&ret.ReqV2.Payload, "Payload", refs.ReqBody)
	parseRefsV2(&ret.RespV2.Headers, "Headers", refs.RespHeaders)
	parseRefsV2(&ret.RespV2.Payload, "Payload", refs.RespBody)

	// parseRefs(&ret.Req, refs.PathParams)
	// parseRefs(&ret.Req, refs.QueryParams)
	// parseRefs(&ret.Req, refs.ReqHeaders)
	// parseIntoObject(refs.ReqBody, &ret.Req, "payload")

	// parseRefs(&ret.Resp, refs.RespHeaders)
	// parseIntoObject(refs.RespBody, &ret.Resp, "payload")

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
