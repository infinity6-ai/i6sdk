package forgepysdkfunc

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"github.com/infinity6-ai/gox/schemaz/schemamold"
	"github.com/infinity6-ai/gox/schemaz/typeastz"
	"github.com/infinity6-ai/i6sdk/forgez/forgepy/forgepyschema"
	"github.com/infinity6-ai/i6sdk/mold/moldsdkfunc"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
)

func NewSdk() *moldsdkfunc.Parser {
	return moldsdkfunc.New(forgepyschema.NewSdk())
}

func writeln(sb *strings.Builder, msg string, args ...any) {
	_, err := fmt.Fprintf(sb, msg, args...)
	errorz.Check(err)
	_, err = sb.WriteRune('\n')
	errorz.Check(err)
}

func indent(s, prefix string) string {
	lines := strings.Split(strings.TrimSuffix(s, "\n"), "\n")
	for i, l := range lines {
		if l != "" {
			lines[i] = prefix + l
		}
	}
	return strings.Join(lines, "\n")
}

func ExtractTypedefs(sdkFuncs []*moldsdkfunc.SdkFunc) typeastz.Typedefs {
	tdfs := typeastz.Typedefs{}
	for _, fn := range sdkFuncs {
		funcId := fn.Spec.CamelId

		hasInput := fn.Req.Input.Head != nil
		hasPartitions := fn.Req.Partitions.Head != nil
		hasReq := hasInput || hasPartitions

		if hasReq {
			reqTypedef := typeastz.Typedef{
				Head:       &typeastz.Tuple{Typ: "object", Name: funcId + "Req", Summary: funcId + " Request"},
				Properties: &typeastz.Tuples{},
			}
			if hasInput {
				reqTypedef.Properties.Add(typeastz.Tuple{Typ: fn.Req.Input.Head.Typ, GoTyp: fn.Req.Input.Head.GoTyp, Name: "input"})
			}
			if hasPartitions {
				reqTypedef.Properties.Add(typeastz.Tuple{Typ: fn.Req.Partitions.Head.Typ, GoTyp: fn.Req.Partitions.Head.GoTyp, Name: "partitions"})
			}
			tdfs.Add(&reqTypedef)
		}

		hasOutput := fn.Resp.Output.Head != nil
		respTypedef := typeastz.Typedef{
			Head:       &typeastz.Tuple{Typ: "object", Name: funcId + "Resp", Summary: funcId + " Response"},
			Properties: &typeastz.Tuples{},
		}
		if hasOutput {
			respTypedef.Properties.Add(typeastz.Tuple{Typ: fn.Resp.Output.Head.Typ, GoTyp: fn.Resp.Output.Head.GoTyp, Name: "output"})
		}
		tdfs.Add(&respTypedef)

		if fn.Spec.Typ == sdkfuncz.TypClientInfer {
			enriched, selected := buildSelectedTypedefs(fn)
			tdfs.Add(&enriched, &selected)
		}
	}
	return tdfs
}

func buildSelectedTypedefs(fn *moldsdkfunc.SdkFunc) (typeastz.Typedef, typeastz.Typedef) {
	funcId := fn.Spec.CamelId
	enrichedName := funcId + "OutputResultsElementSkusElementEnriched"
	enriched := typeastz.Typedef{
		Head: &typeastz.Tuple{Typ: "object", Name: enrichedName, Summary: "One recommended sku to show, as picked by selectResults"},
		Properties: typeastz.NewTuples([]typeastz.Tuple{
			{Typ: "string", Name: "sku_id", Summary: "Id of the sku. Required; the only field sent to i6 (in the recshow event)"},
			{Typ: "number", Name: "rating", Summary: "Model rating, higher is more relevant. The default selectResults sets the rating of the result the sku was picked from"},
			{Typ: "*", Name: "enriched", Summary: "Optional: any data your selectResults attaches to the sku (e.g. name, price, image). Returned to you as is, never sent to i6. The default selectResults does not set it"},
		}),
	}
	selected := typeastz.Typedef{
		Head: &typeastz.Tuple{Typ: "object", Name: funcId + "Selected", Summary: funcId + " Selected: the skus to show, as returned by selectResults"},
		Properties: typeastz.NewTuples([]typeastz.Tuple{
			{Typ: fmt.Sprintf("Array<%s>", enrichedName), Name: "results", Summary: "Skus to show, in the order selectResults returned them"},
		}),
	}
	return enriched, selected
}

func toSnakeCase(s string) string {
	var result strings.Builder
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if unicode.IsUpper(r) {
			if i > 0 && (unicode.IsLower(runes[i-1]) || (i+1 < len(runes) && unicode.IsLower(runes[i+1]))) {
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func funcName(funcId string, version string) string {
	name := strings.TrimPrefix(funcId, "sdkFunc")
	res := toSnakeCase(name)
	if version != "" {
		res += "_" + toSnakeCase(version)
	}
	return res
}

func hasInputParams(fn *moldsdkfunc.SdkFunc) bool {
	if fn.Req.Input.Head == nil {
		return false
	}
	for _, td := range fn.Typedefs.All() {
		if td.Head.Name != fn.Req.Input.Head.Typ || td.Properties == nil {
			continue
		}
		for _, prop := range td.Properties.All() {
			if prop.Name == "params" {
				return true
			}
		}
	}
	return false
}

func Write(ctx context.Context, sdkFuncs []*moldsdkfunc.SdkFunc) string {
	if len(sdkFuncs) == 0 {
		return ""
	}

	var b strings.Builder

	writeln(&b, "class SdkFuncs:")
	writeln(&b, "    \"\"\"Auto-generated SDK workflow functions.\"\"\"")
	writeln(&b, "")
	writeln(&b, "    def __init__(self, services: Any) -> None:")
	writeln(&b, "        self._services = services")

	for _, fn := range sdkFuncs {
		writeln(&b, "")
		var decl strings.Builder
		funcId := fn.Spec.CamelId
		methodName := funcName(funcId, fn.Spec.Version)

		hasInput := fn.Req.Input.Head != nil
		hasPartitions := fn.Req.Partitions.Head != nil
		hasReq := hasInput || hasPartitions

		isInfer := fn.Spec.Typ == sdkfuncz.TypClientInfer

		if isInfer {
			writeln(&decl, "def %s(", methodName)
			writeln(&decl, "    self,")
			writeln(&decl, "    req: %sReq,", funcId)
			writeln(&decl, "    select_results: Optional[Callable[[%sReq, %sOutput], %sSelected]] = None,", funcId, funcId, funcId)
			writeln(&decl, "    opts: Optional[ServiceOptions] = None,")
			writeln(&decl, ") -> %sSelected:", funcId)
		} else if fn.Spec.Typ == sdkfuncz.TypClientQuery {
			writeln(&decl, "def %s(", methodName)
			writeln(&decl, "    self,")
			if hasReq {
				writeln(&decl, "    req: %sReq,", funcId)
			}
			writeln(&decl, "    opts: Optional[ServiceOptions] = None,")
			writeln(&decl, ") -> %sOutput:", funcId)
		} else if hasReq {
			writeln(&decl, "def %s(", methodName)
			writeln(&decl, "    self,")
			writeln(&decl, "    req: %sReq,", funcId)
			writeln(&decl, ") -> Any:")
		} else {
			writeln(&decl, "def %s(self) -> Any:", methodName)
		}

		// Docstring
		var docLines []string
		if fn.Spec.Desc != nil {
			if name := schemamold.DescName(fn.Spec.Desc); name != "" {
				docLines = append(docLines, name)
			}
			if summary := schemamold.DescSummary(fn.Spec.Desc); summary != "" {
				docLines = append(docLines, summary)
			}
		}
		if hasReq {
			docLines = append(docLines, fmt.Sprintf(":param req: %sReq request parameters.", funcId))
		}
		if isInfer {
			docLines = append(docLines, ":param select_results: Optional custom scoring/filtering callback.")
			docLines = append(docLines, ":param opts: Optional per-request execution overrides.")
			docLines = append(docLines, fmt.Sprintf(":return: %sSelected with recommendation results.", funcId))
		} else if fn.Spec.Typ == sdkfuncz.TypClientQuery {
			docLines = append(docLines, ":param opts: Optional per-request execution overrides.")
			docLines = append(docLines, fmt.Sprintf(":return: %sOutput query result.", funcId))
		}
		if len(docLines) > 0 {
			writeln(&decl, "    \"\"\"%s\"\"\"", strings.Join(docLines, "\n\n    "))
		}

		modelName, ok := strings.CutPrefix(fn.Spec.Id, "product-")
		checker.True(ok, "id does not start with 'product-': %s", fn.Spec.Id)

		switch fn.Spec.Typ {
		case sdkfuncz.TypClientInfer:
			writeln(&decl, "    infer = self._services.infer()")
			writeln(&decl, "    req_input = req.get(\"input\") or {}")
			writeln(&decl, "    partitions = dict(req.get(\"partitions\") or {})")
			if fn.Spec.Version != "" {
				writeln(&decl, "    partitions[\"version\"] = %q", fn.Spec.Version)
			}
			writeln(&decl, "    results = infer.run({")
			writeln(&decl, "        \"dataset\": req_input.get(\"dataset\"),")
			writeln(&decl, "        \"model_name\": %q,", modelName)
			writeln(&decl, "        \"params\": req_input,")
			writeln(&decl, "        \"partitions\": partitions,")
			writeln(&decl, "        \"session_id\": req_input.get(\"session_id\"),")
			writeln(&decl, "        \"user_id\": req_input.get(\"user_id\"),")
			writeln(&decl, "        \"service_options\": opts,")
			writeln(&decl, "    })")
			writeln(&decl, "    req_partitions = req.get(\"partitions\") or {}")
			writeln(&decl, "    return infer.select_results({")
			writeln(&decl, "        \"dataset\": req_input.get(\"dataset\"),")
			writeln(&decl, "        \"model_name\": %q,", modelName)
			writeln(&decl, "        \"model_version\": %q,", fn.Spec.Version)
			writeln(&decl, "        \"store_id\": req_partitions.get(\"store_id\"),")
			writeln(&decl, "        \"req\": req,")
			writeln(&decl, "        \"results\": results,")
			writeln(&decl, "        \"rec_id\": results.get(\"rec_id\") if isinstance(results, dict) else getattr(results, \"rec_id\", None),")
			writeln(&decl, "        \"session_id\": req_input.get(\"session_id\"),")
			writeln(&decl, "        \"user_id\": req_input.get(\"user_id\"),")
			writeln(&decl, "        \"select_results\": select_results,")
			writeln(&decl, "        \"service_options\": opts,")
			writeln(&decl, "    })")
		case sdkfuncz.TypClientQuery:
			writeln(&decl, "    dash = self._services.dash()")
			writeln(&decl, "    req_input = req.get(\"input\") or {}")
			writeln(&decl, "    partitions = dict(req.get(\"partitions\") or {})")
			if fn.Spec.Version != "" {
				writeln(&decl, "    partitions[\"version\"] = %q", fn.Spec.Version)
			}
			writeln(&decl, "    query_opts: Dict[str, Any] = {")
			writeln(&decl, "        \"dataset\": req_input.get(\"dataset\"),")
			writeln(&decl, "        \"dash_name\": %q,", fn.Spec.Id)
			writeln(&decl, "        \"partitions\": partitions,")
			writeln(&decl, "        \"service_options\": opts,")
			writeln(&decl, "    }")
			if hasInputParams(fn) {
				writeln(&decl, "    query_opts[\"params\"] = req_input")
			}
			writeln(&decl, "    return dash.query(query_opts)")
		default:
			writeln(&decl, "    return {}")
		}

		writeln(&b, "%s", indent(decl.String(), "    "))
	}

	return b.String()
}
