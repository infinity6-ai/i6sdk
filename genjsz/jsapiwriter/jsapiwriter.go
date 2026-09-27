package jsapiwriter

import (
	"fmt"
	"regexp"
	"strings"

	"go.code.infinity6.ai/sdkgo/camelstring"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypedefwriter"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

var pathParamRe = regexp.MustCompile(`\{([A-Za-z0-9_]+)\}`)

func Write(api *jsparser.Api) string {
	apiId := camelstring.KebabToCamel(api.Spec.Id)

	hasParams := api.Req.Params.Head != nil
	hasQuery := api.Req.Query.Head != nil
	hasReqHeaders := api.Req.Headers.Head != nil
	hasPayload := api.Req.Payload.Head != nil

	hasReq := hasParams || hasQuery || hasReqHeaders || hasPayload

	var b strings.Builder

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

	allTdfs := &jstypez.Typedefs{}
	allTdfs.Add(api.Typedefs.All()...)
	if hasReq {
		allTdfs.Add(&reqTypedef)
	}
	allTdfs.Add(&respTypedef)

	b.WriteString(jstypedefwriter.Write(allTdfs))

	b.WriteString("/**\n")
	b.WriteString(" * @param {import(\"../services/services.js\").Services} services\n")
	if hasReq {
		fmt.Fprintf(&b, " * @param {%sReq} req\n", apiId)
	}
	fmt.Fprintf(&b, " * @returns {Promise<%sResp>}\n", apiId)
	b.WriteString(" */\n")

	if hasReq {
		fmt.Fprintf(&b, "export async function %s(services, req) {\n", apiId)
	} else {
		fmt.Fprintf(&b, "export async function %s(services) {\n", apiId)
	}

	jsPath := strings.TrimPrefix(api.Spec.Path, "/")
	if hasParams {
		jsPath = pathParamRe.ReplaceAllString(jsPath, "${encodeURIComponent(String(req.params.$1))}")
	}

	b.WriteString("  const resp = await services.invoker().invoke({\n")
	fmt.Fprintf(&b, "    method: %q,\n", api.Spec.Method)

	fmt.Fprintf(&b, "    path: `%s`,\n", jsPath)

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
	return b.String()
}
