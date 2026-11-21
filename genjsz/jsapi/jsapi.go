package jsapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/gox/commonz/stringz/camelz"
	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jsschema"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

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

type Apis struct {
	Apis     []*Api           `json:"apis"`
	Typedefs jstypez.Typedefs `json:"typedefs"`
}

func ParseApi(ctx context.Context, api apiz.Api) *Api {
	refs := api.ApiDataRefs()

	apiId := ""
	if id := api.ApiSpec().Id; id != "" {
		apiId = "api" + camelz.P(id).P()
	}

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
	}

	xxx := jstypez.BasicTypes()
	ret.Typedefs.ExtractFrom(xxx)

	validateDataRef(refs.PathParams, "PathParams")
	validateDataRef(refs.QueryParams, "QueryParams")
	validateDataRef(refs.ReqHeaders, "ReqHeaders")
	validateDataRef(refs.RespHeaders, "RespHeaders")

	parseRefsV2 := func(out *jstypez.Typedef, suffix string, s *schemaz.Schema) {
		if s == nil {
			return
		}
		fieldName := ret.Spec.CamelId + strings.ToUpper(suffix[0:1]) + suffix[1:]
		p, tdfs := jsschema.ParseSchema(ctx, fieldName, s)
		checker.Equal("object", p.Head.Typ, "must be object")
		out.Properties = &jstypez.Tuples{}
		out.Head = &jstypez.Tuple{Typ: fieldName, Name: suffix}
		p.Head.Summary = fmt.Sprintf("%s %s", ret.Spec.CamelId, suffix)
		ret.Typedefs.ExtractFrom(tdfs)
		ret.Typedefs.Add(p)
	}

	parseRefsV2(&ret.Req.Params, "params", refs.PathParams)
	parseRefsV2(&ret.Req.Query, "query", refs.QueryParams)
	parseRefsV2(&ret.Req.Headers, "reqHeaders", refs.ReqHeaders)
	parseRefsV2(&ret.Req.Payload, "reqPayload", refs.ReqBody)
	parseRefsV2(&ret.Resp.Headers, "respHeaders", refs.RespHeaders)
	parseRefsV2(&ret.Resp.Payload, "respPayload", refs.RespBody)

	return ret
}

func ParseApis(ctx context.Context, apis []apiz.Api) *Apis {
	ret := Apis{
		Apis:     make([]*Api, len(apis)),
		Typedefs: jstypez.NewTypedefs(nil),
	}
	for idx, api := range apis {
		logger.Info(ctx, "parsring js sdk api", map[string]any{"api": api.ApiSpec().Id})
		parsed := ParseApi(ctx, api)
		ret.Typedefs.Merge(parsed.Typedefs)
		ret.Apis[idx] = parsed
	}
	return &ret
}

func validateDataRef(schema *schemaz.Schema, n string) {
	if schema != nil {
		checker.NotNil(schema.Object, "%s", n)
	}
}
