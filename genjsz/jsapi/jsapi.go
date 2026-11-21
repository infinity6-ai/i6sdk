package jsapi

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/routez/jsapiz"
	"github.com/infinity6-ai/gox/schemaz/schemamold"
	"go.code.infinity6.ai/sdkgo/genjsz/jsschema"
)

func NewSdk() *jsapiz.Parser {
	return jsapiz.New(jsschema.NewSdk())
}

var pathParamRe = regexp.MustCompile(`\{([A-Za-z0-9_]+)\}`)

func writeln(sb *strings.Builder, msg string, args ...any) {
	_, err := fmt.Fprintf(sb, msg, args...)
	errorz.Check(err)
	_, err = sb.WriteRune('\n')
	errorz.Check(err)
}

func Write(ctx context.Context, apis []*jsapiz.Api) string {
	return WriteV2(ctx, apis)
}

func WriteV2(ctx context.Context, apis []*jsapiz.Api) string {
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

func servNameOf(api *jsapiz.Api) string {
	if api.Spec.Desc == nil || api.Spec.Desc.Labels == nil {
		return ""
	}
	return api.Spec.Desc.Labels["servname"]
}

func funcName(apiId string) string {
	name := strings.TrimPrefix(apiId, "api")
	return strings.ToLower(name[0:1]) + name[1:]
}

func writeApiDecl(b *strings.Builder, api *jsapiz.Api) {
	apiId := api.Spec.CamelId

	hasParams := api.Req.Params.Head != nil
	hasQuery := api.Req.Query.Head != nil
	hasReqHeaders := api.Req.Headers.Head != nil
	hasPayload := api.Req.Payload.Head != nil

	hasReq := hasParams || hasQuery || hasReqHeaders || hasPayload

	writeln(b, "/**")
	writeln(b, " * %s", schemamold.DescName(api.Spec.Desc))
	writeln(b, " * ")
	writeln(b, " * %s", schemamold.DescSummary(api.Spec.Desc))
	writeln(b, " * ")
	if hasReq {
		paramDesc := ""
		if api.Req.Payload.Head != nil && api.Req.Payload.Head.GoTyp != "" {
			paramDesc = fmt.Sprintf(" - [Go: %s]", api.Req.Payload.Head.GoTyp)
		}
		writeln(b, " * @param {%sReq} req%s", apiId, paramDesc)
	}
	writeln(b, " * @param {ServiceOptions} [opts]")
	returnDesc := ""
	if api.Resp.Payload.Head != nil && api.Resp.Payload.Head.GoTyp != "" {
		returnDesc = fmt.Sprintf(" - [Go: %s]", api.Resp.Payload.Head.GoTyp)
	}
	writeln(b, " * @returns {Promise<import(\"../invoker/invoker.js\").Resp<%sResp>>}%s", apiId, returnDesc)
	writeln(b, " */")

	if hasReq {
		writeln(b, "async %s(req, opts) {", funcName(apiId))
	} else {
		writeln(b, "async %s(opts) {", funcName(apiId))
	}

	jsPath := strings.TrimPrefix(api.Spec.Path, "/")
	if hasParams {
		jsPath = pathParamRe.ReplaceAllString(jsPath, "${encodeURIComponent(String(req.params.$1))}")
	}

	if api.Spec.Desc != nil && len(api.Spec.Desc.Labels) > 0 {
		var labelKeys []string
		for k := range api.Spec.Desc.Labels {
			labelKeys = append(labelKeys, k)
		}
		sort.Strings(labelKeys)

		writeln(b, "  opts = opts || {};")
		writeln(b, "  opts.labels = {")
		for _, k := range labelKeys {
			writeln(b, "    %q: %q,", k, api.Spec.Desc.Labels[k])
		}
		writeln(b, "    ...(opts.labels || {}),")
		writeln(b, "  };")
	}

	writeln(b, "  return this._services.invoker().invokeResult({")
	writeln(b, "    method: %q,", api.Spec.Method)

	if servName := servNameOf(api); servName != "" {
		// The api lives in another service: use its url from the config (servUrls), falling back to the baseUrl.
		writeln(b, "    url: `${this._services.config().servUrls?.[%q] ?? this._services.config().baseUrl}/%s`,", servName, jsPath)
	} else {
		writeln(b, "    path: `%s`,", jsPath)
	}

	if hasQuery {
		writeln(b, "    query: req.query,")
	}
	if hasReqHeaders {
		writeln(b, "    headers: req.headers,")
	}
	if hasPayload {
		writeln(b, "    json: req.payload,")
	}

	writeln(b, "    opts: opts,")

	if api.Spec.Desc != nil && api.Spec.Desc.Labels != nil && api.Spec.Desc.Labels["sinkable"] == "true" {
		writeln(b, "    sinkable: true,")
	}

	writeln(b, "  });")
	writeln(b, "}")
}
