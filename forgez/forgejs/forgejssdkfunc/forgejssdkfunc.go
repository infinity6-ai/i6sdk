package forgejssdkfunc

import (
	"context"
	"fmt"
	"strings"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/commonz/stringz/camelz"
	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"github.com/infinity6-ai/gox/schemaz/schemamold"
	"github.com/infinity6-ai/gox/schemaz/typeastz"
	"go.code.infinity6.ai/sdkgo/forgez/forgejs/forgejsschema"
	"go.code.infinity6.ai/sdkgo/mold/moldsdkfunc"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
)

func NewSdk() *moldsdkfunc.Parser {
	return moldsdkfunc.New(forgejsschema.NewSdk())
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

func funcName(funcId string, version string) string {
	name := strings.TrimPrefix(funcId, "sdkFunc")
	name = strings.ToLower(name[0:1]) + name[1:]
	if version != "" {
		name += camelz.P(version).P()
	}
	return name
}

func Write(ctx context.Context, sdkFuncs []*moldsdkfunc.SdkFunc) string {
	if len(sdkFuncs) == 0 {
		return ""
	}

	var b strings.Builder

	writeln(&b, "class SdkFuncs {")
	writeln(&b, "")
	writeln(&b, "  /**")
	writeln(&b, "   * Constructs an SdkFuncs instance.")
	writeln(&b, "   * @param {Services} services - The services manager instance.")
	writeln(&b, "   */")
	writeln(&b, "  constructor(services) {")
	writeln(&b, "    /**")
	writeln(&b, "     * @private")
	writeln(&b, "     * @type {Services}")
	writeln(&b, "     */")
	writeln(&b, "    this._services = services;")
	writeln(&b, "  }")
	writeln(&b, "")

	for idx, fn := range sdkFuncs {
		if idx > 0 {
			writeln(&b, "")
		}
		var decl strings.Builder
		funcId := fn.Spec.CamelId

		hasInput := fn.Req.Input.Head != nil
		hasPartitions := fn.Req.Partitions.Head != nil
		hasReq := hasInput || hasPartitions

		writeln(&decl, "/**")
		if fn.Spec.Desc != nil {
			writeln(&decl, " * %s", schemamold.DescName(fn.Spec.Desc))
			writeln(&decl, " * ")
			writeln(&decl, " * %s", schemamold.DescSummary(fn.Spec.Desc))
			writeln(&decl, " * ")
		}
		if hasReq {
			paramDesc := ""
			if fn.Req.Input.Head != nil && fn.Req.Input.Head.GoTyp != "" {
				paramDesc = fmt.Sprintf(" - [Go: %s]", fn.Req.Input.Head.GoTyp)
			}
			writeln(&decl, " * @param {%sReq} req%s", funcId, paramDesc)
		}
		isInfer := fn.Spec.Typ == sdkfuncz.TypClientInfer
		if isInfer {
			writeln(&decl, " * @param {(req: %sReq, results: %sOutput) => Promise<%sSelected>} [selectResults] - picks which of the recommended skus are shown and may attach your own data to each in enriched; default: for each result, in order, the first sku not in req.input.skus (the cart) and not already picked, with the rating of that result; no duplicates, and results with no such sku are skipped", funcId, funcId, funcId)
			writeln(&decl, " * @param {ServiceOptions} [opts]")
			writeln(&decl, " * @returns {Promise<%sSelected>}", funcId)
		} else {
			returnDesc := ""
			if fn.Resp.Output.Head != nil && fn.Resp.Output.Head.GoTyp != "" {
				returnDesc = fmt.Sprintf(" - [Go: %s]", fn.Resp.Output.Head.GoTyp)
			}
			if fn.Spec.Typ == sdkfuncz.TypClientQuery {
				writeln(&decl, " * @param {ServiceOptions} [opts]")
			}
			writeln(&decl, " * @returns {Promise<%sOutput>}%s", funcId, returnDesc)
		}
		writeln(&decl, " */")

		if isInfer {
			writeln(&decl, "async %s(req, selectResults, opts) {", funcName(funcId, fn.Spec.Version))
		} else if fn.Spec.Typ == sdkfuncz.TypClientQuery {
			writeln(&decl, "async %s(req, opts) {", funcName(funcId, fn.Spec.Version))
		} else if hasReq {
			writeln(&decl, "async %s(req) {", funcName(funcId, fn.Spec.Version))
		} else {
			writeln(&decl, "async %s() {", funcName(funcId, fn.Spec.Version))
		}

		modelName, ok := strings.CutPrefix(fn.Spec.Id, "product-")
		checker.True(ok, "id does not start with 'product-': %s", fn.Spec.Id)

		switch fn.Spec.Typ {
		case sdkfuncz.TypClientInfer:
			writeln(&decl, "  const infer = this._services.infer()")
			writeln(&decl, "  const results = await infer.run({")
			writeln(&decl, "    dataset: req.input.dataset,")
			writeln(&decl, "    modelName: \"%s\",", modelName)
			writeln(&decl, "    params: req.input,")
			if fn.Spec.Version != "" {
				writeln(&decl, "    partitions: {\"version\": \"%s\", ...req.partitions},", fn.Spec.Version)
			} else {
				writeln(&decl, "    partitions: req.partitions,")
			}
			writeln(&decl, "    sessionId: req.input.session_id,")
			writeln(&decl, "    userId: req.input.user_id,")
			writeln(&decl, "    serviceOptions: opts,")
			writeln(&decl, "  })")
			writeln(&decl, "  return infer.selectResults({")
			writeln(&decl, "    dataset: req.input.dataset,")
			writeln(&decl, "    modelName: \"%s\",", modelName)
			writeln(&decl, "    modelVersion: \"%s\",", fn.Spec.Version)
			writeln(&decl, "    storeId: req.partitions.store_id,")
			writeln(&decl, "    req,")
			writeln(&decl, "    results,")
			writeln(&decl, "    recId: results.rec_id,")
			writeln(&decl, "    sessionId: req.input.session_id,")
			writeln(&decl, "    userId: req.input.user_id,")
			writeln(&decl, "    selectResults,")
			writeln(&decl, "    serviceOptions: opts,")
			writeln(&decl, "  })")
		case sdkfuncz.TypClientQuery:
			writeln(&decl, "  const dash = this._services.dash()")
			writeln(&decl, "  const ret = await dash.query({")
			writeln(&decl, "    dataset: req.input.dataset,")
			writeln(&decl, "    dashName: \"%s\",", fn.Spec.Id)
			if hasInputParams(fn) {
				writeln(&decl, "    params: req.input,")
			}
			if fn.Spec.Version != "" {
				writeln(&decl, "    partitions: {\"version\": \"%s\", ...req.partitions},", fn.Spec.Version)
			} else {
				writeln(&decl, "    partitions: req.partitions,")
			}
			writeln(&decl, "    serviceOptions: opts,")
			writeln(&decl, "  })")
			writeln(&decl, "  return ret;")
		default:
			if hasReq {
				writeln(&decl, "  console.log(req);")
			} else {
				writeln(&decl, "  console.log(\"no req\");")
			}
			writeln(&decl, "  return {};")
		}
		writeln(&decl, "}")

		writeln(&b, "%s", indent(decl.String(), "  "))
	}
	writeln(&b, "}")
	writeln(&b, "")
	writeln(&b, "export { SdkFuncs };")

	return b.String()
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
