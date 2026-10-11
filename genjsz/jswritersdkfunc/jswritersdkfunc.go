package jswritersdkfunc

import (
	"context"
	"fmt"
	"strings"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"go.code.infinity6.ai/sdkgo/camelstring"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparsersdkfunc"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
)

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

func ExtractTypedefs(sdkFuncs []*jsparsersdkfunc.SdkFunc) jstypez.Typedefs {
	tdfs := jstypez.Typedefs{}
	for _, fn := range sdkFuncs {
		funcId := fn.Spec.CamelId

		hasInput := fn.Req.Input.Head != nil
		hasPartitions := fn.Req.Partitions.Head != nil
		hasReq := hasInput || hasPartitions

		if hasReq {
			reqTypedef := jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "object", Name: funcId + "Req", Summary: funcId + " Request"},
				Properties: &jstypez.Tuples{},
			}
			if hasInput {
				reqTypedef.Properties.Add(jstypez.Tuple{Typ: fn.Req.Input.Head.Typ, GoTyp: fn.Req.Input.Head.GoTyp, Name: "input"})
			}
			if hasPartitions {
				reqTypedef.Properties.Add(jstypez.Tuple{Typ: fn.Req.Partitions.Head.Typ, GoTyp: fn.Req.Partitions.Head.GoTyp, Name: "partitions"})
			}
			tdfs.Add(&reqTypedef)
		}

		hasOutput := fn.Resp.Output.Head != nil
		respTypedef := jstypez.Typedef{
			Head:       &jstypez.Tuple{Typ: "object", Name: funcId + "Resp", Summary: funcId + " Response"},
			Properties: &jstypez.Tuples{},
		}
		if hasOutput {
			respTypedef.Properties.Add(jstypez.Tuple{Typ: fn.Resp.Output.Head.Typ, GoTyp: fn.Resp.Output.Head.GoTyp, Name: "output"})
		}
		tdfs.Add(&respTypedef)
	}
	return tdfs
}

func funcName(funcId string, version string) string {
	name := strings.TrimPrefix(funcId, "sdkFunc")
	name = strings.ToLower(name[0:1]) + name[1:]
	if version != "" {
		name += camelstring.ToCamelUpper(version)
	}
	return name
}

func Write(ctx context.Context, sdkFuncs []*jsparsersdkfunc.SdkFunc) string {
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
			writeln(&decl, " * %s", jsparser.DescName(fn.Spec.Desc))
			writeln(&decl, " * ")
			writeln(&decl, " * %s", jsparser.DescSummary(fn.Spec.Desc))
			writeln(&decl, " * ")
		}
		if hasReq {
			paramDesc := ""
			if fn.Req.Input.Head != nil && fn.Req.Input.Head.GoTyp != "" {
				paramDesc = fmt.Sprintf(" - [Go: %s]", fn.Req.Input.Head.GoTyp)
			}
			writeln(&decl, " * @param {%sReq} req%s", funcId, paramDesc)
		}
		returnDesc := ""
		if fn.Resp.Output.Head != nil && fn.Resp.Output.Head.GoTyp != "" {
			returnDesc = fmt.Sprintf(" - [Go: %s]", fn.Resp.Output.Head.GoTyp)
		}
		writeln(&decl, " * @returns {Promise<%sOutput>}%s", funcId, returnDesc)
		writeln(&decl, " */")

		if hasReq {
			writeln(&decl, "async %s(req) {", funcName(funcId, fn.Spec.Version))
		} else {
			writeln(&decl, "async %s() {", funcName(funcId, fn.Spec.Version))
		}

		modelName, ok := strings.CutPrefix(fn.Spec.Id, "product-")
		checker.True(ok, "id does not start with 'product-': %s", fn.Spec.Id)

		switch fn.Spec.Typ {
		case sdkfuncz.TypClientInfer:
			writeln(&decl, "  const infer = this._services.infer()")
			writeln(&decl, "  const ret = await infer.run({")
			writeln(&decl, "    dataset: req.input.dataset,")
			writeln(&decl, "    modelName: \"%s\",", modelName)
			writeln(&decl, "    params: req.input,")
			if fn.Spec.Version != "" {
				writeln(&decl, "    partitions: {\"version\": \"%s\", ...req.partitions},", fn.Spec.Version)
			} else {
				writeln(&decl, "    partitions: req.partitions,")
			}
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
