package jssdkfunc

import (
	"context"
	"fmt"
	"strings"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/gox/commonz/stringz/camelz"
	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jsschema"
	"github.com/infinity6-ai/gox/schemaz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

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
	Version  string        `json:"version,omitempty"`
	Desc     *schemaz.Desc `json:"desc,omitempty"`
	Typ      sdkfuncz.Typ  `json:"typ,omitempty"`
	FuncName string        `json:"func_name,omitempty"`
}

type SdkFuncs struct {
	SdkFuncs []*SdkFunc       `json:"sdk_funcs"`
	Typedefs jstypez.Typedefs `json:"typedefs"`
}

func (s *SdkFuncs) Funcs() []*SdkFunc {
	return s.SdkFuncs
}

type Parser struct {
	schemaParser *jsschema.Parser
}

func New(schemaParser *jsschema.Parser) *Parser {
	if schemaParser == nil {
		schemaParser = jsschema.NewSdk()
	}
	return &Parser{schemaParser: schemaParser}
}

func NewSdk() *Parser {
	return New(jsschema.NewSdk())
}

func (p *Parser) ParseSdkFunc(ctx context.Context, fn sdkfuncz.Func) *SdkFunc {
	refs := fn.SdkFuncDataRefs()

	funcId := ""
	if id := fn.SdkFuncSpec().Id; id != "" {
		funcId = "sdkFunc" + camelz.P(id).P()
	}

	parseSpec := func() SdkFuncSpec {
		s := fn.SdkFuncSpec()
		ret := SdkFuncSpec{
			Id:      s.Id,
			CamelId: funcId,
			Version: s.Version,
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
		validateDataRef(refs.Partitions, "Partitions")
		validateDataRef(refs.Output, "Output")

		parseRefsV2 := func(out *jstypez.Typedef, suffix string, s *schemaz.Schema) {
			if s == nil {
				return
			}
			fieldName := ret.Spec.CamelId + strings.ToUpper(suffix[0:1]) + suffix[1:]
			td, tdfs := p.schemaParser.ParseSchema(ctx, fieldName, s)
			checker.Equal("object", td.Head.Typ, "must be object")
			out.Properties = &jstypez.Tuples{}
			out.Head = &jstypez.Tuple{Typ: fieldName, Name: suffix}
			td.Head.Summary = fmt.Sprintf("%s %s", ret.Spec.CamelId, suffix)
			ret.Typedefs.ExtractFrom(tdfs)
			ret.Typedefs.Add(td)
		}

		parseRefsV2(&ret.Req.Input, "input", refs.Input)
		parseRefsV2(&ret.Req.Partitions, "partitions", refs.Partitions)
		parseRefsV2(&ret.Resp.Output, "output", refs.Output)
	}

	return ret
}

func (p *Parser) ParseSdkFuncs(ctx context.Context, funcs []sdkfuncz.Func) *SdkFuncs {
	ret := SdkFuncs{
		SdkFuncs: make([]*SdkFunc, len(funcs)),
		Typedefs: jstypez.NewTypedefs(nil),
	}
	for idx, fn := range funcs {
		logger.Info(ctx, "parsing js sdk func", map[string]any{"func": fn.SdkFuncSpec().Id})
		parsed := p.ParseSdkFunc(ctx, fn)
		ret.Typedefs.Merge(parsed.Typedefs)
		ret.SdkFuncs[idx] = parsed
	}
	return &ret
}

func validateDataRef(schema *schemaz.Schema, n string) {
	if schema != nil {
		checker.NotNil(schema.Object, "%s", n)
	}
}
