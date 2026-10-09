package jsparsersdkfunc

import (
	"fmt"
	"strings"

	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/camelstring"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
)

type Req struct {
	Input      jstypez.Typedef `json:"input"`
	Partitions jstypez.Typedef `json:"partitions"`
}

type Resp struct {
	Output jstypez.Typedef `json:"output"`
}

type SdkFunc struct {
	Spec     SdkFuncSpec      `json:"spec"`
	Typedefs jstypez.Typedefs `json:"typedefs"`

	Req  Req  `json:"req"`
	Resp Resp `json:"resp"`
}

type SdkFuncSpec struct {
	Id       string        `json:"id,omitempty"`
	CamelId  string        `json:"camel_id"`
	Desc     *schemaz.Desc `json:"desc,omitempty"`
	Typ      sdkfuncz.Typ  `json:"typ,omitempty"`
	FuncName string        `json:"func_name,omitempty"`
}

func ParseSdkFunc(fn sdkfuncz.Func) *SdkFunc {
	refs := fn.SdkFuncDataRefs()

	funcId := camelstring.ToCamel(fn.SdkFuncSpec().Id)
	if funcId != "" {
		funcId = "sdkFunc" + strings.ToUpper(funcId[0:1]) + funcId[1:]
	}

	parseSpec := func() SdkFuncSpec {
		s := fn.SdkFuncSpec()
		ret := SdkFuncSpec{
			Id:      s.Id,
			CamelId: funcId,
			Typ:     s.Typ,
		}
		if s.Desc != nil {
			ret.Desc = s.Desc()
		}
		return ret
	}

	ret := &SdkFunc{
		Spec: parseSpec(),
	}

	xxx := jstypez.BasicTypes()
	ret.Typedefs.ExtractFrom(xxx)

	if refs != nil {
		validateDataRef(refs.Input, "Input")
		validateDataRef(refs.ModelPartitions, "Partitions")
		validateDataRef(refs.Output, "Output")

		parseRefsV2 := func(out *jstypez.Typedef, suffix string, s *schemaz.Schema) {
			if s == nil {
				return
			}
			fieldName := ret.Spec.CamelId + strings.ToUpper(suffix[0:1]) + suffix[1:]
			p, tdfs := jsparser.ParseSchema(fieldName, s)
			checker.Equal("object", p.Head.Typ, "must be object")
			out.Properties = &jstypez.Tuples{}
			out.Head = &jstypez.Tuple{Typ: fieldName, Name: suffix}
			p.Head.Summary = fmt.Sprintf("%s %s", ret.Spec.CamelId, suffix)
			ret.Typedefs.ExtractFrom(tdfs)
			ret.Typedefs.Add(p)
		}

		parseRefsV2(&ret.Req.Input, "input", refs.Input)
		parseRefsV2(&ret.Req.Partitions, "partitions", refs.ModelPartitions)
		parseRefsV2(&ret.Resp.Output, "output", refs.Output)
	}

	return ret
}

func validateDataRef(schema *schemaz.Schema, n string) {
	if schema != nil {
		checker.NotNil(schema.Object, "%s", n)
	}
}
