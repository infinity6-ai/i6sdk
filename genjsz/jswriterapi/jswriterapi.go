package jswriterapi

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"go.code.infinity6.ai/platform/logz"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparserapi"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

var pathParamRe = regexp.MustCompile(`\{([A-Za-z0-9_]+)\}`)

func writeln(sb *strings.Builder, msg string, args ...any) {
	_, err := fmt.Fprintf(sb, msg, args...)
	errorz.Check(err)
	_, err = sb.WriteRune('\n')
	errorz.Check(err)
}

func ExtractTypedefs(apis []*jsparserapi.Api) jstypez.Typedefs {
	tdfs := jstypez.Typedefs{}
	for _, api := range apis {
		apiId := api.Spec.CamelId
		reqTypedef, hasReq := buildReqTypedef(apiId, api)
		respTypedef := buildRespTypedef(apiId, api)

		if hasReq {
			tdfs.Add(&reqTypedef)
		}
		tdfs.Add(&respTypedef)
	}
	return tdfs
}

func WriteV2(ctx context.Context, apis []*jsparserapi.Api) string {
	if len(apis) == 0 {
		return ""
	}

	var b strings.Builder

	writeln(&b, "class Apis {")
	writeln(&b, "")
	writeln(&b, "  /**")
	writeln(&b, "   * Constructs an Ingest client instance.")
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

	for idx, api := range apis {
		if idx > 0 {
			writeln(&b, "")
		}
		var decl strings.Builder
		writeApiDecl(&decl, api)
		writeln(&b, "%s", indent(decl.String(), "  "))
	}
	writeln(&b, "}")
	writeln(&b, "")
	writeln(&b, "export { Apis };")

	return b.String()
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

func buildReqTypedef(apiId string, api *jsparserapi.Api) (jstypez.Typedef, bool) {
	hasParams := api.Req.Params.Head != nil
	hasQuery := api.Req.Query.Head != nil
	hasReqHeaders := api.Req.Headers.Head != nil
	hasPayload := api.Req.Payload.Head != nil

	hasReq := hasParams || hasQuery || hasReqHeaders || hasPayload

	reqTypedef := jstypez.Typedef{
		Head:       &jstypez.Tuple{Typ: "object", Name: apiId + "Req", Summary: apiId + " Request"},
		Properties: &jstypez.Tuples{},
	}
	if hasParams {
		reqTypedef.Properties.Add(jstypez.Tuple{Typ: api.Req.Params.Head.Typ, Name: "params"})
	}
	if hasQuery {
		reqTypedef.Properties.Add(jstypez.Tuple{Typ: api.Req.Query.Head.Typ, Name: "query"})
	}
	if hasReqHeaders {
		reqTypedef.Properties.Add(jstypez.Tuple{Typ: api.Req.Headers.Head.Typ, Name: "headers"})
	}
	if hasPayload {
		reqTypedef.Properties.Add(jstypez.Tuple{Typ: api.Req.Payload.Head.Typ, Name: "payload"})
	}
	return reqTypedef, hasReq
}

func buildRespTypedef(apiId string, api *jsparserapi.Api) jstypez.Typedef {
	respTypedef := jstypez.Typedef{
		Head:       &jstypez.Tuple{Typ: "object", Name: apiId + "Resp", Summary: apiId + " Response"},
		Properties: &jstypez.Tuples{},
	}
	if api.Resp.Headers.Head != nil {
		respTypedef.Properties.Add(jstypez.Tuple{Typ: api.Resp.Headers.Head.Typ, Name: "headers"})
	}
	if api.Resp.Payload.Head != nil {
		respTypedef.Properties.Add(jstypez.Tuple{Typ: api.Resp.Payload.Head.Typ, Name: "payload"})
	}
	return respTypedef
}

func funcName(apiId string) string {
	name := strings.TrimPrefix(apiId, "api")
	return strings.ToLower(name[0:1]) + name[1:]
}

func writeApiDecl(b *strings.Builder, api *jsparserapi.Api) {
	apiId := api.Spec.CamelId

	hasParams := api.Req.Params.Head != nil
	hasQuery := api.Req.Query.Head != nil
	hasReqHeaders := api.Req.Headers.Head != nil
	hasPayload := api.Req.Payload.Head != nil

	hasReq := hasParams || hasQuery || hasReqHeaders || hasPayload

	writeln(b, "/**")
	writeln(b, " * %s", jsparser.DescName(api.Spec.Desc))
	writeln(b, " * ")
	writeln(b, " * %s", jsparser.DescSummary(api.Spec.Desc))
	writeln(b, " * ")
	if hasReq {
		writeln(b, " * @param {%sReq} req", apiId)
	}
	writeln(b, " * @returns {Promise<%sResp>}", apiId)
	writeln(b, " */")

	if hasReq {
		writeln(b, "async %s(req) {", funcName(apiId))
	} else {
		writeln(b, "async %s() {", funcName(apiId))
	}

	jsPath := strings.TrimPrefix(api.Spec.Path, "/")
	if hasParams {
		jsPath = pathParamRe.ReplaceAllString(jsPath, "${encodeURIComponent(String(req.params.$1))}")
	}

	writeln(b, "  const resp = await this._services.invoker().invoke({")
	writeln(b, "    method: %q,", api.Spec.Method)

	writeln(b, "    path: `%s`,", jsPath)

	if hasQuery {
		writeln(b, "    query: req.query,")
	}
	if hasReqHeaders {
		writeln(b, "    headers: req.headers,")
	}
	if hasPayload {
		writeln(b, "    json: req.payload,")
	}

	writeln(b, "  });")
	writeln(b, "  const result = {status: resp.status, ok: resp.ok};")
	writeln(b, "  if (resp.ok) {")

	if api.Resp.Payload.Head != nil {
		writeln(b, "    result.payload = await resp.json();")
	} else {
		writeln(b, "    result.body = await resp.json();")
	}

	writeln(b, "  }")
	writeln(b, "  if (!resp.ok) {")
	writeln(b, "    const text = await resp.text();")
	writeln(b, "    try {")
	writeln(b, "      result.error = JSON.parse(text);")
	writeln(b, "    } catch {")
	writeln(b, "      result.error = text;")
	writeln(b, "    }")
	writeln(b, "  }")
	writeln(b, "  return result;")
	writeln(b, "}")
}
