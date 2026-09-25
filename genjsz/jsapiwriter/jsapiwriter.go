package jsapiwriter

import (
	"fmt"
	"regexp"
	"strings"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/camelstring"
)

var pathParamRe = regexp.MustCompile(`\{([A-Za-z0-9_]+)\}`)

func Write(api *jsparser.Api) string {
	apiId := camelstring.KebabToCamel(api.Spec.Id)
	
	hasParams := api.Req.Properties != nil && len(api.Req.Properties.All()) > 0
	
	var b strings.Builder
	b.WriteString("/**\n")
	b.WriteString(" * @param {import(\"../services/services.js\").Services} services\n")
	if hasParams {
		fmt.Fprintf(&b, " * @param {%s} params\n", api.Req.Head.Name)
	}
	fmt.Fprintf(&b, " * @returns {Promise<%s>}\n", api.Resp.Head.Name)
	b.WriteString(" */\n")
	
	if hasParams {
		fmt.Fprintf(&b, "export async function %s(services, params) {\n", apiId)
	} else {
		fmt.Fprintf(&b, "export async function %s(services) {\n", apiId)
	}
	
	jsPath := pathParamRe.ReplaceAllString(strings.TrimPrefix(api.Spec.Path, "/"), "${encodeURIComponent(String(params.$1))}")
	
	b.WriteString("  const resp = await services.invoker().invoke({\n")
	fmt.Fprintf(&b, "    method: %q,\n", api.Spec.Method)
	
	fmt.Fprintf(&b, "    path: `%s`,\n", strings.TrimPrefix(jsPath, "/"))
	
	hasPayload := false
	if hasParams {
		for _, p := range api.Req.Properties.All() {
			if p.Name == "payload" {
				hasPayload = true
				break
			}
		}
	}
	
	if hasPayload {
		b.WriteString("    json: params.payload,\n")
	}
	
	b.WriteString("  });\n")
	b.WriteString("  const result = {status: resp.status, ok: resp.ok};\n")
	b.WriteString("  if (resp.ok) {\n")
	
	hasRespPayload := false
	if api.Resp.Properties != nil {
		for _, p := range api.Resp.Properties.All() {
			if p.Name == "payload" {
				hasRespPayload = true
				break
			}
		}
	}
	if hasRespPayload {
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
