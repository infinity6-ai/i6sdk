package jsapiwriter

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"go.code.infinity6.ai/platform/logz"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypedefwriter"
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

func Write(ctx context.Context, apis *jsparser.Apis) string {
	if apis == nil {
		return ""
	}

	var b strings.Builder

	allTdfs := &jstypez.Typedefs{}
	allTdfs.Add(apis.Typedefs.All()...)

	for _, api := range apis.Apis {
		logger.Info(ctx, "genereating js sdk api", map[string]any{"api": api.Spec.Id})
		allTdfs.Add(api.Typedefs.All()...)

		apiId := api.Spec.CamelId
		reqTypedef, hasReq := buildReqTypedef(apiId, api)
		respTypedef := buildRespTypedef(apiId, api)

		if hasReq {
			allTdfs.Add(&reqTypedef)
		}
		allTdfs.Add(&respTypedef)
	}

	b.WriteString("/**\n")
	b.WriteString(" * @typedef {import(\"../services/services.js\").Services} Services\n")
	b.WriteString(" */\n\n")
	b.WriteString(jstypedefwriter.Write(allTdfs))

	if len(apis.Apis) == 0 {
		return b.String()
	}

	b.WriteString("class Apis {\n\n")

	b.WriteString("  /**\n")
	b.WriteString("   * Constructs an Ingest client instance.\n")
	b.WriteString("   * @param {Services} services - The services manager instance.\n")
	b.WriteString("   */\n")
	b.WriteString("  constructor(services) {\n")
	b.WriteString("    /**\n")
	b.WriteString("     * @private\n")
	b.WriteString("     * @type {Services}\n")
	b.WriteString("     */\n")
	b.WriteString("    this._services = services;\n")
	b.WriteString("  }\n\n")

	for idx, api := range apis.Apis {
		if idx > 0 {
			b.WriteString("\n")
		}
		var decl strings.Builder
		writeApiDecl(&decl, api)
		b.WriteString(indent(decl.String(), "  "))
	}
	b.WriteString("}\n\n")
	b.WriteString("export { Apis };\n")

	return b.String()
}

func indent(s, prefix string) string {
	lines := strings.Split(strings.TrimSuffix(s, "\n"), "\n")
	for i, l := range lines {
		if l != "" {
			lines[i] = prefix + l
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

func buildReqTypedef(apiId string, api *jsparser.Api) (jstypez.Typedef, bool) {
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

func buildRespTypedef(apiId string, api *jsparser.Api) jstypez.Typedef {
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

func writeApiDecl(b *strings.Builder, api *jsparser.Api) {
	apiId := api.Spec.CamelId

	hasParams := api.Req.Params.Head != nil
	hasQuery := api.Req.Query.Head != nil
	hasReqHeaders := api.Req.Headers.Head != nil
	hasPayload := api.Req.Payload.Head != nil

	hasReq := hasParams || hasQuery || hasReqHeaders || hasPayload

	b.WriteString("/**\n")
	if hasReq {
		fmt.Fprintf(b, " * @param {%sReq} req\n", apiId)
	}
	fmt.Fprintf(b, " * @returns {Promise<%sResp>}\n", apiId)
	b.WriteString(" */\n")

	if hasReq {
		fmt.Fprintf(b, "async %s(req) {\n", funcName(apiId))
	} else {
		fmt.Fprintf(b, "async %s() {\n", funcName(apiId))
	}

	jsPath := strings.TrimPrefix(api.Spec.Path, "/")
	if hasParams {
		jsPath = pathParamRe.ReplaceAllString(jsPath, "${encodeURIComponent(String(req.params.$1))}")
	}

	b.WriteString("  const resp = await this._services.invoker().invoke({\n")
	fmt.Fprintf(b, "    method: %q,\n", api.Spec.Method)

	fmt.Fprintf(b, "    path: `%s`,\n", jsPath)

	if hasQuery {
		b.WriteString("    query: req.query,\n")
	}
	if hasReqHeaders {
		b.WriteString("    headers: req.headers,\n")
	}
	if hasPayload {
		b.WriteString("    json: req.payload,\n")
	}

	b.WriteString("  });\n")
	b.WriteString("  const result = {status: resp.status, ok: resp.ok};\n")
	b.WriteString("  if (resp.ok) {\n")

	if api.Resp.Payload.Head != nil {
		b.WriteString("    result.payload = await resp.json();\n")
	} else {
		b.WriteString("    result.body = await resp.json();\n")
	}

	b.WriteString("  }\n")
	b.WriteString("  if (!resp.ok) {\n")
	b.WriteString("    const text = await resp.text();\n")
	b.WriteString("    try {\n")
	b.WriteString("      result.error = JSON.parse(text);\n")
	b.WriteString("    } catch {\n")
	b.WriteString("      result.error = text;\n")
	b.WriteString("    }\n")
	b.WriteString("  }\n")
	b.WriteString("  return result;\n")
	b.WriteString("}\n")
}
