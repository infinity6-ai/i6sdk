package jsparser

import (
	"strings"

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
	Typedefs jstypez.Typedefs `json:"typedefs"`

	Req  Req  `json:"req"`
	Resp Resp `json:"resp"`
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
	apiId = "api" + strings.ToUpper(apiId[0:1]) + apiId[1:]

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
		fieldName := ret.Spec.CamelId + suffix[0:1] + suffix[1:]
		// out.Head = &jstypez.Tuple{Typ: "object", Name: fieldName}
		out.Properties = &jstypez.Tuples{}
		p, tdfs := ParseSchema(s)
		checker.Equal("object", p.Head.Typ, "must be object")
		out.Head = &jstypez.Tuple{Typ: fieldName, Name: suffix}
		p.Head.Name = out.Head.Typ
		// out.Properties.Add(p.Properties.All()...)
		ret.Typedefs.ExtractFrom(tdfs)
		ret.Typedefs.Add(p)
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

	parseRefsV2(&ret.Req.Params, "params", refs.PathParams)
	parseRefsV2(&ret.Req.Query, "query", refs.QueryParams)
	parseRefsV2(&ret.Req.Headers, "reqHeaders", refs.ReqHeaders)
	parseRefsV2(&ret.Req.Payload, "reqPayload", refs.ReqBody)
	parseRefsV2(&ret.Resp.Headers, "respHeaders", refs.RespHeaders)
	parseRefsV2(&ret.Resp.Payload, "respPayload", refs.RespBody)

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
