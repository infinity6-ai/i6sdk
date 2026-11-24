package forgepyapi

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/routez/apispecz"
	"github.com/infinity6-ai/gox/schemaz/schemamold"
	"github.com/infinity6-ai/i6sdk/forgez/forgepy/forgepyschema"
)

func NewSdk() *apispecz.Parser {
	return apispecz.New(forgepyschema.NewSdk())
}

var pathParamRe = regexp.MustCompile(`\{([A-Za-z0-9_]+)(\.\.\.)?\}`)

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

func servNameOf(api *apispecz.Api) string {
	if api.Spec.Desc == nil || api.Spec.Desc.Labels == nil {
		return ""
	}
	return api.Spec.Desc.Labels["servname"]
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

func funcName(apiId string) string {
	name := strings.TrimPrefix(apiId, "api")
	return toSnakeCase(name)
}

func camelFuncName(apiId string) string {
	name := strings.TrimPrefix(apiId, "api")
	if name == "" {
		return ""
	}
	return strings.ToLower(name[0:1]) + name[1:]
}

func Write(ctx context.Context, apis []*apispecz.Api) string {
	if len(apis) == 0 {
		return ""
	}

	var b strings.Builder

	writeln(&b, "class Apis:")
	writeln(&b, "    \"\"\"Auto-generated API client methods for i6 services.\"\"\"")
	writeln(&b, "")
	writeln(&b, "    def __init__(self, services: Any) -> None:")
	writeln(&b, "        self._services = services")

	for _, api := range apis {
		writeln(&b, "")
		var decl strings.Builder
		writeApiDecl(&decl, api)
		writeln(&b, "%s", indent(decl.String(), "    "))
	}

	return b.String()
}

func writeApiDecl(b *strings.Builder, api *apispecz.Api) {
	apiId := api.Spec.CamelId
	methodName := funcName(apiId)
	camelName := camelFuncName(apiId)

	hasParams := api.Req.Params.Head != nil
	hasQuery := api.Req.Query.Head != nil
	hasReqHeaders := api.Req.Headers.Head != nil
	hasPayload := api.Req.Payload.Head != nil
	hasReq := hasParams || hasQuery || hasReqHeaders || hasPayload

	reqType := apiId + "Req"
	respType := apiId + "Resp"

	writeln(b, "def %s(", methodName)
	writeln(b, "    self,")
	if hasReq {
		writeln(b, "    req: Optional[%s] = None,", reqType)
	}
	writeln(b, "    opts: Optional[ServiceOptions] = None,")
	writeln(b, ") -> Resp[%s]:", respType)

	// Docstring
	var docLines []string
	if name := schemamold.DescName(api.Spec.Desc); name != "" {
		docLines = append(docLines, name)
	}
	if summary := schemamold.DescSummary(api.Spec.Desc); summary != "" {
		docLines = append(docLines, summary)
	}
	if hasReq {
		docLines = append(docLines, fmt.Sprintf(":param req: %s request payload.", reqType))
	}
	docLines = append(docLines, ":param opts: Optional per-request execution overrides.")
	docLines = append(docLines, fmt.Sprintf(":return: Resp containing parsed %s.", respType))

	writeln(b, "    \"\"\"%s\"\"\"", strings.Join(docLines, "\n\n    "))

	if hasReq {
		writeln(b, "    req = req or {}")
	}

	pyPath := strings.TrimPrefix(api.Spec.Path, "/")
	if hasParams {
		// Use urllib.parse.quote for path parameters
		interpolated := pathParamRe.ReplaceAllStringFunc(pyPath, func(m string) string {
			sub := pathParamRe.FindStringSubmatch(m)
			paramName := sub[1]
			if len(sub) > 2 && sub[2] == "..." {
				return fmt.Sprintf("{quote(str((req.get('params') or {}).get('%s', '')), safe='/')}", paramName)
			}
			return fmt.Sprintf("{quote(str((req.get('params') or {}).get('%s', '')))}", paramName)
		})
		writeln(b, "    path = f%q", interpolated)
	} else {
		writeln(b, "    path = %q", pyPath)
	}

	if api.Spec.Desc != nil && len(api.Spec.Desc.Labels) > 0 {
		var labelKeys []string
		for k := range api.Spec.Desc.Labels {
			labelKeys = append(labelKeys, k)
		}
		sort.Strings(labelKeys)

		writeln(b, "    effective_opts: Dict[str, Any] = dict(opts or {})")
		writeln(b, "    labels = {")
		for _, k := range labelKeys {
			writeln(b, "        %q: %q,", k, api.Spec.Desc.Labels[k])
		}
		writeln(b, "    }")
		writeln(b, "    labels.update(effective_opts.get(\"labels\") or {})")
		writeln(b, "    effective_opts[\"labels\"] = labels")
	}

	writeln(b, "    invoke_kwargs: Dict[str, Any] = {")
	writeln(b, "        \"method\": %q,", api.Spec.Method)

	servName := servNameOf(api)
	if servName == "" {
		servName = "decsuite"
	}
	writeln(b, "        \"url\": f\"{((self._services.config().serv_urls or {}).get('%s') or '')}/{path}\",", servName)

	if hasQuery {
		writeln(b, "        \"query\": req.get(\"query\"),")
	}
	if hasReqHeaders {
		writeln(b, "        \"headers\": req.get(\"headers\"),")
	}
	if hasPayload {
		writeln(b, "        \"json\": req.get(\"payload\"),")
	}

	if api.Spec.Desc != nil && len(api.Spec.Desc.Labels) > 0 {
		writeln(b, "        \"opts\": effective_opts,")
	} else {
		writeln(b, "        \"opts\": opts,")
	}

	if api.Spec.Desc != nil && api.Spec.Desc.Labels != nil && api.Spec.Desc.Labels["sinkable"] == "true" {
		writeln(b, "        \"sinkable\": True,")
	}

	writeln(b, "    }")
	writeln(b, "    return self._services.invoker().invoke_result(**invoke_kwargs)")

	// Provide camelCase alias if different from snake_case
	if camelName != "" && camelName != methodName {
		writeln(b, "")
		writeln(b, "%s = %s", camelName, methodName)
	}
}
