package cmdsdkweb

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	texttemplate "text/template"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/commonz/filez"
	"github.com/infinity6-ai/gox/commonz/stringz/camelz"
	"github.com/infinity6-ai/gox/routez/apispecz"
	"github.com/infinity6-ai/gox/schemaz/schemamold"
	"github.com/infinity6-ai/gox/schemaz/typeastz"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/forgez/forgejs/forgejsschema"
	"github.com/infinity6-ai/i6sdk/mold/moldsdkfunc"
	"github.com/spf13/cobra"
)

//go:embed docstpl/*.tmpl
var docsTemplatesFS embed.FS

// SdkJson represents the data extracted from handwritten modules by docgen.mjs
type SdkJson struct {
	Modules  []SdkModule           `json:"modules"`
	Typedefs map[string]SdkTypedef `json:"typedefs"`
}

type SdkModule struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	CallPath    string      `json:"callPath"`
	Summary     string      `json:"summary"`
	Description string      `json:"description"`
	Methods     []SdkMethod `json:"methods"`
}

type SdkMethod struct {
	Name        string     `json:"name"`
	CallPath    string     `json:"callPath"`
	Summary     string     `json:"summary"`
	Description string     `json:"description"`
	Params      []SdkParam `json:"params"`
	ReturnType  string     `json:"returnType"`
	ReturnDoc   string     `json:"returnDoc"`
}

type SdkParam struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Optional bool   `json:"optional"`
	Doc      string `json:"doc"`
}

type SdkTypedef struct {
	Name        string            `json:"name"`
	Summary     string            `json:"summary"`
	Description string            `json:"description"`
	Properties  []SdkTypeProperty `json:"properties"`
}

type SdkTypeProperty struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Optional bool   `json:"optional"`
	Doc      string `json:"doc"`
}

// NavItem represents an entry in the documentation sidebar navigation
type NavItem struct {
	Id          string
	Title       string
	Url         string
	Method      string
	MethodLower string
	Summary     string
	SubItems    []NavItem
}

// FieldDoc represents a flattened request or response field row in documentation tables
type FieldDoc struct {
	Name        string
	Type        string
	Optional    bool
	SummaryHtml template.HTML
	Anchor      string
}

// BasePageData contains common navigation and layout fields for all HTML pages
type BasePageData struct {
	Title            string
	Description      string
	MethodId         string
	HasTryIt         bool
	BaseUrl          string
	BaseRelativePath string
	CurrentSection   string
	CurrentMethodId  string
	NavModules       []NavItem
	NavApis          []NavItem
	NavFunctions     []NavItem
}

// MethodPageData is passed to method.html.tmpl
type MethodPageData struct {
	BasePageData
	CallPath        string
	Name            string
	SummaryHtml     template.HTML
	DescriptionHtml template.HTML
	HttpMethod      string
	HttpMethodLower string
	HttpPath        string
	TokenType       string
	RequestFields   []FieldDoc
	ResponseFields  []FieldDoc
	ExampleCode     string
}

// SectionPageData is passed to section.html.tmpl
type SectionPageData struct {
	BasePageData
	CallPath        string
	SummaryHtml     template.HTML
	DescriptionHtml template.HTML
	Methods         []MethodSummary
	Typedefs        []SdkTypedef
}

type MethodSummary struct {
	CallPath        string
	Url             string
	Summary         string
	HttpMethod      string
	HttpMethodLower string
}

// IndexPageData is passed to index.html.tmpl
type IndexPageData struct {
	BasePageData
}

// MethodCatalogItem is written to methods.json for inspect(x)
type MethodCatalogItem struct {
	CallPath       string     `json:"callPath"`
	Name           string     `json:"name"`
	Summary        string     `json:"summary"`
	HttpMethod     string     `json:"httpMethod,omitempty"`
	HttpPath       string     `json:"httpPath,omitempty"`
	TokenType      string     `json:"tokenType,omitempty"`
	Url            string     `json:"url"`
	RequestFields  []FieldDoc `json:"requestFields,omitempty"`
	ResponseFields []FieldDoc `json:"responseFields,omitempty"`
}

type DocsOptions struct {
	SdkJsonPath  string
	OutDir       string
	BaseUrl      string
	AllowMissing bool
}

func prepareDocsCmd() *cobra.Command {
	var sdkPath string
	var outDir string
	var baseUrl string
	var allowMissing bool

	cmd := &cobra.Command{
		Use:   "docs",
		Short: "generate static HTML docs site for sdkweb",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return RunDocs(ctx, DocsOptions{
				SdkJsonPath:  sdkPath,
				OutDir:       outDir,
				BaseUrl:      baseUrl,
				AllowMissing: allowMissing,
			})
		},
	}

	cmd.Flags().StringVar(&sdkPath, "sdk", "target/sdk.json", "path to sdk.json containing handwritten modules")
	cmd.Flags().StringVar(&outDir, "out", "target/docs", "output directory for static HTML pages")
	cmd.Flags().StringVar(&baseUrl, "base-url", "https://cdn.infinity6.ai/docs", "canonical base URL for documentation links")
	cmd.Flags().BoolVar(&allowMissing, "allow-missing", false, "continue generating even if missing docs are found")

	return cmd
}

func RunDocs(ctx context.Context, opts DocsOptions) error {
	// 1. Parse APIs and SdkFuncs in-process
	schemaParser := forgejsschema.NewSdk()
	apiParser := apispecz.New(schemaParser)
	sdkFuncParser := moldsdkfunc.New(schemaParser)

	apis := apicollectionz.Apis()
	parsedApis := apiParser.ParseApis(ctx, apis)

	funcs := apicollectionz.Funcs()
	parsedSdkFuncs := sdkFuncParser.ParseSdkFuncs(ctx, funcs)

	// 2. Read sdk.json (if exists, or empty)
	var sdkData SdkJson
	if opts.SdkJsonPath != "" && filez.FileExists(opts.SdkJsonPath) {
		bytes, err := os.ReadFile(opts.SdkJsonPath)
		if err != nil {
			return fmt.Errorf("failed to read sdk.json: %w", err)
		}
		if err := json.Unmarshal(bytes, &sdkData); err != nil {
			return fmt.Errorf("failed to parse sdk.json: %w", err)
		}
	}

	// 3. Check docs: missing summaries fail build
	var missingDocs []string
	checkFieldDoc := func(apiId, fieldName, summary string) {
		if strings.TrimSpace(summary) == "" {
			missingDocs = append(missingDocs, fmt.Sprintf("%s: %s: missing doc", apiId, fieldName))
		}
	}

	for _, api := range parsedApis.Apis {
		if strings.TrimSpace(schemamold.DescSummary(api.Spec.Desc)) == "" {
			missingDocs = append(missingDocs, fmt.Sprintf("%s: missing doc", api.Spec.Id))
		}
		CheckApiTypedefFields(api, checkFieldDoc)
	}

	for _, fn := range parsedSdkFuncs.SdkFuncs {
		if strings.TrimSpace(schemamold.DescSummary(fn.Spec.Desc)) == "" {
			missingDocs = append(missingDocs, fmt.Sprintf("%s: missing doc", fn.Spec.Id))
		}
		CheckSdkFuncTypedefFields(fn, checkFieldDoc)
	}

	if len(missingDocs) > 0 {
		for _, m := range missingDocs {
			fmt.Fprintln(os.Stderr, m)
		}
		if !opts.AllowMissing {
			return fmt.Errorf("doc check failed with %d missing doc(s)", len(missingDocs))
		}
	}

	// 4. Generate static HTML pages
	generator := newSiteGenerator(opts, parsedApis.Apis, parsedSdkFuncs.SdkFuncs, sdkData)
	return generator.Generate(ctx)
}

func findTypedef(tdfs typeastz.Typedefs, name string) *typeastz.Typedef {
	for _, td := range tdfs.All() {
		if td != nil && td.Head != nil && td.Head.Name == name {
			return td
		}
	}
	return nil
}

// fieldGroup is one top-level part of a request or response (params,
// payload, ...) and the typedef that describes it.
type fieldGroup struct {
	prefix string
	td     typeastz.Typedef
}

func (fg fieldGroup) typeName() string {
	if fg.td.Head == nil {
		return ""
	}
	return fg.td.Head.Typ
}

func CheckApiTypedefFields(api *apispecz.Api, report func(id, field, summary string)) {
	checkTypedefFields(api.Spec.Id, api.Typedefs, []fieldGroup{
		{"params", api.Req.Params},
		{"query", api.Req.Query},
		{"headers", api.Req.Headers},
		{"payload", api.Req.Payload},
		{"respPayload", api.Resp.Payload},
	}, report)
}

func CheckSdkFuncTypedefFields(fn *moldsdkfunc.SdkFunc, report func(id, field, summary string)) {
	checkTypedefFields(fn.Spec.Id, fn.Typedefs, []fieldGroup{
		{"input", fn.Req.Input},
		{"partitions", fn.Req.Partitions},
		{"output", fn.Resp.Output},
	}, report)
}

// checkTypedefFields reports every leaf field of the groups. A field whose
// type (or array element type) is a record with properties is not a leaf;
// its properties are reported instead.
func checkTypedefFields(id string, tdfs typeastz.Typedefs, groups []fieldGroup, report func(id, field, summary string)) {
	for _, fg := range groups {
		td := findTypedef(tdfs, fg.typeName())
		if td == nil || td.Properties == nil {
			continue
		}
		for _, p := range td.Properties.All() {
			fullName := fg.prefix + "." + p.Name
			elemType := p.Typ
			if strings.HasPrefix(elemType, "Array<") && strings.HasSuffix(elemType, ">") {
				elemType = strings.TrimPrefix(strings.TrimSuffix(elemType, ">"), "Array<")
			}
			if nested := findTypedef(tdfs, elemType); nested != nil && nested.Properties != nil && len(nested.Properties.All()) > 0 {
				for _, np := range nested.Properties.All() {
					report(id, fullName+"."+np.Name, np.Summary)
				}
				continue
			}
			report(id, fullName, p.Summary)
		}
	}
}

type siteGenerator struct {
	opts     DocsOptions
	apis     []*apispecz.Api
	sdkFuncs []*moldsdkfunc.SdkFunc
	sdkData  SdkJson

	navModules   []NavItem
	navApis      []NavItem
	navFunctions []NavItem

	methodCatalog map[string]MethodCatalogItem
	sitemapUrls   []string
}

func newSiteGenerator(opts DocsOptions, apis []*apispecz.Api, funcs []*moldsdkfunc.SdkFunc, sdk SdkJson) *siteGenerator {
	g := &siteGenerator{
		opts:          opts,
		apis:          apis,
		sdkFuncs:      funcs,
		sdkData:       sdk,
		methodCatalog: make(map[string]MethodCatalogItem),
	}
	g.buildNavigation()
	return g
}

func (g *siteGenerator) buildNavigation() {
	// 1. Handwritten Modules
	for _, mod := range g.sdkData.Modules {
		item := NavItem{
			Id:      mod.Id,
			Title:   mod.Name,
			Url:     mod.Id + ".html",
			Summary: mod.Summary,
		}
		for _, m := range mod.Methods {
			methodSlug := mod.Id + "-" + kebabCase(m.Name) + ".html"
			item.SubItems = append(item.SubItems, NavItem{
				Id:      m.Name,
				Title:   m.Name,
				Url:     methodSlug,
				Summary: m.Summary,
			})
		}
		g.navModules = append(g.navModules, item)
	}

	// Types helper section
	g.navModules = append(g.navModules, NavItem{
		Id:      "types",
		Title:   "Types & Helpers",
		Url:     "types.html",
		Summary: "Types, constants (Grants, bitmapz), and response envelopes (Resp).",
	})

	// 2. APIs
	for _, api := range g.apis {
		method := api.Spec.Method
		fnName := apiFuncName(api.Spec.CamelId)
		slug := "apis/" + api.Spec.Id + ".html"
		g.navApis = append(g.navApis, NavItem{
			Id:          api.Spec.Id,
			Title:       fnName,
			Url:         slug,
			Method:      method,
			MethodLower: strings.ToLower(method),
			Summary:     schemamold.DescSummary(api.Spec.Desc),
		})
	}

	// 3. Functions (sdk.funcs())
	for _, fn := range g.sdkFuncs {
		fnName := sdkFuncName(fn.Spec.CamelId, fn.Spec.Version)
		slug := "functions/" + fn.Spec.Id + ".html"
		g.navFunctions = append(g.navFunctions, NavItem{
			Id:      fn.Spec.Id,
			Title:   fnName,
			Url:     slug,
			Summary: schemamold.DescSummary(fn.Spec.Desc),
		})
	}
}

func (g *siteGenerator) linkResolver(target string) string {
	parts := strings.Split(target, "#")
	if len(parts) == 2 {
		return strings.ToLower(parts[0]) + "-" + kebabCase(parts[1]) + ".html"
	}
	return strings.ToLower(target) + ".html"
}

func (g *siteGenerator) Generate(ctx context.Context) error {
	errorz.Check(os.MkdirAll(g.opts.OutDir, os.ModePerm))
	errorz.Check(os.MkdirAll(filepath.Join(g.opts.OutDir, "apis"), os.ModePerm))
	if len(g.sdkFuncs) > 0 {
		errorz.Check(os.MkdirAll(filepath.Join(g.opts.OutDir, "functions"), os.ModePerm))
	}

	indexTmpl, err := template.ParseFS(docsTemplatesFS, "docstpl/layout.html.tmpl", "docstpl/index.html.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse index template: %w", err)
	}

	sectionTmpl, err := template.ParseFS(docsTemplatesFS, "docstpl/layout.html.tmpl", "docstpl/section.html.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse section template: %w", err)
	}

	methodTmpl, err := template.ParseFS(docsTemplatesFS, "docstpl/layout.html.tmpl", "docstpl/method.html.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse method template: %w", err)
	}

	sitemapTmpl, err := texttemplate.ParseFS(docsTemplatesFS, "docstpl/sitemap.xml.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse sitemap template: %w", err)
	}

	callbackTmpl, err := template.ParseFS(docsTemplatesFS, "docstpl/callback.html.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse callback template: %w", err)
	}

	notFoundTmpl, err := template.ParseFS(docsTemplatesFS, "docstpl/404.html.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse 404 template: %w", err)
	}

	// Generate Index
	if err := g.renderIndex(indexTmpl); err != nil {
		return err
	}

	// Generate Section Pages
	if err := g.renderSections(sectionTmpl); err != nil {
		return err
	}

	// Generate API Method Pages
	if err := g.renderApiMethods(methodTmpl); err != nil {
		return err
	}

	// Generate Function Method Pages
	if err := g.renderFunctionMethods(methodTmpl); err != nil {
		return err
	}

	// Generate Handwritten Method Pages
	if err := g.renderHandwrittenMethods(methodTmpl); err != nil {
		return err
	}

	// Write methods.json
	catalogBytes, err := json.MarshalIndent(g.methodCatalog, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal methods.json: %w", err)
	}
	if err := os.WriteFile(filepath.Join(g.opts.OutDir, "methods.json"), catalogBytes, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write methods.json: %w", err)
	}

	// Write sitemap.xml
	if err := g.renderSitemap(sitemapTmpl); err != nil {
		return err
	}

	// Write callback.html and 404.html
	if err := g.renderStaticHtml(callbackTmpl, "callback.html.tmpl", "callback.html"); err != nil {
		return err
	}
	if err := g.renderStaticHtml(notFoundTmpl, "404.html.tmpl", "404.html"); err != nil {
		return err
	}

	// Copy site assets (docs.css, bar.js, inspect.js, tryit.js, search.js) if present
	g.copyAssets()

	return writeTypesJson(filepath.Join("dist", "types"), filepath.Join(g.opts.OutDir, "types.json"))
}

// writeTypesJson writes every .d.ts under typesDir as {relative path: content},
// for the editor autocomplete. Nothing is written if typesDir is missing.
func writeTypesJson(typesDir, outPath string) error {
	if !dirExists(typesDir) {
		return nil
	}
	files := map[string]string{}
	err := filepath.WalkDir(typesDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".d.ts") {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(typesDir, path)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(rel)] = string(content)
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", typesDir, err)
	}
	b, err := json.Marshal(files)
	if err != nil {
		return fmt.Errorf("failed to marshal types.json: %w", err)
	}
	if err := os.WriteFile(outPath, b, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write types.json: %w", err)
	}
	return nil
}

func (g *siteGenerator) renderIndex(tmpl *template.Template) error {
	outPath := filepath.Join(g.opts.OutDir, "index.html")
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create index.html: %w", err)
	}
	defer f.Close()

	data := IndexPageData{
		BasePageData: BasePageData{
			Title:            "i6 Web SDK Reference",
			Description:      "Interactive API reference and documentation for the i6 Web SDK.",
			BaseUrl:          g.opts.BaseUrl,
			BaseRelativePath: "",
			CurrentSection:   "",
			NavModules:       g.navModules,
			NavApis:          g.navApis,
			NavFunctions:     g.navFunctions,
		},
	}

	g.recordSitemapUrl("index.html")
	return tmpl.ExecuteTemplate(f, "layout", data)
}

func (g *siteGenerator) renderSections(tmpl *template.Template) error {
	// 1. APIs overview section
	var apiMethods []MethodSummary
	for _, api := range g.apis {
		apiMethods = append(apiMethods, MethodSummary{
			CallPath:        "sdk.apis()." + apiFuncName(api.Spec.CamelId),
			Url:             "apis/" + api.Spec.Id + ".html",
			Summary:         schemamold.DescSummary(api.Spec.Desc),
			HttpMethod:      api.Spec.Method,
			HttpMethodLower: strings.ToLower(api.Spec.Method),
		})
	}
	if err := g.renderSectionPage(tmpl, "apis.html", SectionPageData{
		BasePageData: BasePageData{
			Title:            "API Reference",
			Description:      "Direct JSON APIs for Infinity6 platform services.",
			BaseUrl:          g.opts.BaseUrl,
			BaseRelativePath: "",
			CurrentSection:   "apis",
			NavModules:       g.navModules,
			NavApis:          g.navApis,
			NavFunctions:     g.navFunctions,
		},
		CallPath:    "sdk.apis()",
		SummaryHtml: template.HTML("<p>Direct JSON APIs for Infinity6 platform services.</p>"),
		Methods:     apiMethods,
	}); err != nil {
		return err
	}

	// 2. Functions overview section
	if len(g.sdkFuncs) > 0 {
		var functionMethods []MethodSummary
		for _, fn := range g.sdkFuncs {
			functionMethods = append(functionMethods, MethodSummary{
				CallPath: "sdk.funcs()." + sdkFuncName(fn.Spec.CamelId, fn.Spec.Version),
				Url:      "functions/" + fn.Spec.Id + ".html",
				Summary:  schemamold.DescSummary(fn.Spec.Desc),
			})
		}
		if err := g.renderSectionPage(tmpl, "functions.html", SectionPageData{
			BasePageData: BasePageData{
				Title:            "Functions Reference",
				Description:      "Client-side product routines and domain queries.",
				BaseUrl:          g.opts.BaseUrl,
				BaseRelativePath: "",
				CurrentSection:   "functions",
				NavModules:       g.navModules,
				NavApis:          g.navApis,
				NavFunctions:     g.navFunctions,
			},
			CallPath:    "sdk.funcs()",
			SummaryHtml: template.HTML("<p>Client-side product routines and domain queries.</p>"),
			Methods:     functionMethods,
		}); err != nil {
			return err
		}
	}

	// 3. Handwritten module section pages
	for _, mod := range g.sdkData.Modules {
		var modMethods []MethodSummary
		for _, m := range mod.Methods {
			modMethods = append(modMethods, MethodSummary{
				CallPath: m.CallPath,
				Url:      mod.Id + "-" + kebabCase(m.Name) + ".html",
				Summary:  m.Summary,
			})
		}

		if err := g.renderSectionPage(tmpl, mod.Id+".html", SectionPageData{
			BasePageData: BasePageData{
				Title:            mod.Name + " Reference",
				Description:      mod.Summary,
				BaseUrl:          g.opts.BaseUrl,
				BaseRelativePath: "",
				CurrentSection:   mod.Id,
				NavModules:       g.navModules,
				NavApis:          g.navApis,
				NavFunctions:     g.navFunctions,
			},
			CallPath:        mod.CallPath,
			SummaryHtml:     template.HTML(RenderMarkdown(mod.Summary, g.linkResolver)),
			DescriptionHtml: template.HTML(RenderMarkdown(mod.Description, g.linkResolver)),
			Methods:         modMethods,
		}); err != nil {
			return err
		}
	}

	// 4. Types section page
	var allTypedefs []SdkTypedef
	for _, t := range g.sdkData.Typedefs {
		allTypedefs = append(allTypedefs, t)
	}
	sort.Slice(allTypedefs, func(i, j int) bool {
		return allTypedefs[i].Name < allTypedefs[j].Name
	})

	return g.renderSectionPage(tmpl, "types.html", SectionPageData{
		BasePageData: BasePageData{
			Title:            "Types & Helpers Reference",
			Description:      "Reference for shared SDK types, constants, and envelopes.",
			BaseUrl:          g.opts.BaseUrl,
			BaseRelativePath: "",
			CurrentSection:   "types",
			NavModules:       g.navModules,
			NavApis:          g.navApis,
			NavFunctions:     g.navFunctions,
		},
		CallPath:    "I6Sdk.Resp / I6Sdk.Grants",
		SummaryHtml: template.HTML("<p>Catalog of shared types, configuration options, and permission bitmasks.</p>"),
		Typedefs:    allTypedefs,
	})
}

func (g *siteGenerator) renderSectionPage(tmpl *template.Template, fileName string, data SectionPageData) error {
	outPath := filepath.Join(g.opts.OutDir, fileName)
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", fileName, err)
	}
	defer f.Close()

	g.recordSitemapUrl(fileName)
	return tmpl.ExecuteTemplate(f, "layout", data)
}

func (g *siteGenerator) renderApiMethods(tmpl *template.Template) error {
	for _, api := range g.apis {
		fnName := apiFuncName(api.Spec.CamelId)
		callPath := "sdk.apis()." + fnName
		name := schemamold.DescName(api.Spec.Desc)
		summary := schemamold.DescSummary(api.Spec.Desc)
		tokenType := apiTokenType(api)

		reqFields := g.flattenApiReqFields(api)
		respFields := g.flattenApiRespFields(api)

		exampleCode := g.buildApiExampleCode(fnName, api)

		slug := "apis/" + api.Spec.Id + ".html"
		pageData := MethodPageData{
			BasePageData: BasePageData{
				Title:            name + " (" + callPath + ")",
				Description:      summary,
				MethodId:         api.Spec.Id,
				HasTryIt:         true,
				BaseUrl:          g.opts.BaseUrl,
				BaseRelativePath: "../",
				CurrentSection:   "apis",
				CurrentMethodId:  api.Spec.Id,
				NavModules:       g.navModules,
				NavApis:          g.navApis,
				NavFunctions:     g.navFunctions,
			},
			CallPath:        callPath,
			Name:            name,
			SummaryHtml:     template.HTML(RenderMarkdown(summary, g.linkResolver)),
			HttpMethod:      api.Spec.Method,
			HttpMethodLower: strings.ToLower(api.Spec.Method),
			HttpPath:        api.Spec.Path,
			TokenType:       tokenType,
			RequestFields:   reqFields,
			ResponseFields:  respFields,
			ExampleCode:     exampleCode,
		}

		outPath := filepath.Join(g.opts.OutDir, slug)
		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("failed to create %s: %w", slug, err)
		}

		err = tmpl.ExecuteTemplate(f, "layout", pageData)
		f.Close()
		if err != nil {
			return fmt.Errorf("failed to render %s: %w", slug, err)
		}

		g.recordSitemapUrl(slug)
		g.methodCatalog[callPath] = MethodCatalogItem{
			CallPath:       callPath,
			Name:           name,
			Summary:        summary,
			HttpMethod:     api.Spec.Method,
			HttpPath:       api.Spec.Path,
			TokenType:      tokenType,
			Url:            slug,
			RequestFields:  reqFields,
			ResponseFields: respFields,
		}
	}
	return nil
}

func (g *siteGenerator) renderFunctionMethods(tmpl *template.Template) error {
	for _, fn := range g.sdkFuncs {
		fnName := sdkFuncName(fn.Spec.CamelId, fn.Spec.Version)
		callPath := "sdk.funcs()." + fnName
		name := schemamold.DescName(fn.Spec.Desc)
		summary := schemamold.DescSummary(fn.Spec.Desc)

		reqFields := g.flattenSdkFuncReqFields(fn)
		respFields := g.flattenSdkFuncRespFields(fn)

		exampleCode := g.buildSdkFuncExampleCode(callPath, fn)

		slug := "functions/" + fn.Spec.Id + ".html"
		pageData := MethodPageData{
			BasePageData: BasePageData{
				Title:            name + " (" + callPath + ")",
				Description:      summary,
				MethodId:         fn.Spec.Id,
				HasTryIt:         true,
				BaseUrl:          g.opts.BaseUrl,
				BaseRelativePath: "../",
				CurrentSection:   "functions",
				CurrentMethodId:  fn.Spec.Id,
				NavModules:       g.navModules,
				NavApis:          g.navApis,
				NavFunctions:     g.navFunctions,
			},
			CallPath:       callPath,
			Name:           name,
			SummaryHtml:    template.HTML(RenderMarkdown(summary, g.linkResolver)),
			RequestFields:  reqFields,
			ResponseFields: respFields,
			ExampleCode:    exampleCode,
		}

		outPath := filepath.Join(g.opts.OutDir, slug)
		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("failed to create %s: %w", slug, err)
		}

		err = tmpl.ExecuteTemplate(f, "layout", pageData)
		f.Close()
		if err != nil {
			return fmt.Errorf("failed to render %s: %w", slug, err)
		}

		g.recordSitemapUrl(slug)
		g.methodCatalog[callPath] = MethodCatalogItem{
			CallPath:       callPath,
			Name:           name,
			Summary:        summary,
			Url:            slug,
			RequestFields:  reqFields,
			ResponseFields: respFields,
		}
	}
	return nil
}

func (g *siteGenerator) renderHandwrittenMethods(tmpl *template.Template) error {
	for _, mod := range g.sdkData.Modules {
		for _, m := range mod.Methods {
			callPath := m.CallPath
			slug := mod.Id + "-" + kebabCase(m.Name) + ".html"

			var reqFields []FieldDoc
			for _, p := range m.Params {
				reqFields = append(reqFields, FieldDoc{
					Name:        p.Name,
					Type:        p.Type,
					Optional:    p.Optional,
					SummaryHtml: template.HTML(RenderMarkdown(p.Doc, g.linkResolver)),
					Anchor:      "param-" + p.Name,
				})

				// If param is a typedef, also list its properties
				if tdef, exists := g.sdkData.Typedefs[p.Type]; exists {
					for _, prop := range tdef.Properties {
						reqFields = append(reqFields, FieldDoc{
							Name:        p.Name + "." + prop.Name,
							Type:        prop.Type,
							Optional:    prop.Optional,
							SummaryHtml: template.HTML(RenderMarkdown(prop.Doc, g.linkResolver)),
							Anchor:      "field-" + p.Name + "-" + prop.Name,
						})
					}
				}
			}

			var respFields []FieldDoc
			if m.ReturnType != "" && m.ReturnType != "void" {
				respFields = append(respFields, FieldDoc{
					Name:        "return",
					Type:        m.ReturnType,
					SummaryHtml: template.HTML(RenderMarkdown(m.ReturnDoc, g.linkResolver)),
					Anchor:      "return-val",
				})
			}

			// Try it allowed for infer and dash
			hasTryIt := mod.Id == "infer" || mod.Id == "dash"
			exampleCode := g.buildHandwrittenExampleCode(m)

			pageData := MethodPageData{
				BasePageData: BasePageData{
					Title:            m.Name + " (" + callPath + ")",
					Description:      m.Summary,
					MethodId:         mod.Id + "-" + m.Name,
					HasTryIt:         hasTryIt,
					BaseUrl:          g.opts.BaseUrl,
					BaseRelativePath: "",
					CurrentSection:   mod.Id,
					CurrentMethodId:  m.Name,
					NavModules:       g.navModules,
					NavApis:          g.navApis,
					NavFunctions:     g.navFunctions,
				},
				CallPath:        callPath,
				Name:            m.Name,
				SummaryHtml:     template.HTML(RenderMarkdown(m.Summary, g.linkResolver)),
				DescriptionHtml: template.HTML(RenderMarkdown(m.Description, g.linkResolver)),
				RequestFields:   reqFields,
				ResponseFields:  respFields,
				ExampleCode:     exampleCode,
			}

			outPath := filepath.Join(g.opts.OutDir, slug)
			f, err := os.Create(outPath)
			if err != nil {
				return fmt.Errorf("failed to create %s: %w", slug, err)
			}

			err = tmpl.ExecuteTemplate(f, "layout", pageData)
			f.Close()
			if err != nil {
				return fmt.Errorf("failed to render %s: %w", slug, err)
			}

			g.recordSitemapUrl(slug)
			g.methodCatalog[callPath] = MethodCatalogItem{
				CallPath:       callPath,
				Name:           m.Name,
				Summary:        m.Summary,
				Url:            slug,
				RequestFields:  reqFields,
				ResponseFields: respFields,
			}
		}
	}
	return nil
}

func (g *siteGenerator) renderSitemap(tmpl *texttemplate.Template) error {
	outPath := filepath.Join(g.opts.OutDir, "sitemap.xml")
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create sitemap.xml: %w", err)
	}
	defer f.Close()

	var urls []string
	sort.Strings(g.sitemapUrls)
	for _, u := range g.sitemapUrls {
		urls = append(urls, g.opts.BaseUrl+"/"+u)
	}

	return tmpl.ExecuteTemplate(f, "sitemap.xml.tmpl", map[string]any{"Urls": urls})
}

func (g *siteGenerator) renderStaticHtml(tmpl *template.Template, templateName, outputName string) error {
	outPath := filepath.Join(g.opts.OutDir, outputName)
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", outputName, err)
	}
	defer f.Close()

	g.recordSitemapUrl(outputName)
	return tmpl.ExecuteTemplate(f, templateName, nil)
}

func (g *siteGenerator) recordSitemapUrl(u string) {
	g.sitemapUrls = append(g.sitemapUrls, strings.TrimPrefix(u, "/"))
}

func (g *siteGenerator) copyAssets() {
	// Copy src/site/docs/* if present
	siteDocsDir := filepath.Join("src", "site", "docs")
	if dirExists(siteDocsDir) {
		entries, err := os.ReadDir(siteDocsDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				srcFile := filepath.Join(siteDocsDir, entry.Name())
				dstFile := filepath.Join(g.opts.OutDir, entry.Name())
				copyFile(srcFile, dstFile)
			}
		}
	}

	// Copy the SDK min.js if present
	copyFile(filepath.Join("dist", "bundle", "i6sdk-web.min.js"), filepath.Join(g.opts.OutDir, "i6sdk-web.min.js"))
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func copyFile(src, dst string) {
	if !filez.FileExists(src) {
		return
	}
	sf, err := os.Open(src)
	if err != nil {
		return
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return
	}
	defer df.Close()

	_, _ = io.Copy(df, sf)
}

// typedefFields flattens the groups into dotted rows (params.dataset,
// payload.store.id), following a nested record one level. anchorPrefix
// keeps request and response anchors apart on the same page.
func (g *siteGenerator) typedefFields(tdfs typeastz.Typedefs, groups []fieldGroup, anchorPrefix string) []FieldDoc {
	var fields []FieldDoc
	add := func(fullName string, p typeastz.Tuple) {
		fields = append(fields, FieldDoc{
			Name:        fullName,
			Type:        strings.TrimSuffix(p.Typ, "?"),
			Optional:    strings.HasSuffix(p.Typ, "?") || strings.Contains(strings.ToLower(p.Summary), "optional"),
			SummaryHtml: template.HTML(RenderMarkdown(p.Summary, g.linkResolver)),
			Anchor:      anchorPrefix + strings.ReplaceAll(fullName, ".", "-"),
		})
	}
	for _, fg := range groups {
		td := findTypedef(tdfs, fg.typeName())
		if td == nil || td.Properties == nil {
			continue
		}
		for _, p := range td.Properties.All() {
			fullName := fg.prefix + "." + p.Name
			add(fullName, p)
			if nested := findTypedef(tdfs, p.Typ); nested != nil && nested.Properties != nil {
				for _, np := range nested.Properties.All() {
					add(fullName+"."+np.Name, np)
				}
			}
		}
	}
	return fields
}

func (g *siteGenerator) flattenApiReqFields(api *apispecz.Api) []FieldDoc {
	return g.typedefFields(api.Typedefs, []fieldGroup{
		{"params", api.Req.Params},
		{"query", api.Req.Query},
		{"headers", api.Req.Headers},
		{"payload", api.Req.Payload},
	}, "field-")
}

func (g *siteGenerator) flattenApiRespFields(api *apispecz.Api) []FieldDoc {
	return g.typedefFields(api.Typedefs, []fieldGroup{
		{"payload", api.Resp.Payload},
		{"headers", api.Resp.Headers},
	}, "resp-")
}

func (g *siteGenerator) flattenSdkFuncReqFields(fn *moldsdkfunc.SdkFunc) []FieldDoc {
	return g.typedefFields(fn.Typedefs, []fieldGroup{
		{"input", fn.Req.Input},
		{"partitions", fn.Req.Partitions},
	}, "field-")
}

func (g *siteGenerator) flattenSdkFuncRespFields(fn *moldsdkfunc.SdkFunc) []FieldDoc {
	return g.typedefFields(fn.Typedefs, []fieldGroup{
		{"output", fn.Resp.Output},
	}, "resp-")
}

// exampleProp is one property of a type, for building example values.
type exampleProp struct {
	name  string
	typ   string
	value string // code to use instead of a placeholder
}

// exampleTypes finds a type by name: its properties if it is a record, or
// the type it stands for if it is an alias.
type exampleTypes func(name string) (props []exampleProp, alias string, ok bool)

func forgeExampleTypes(tdfs typeastz.Typedefs) exampleTypes {
	return func(name string) ([]exampleProp, string, bool) {
		td := findTypedef(tdfs, name)
		if td == nil {
			return nil, "", false
		}
		if td.Properties == nil {
			return nil, td.Head.Typ, true
		}
		var props []exampleProp
		for _, p := range td.Properties.All() {
			props = append(props, exampleProp{name: p.Name, typ: p.Typ})
		}
		return props, "", true
	}
}

func (g *siteGenerator) sdkExampleTypes(name string) ([]exampleProp, string, bool) {
	td, ok := g.sdkData.Typedefs[name]
	if !ok || len(td.Properties) == 0 {
		return nil, "", false
	}
	var props []exampleProp
	for _, p := range td.Properties {
		typ := p.Type
		if p.Optional {
			typ += "?"
		}
		props = append(props, exampleProp{name: p.Name, typ: typ})
	}
	return props, "", true
}

// exampleValue is a placeholder of type t that type checks: "" for strings,
// 0 for numbers, and records filled in property by property.
func exampleValue(t string, types exampleTypes, indent string, depth int) string {
	t = strings.TrimSpace(strings.TrimSuffix(t, "?"))
	if !strings.HasPrefix(t, "{") {
		switch {
		case strings.Contains(t, "=> Promise<"):
			return "async () => {}"
		case strings.Contains(t, "=>") || strings.Contains(t, "Function"):
			return "() => {}"
		}
	}
	if i := strings.Index(t, "|"); i >= 0 {
		t = strings.TrimSpace(t[:i])
	}
	switch {
	case t == "string":
		return `""`
	case t == "number":
		return "0"
	case t == "boolean":
		return "false"
	case t == "any" || t == "unknown":
		return "{}"
	case strings.HasPrefix(t, `"`) || (t != "" && strings.ContainsRune("0123456789-", rune(t[0]))):
		return t // a literal
	case strings.HasPrefix(t, "Array<") || strings.HasSuffix(t, "[]"):
		return "[]"
	case strings.HasPrefix(t, "Record<") || strings.HasPrefix(t, "{"):
		return "{}"
	case t == "File":
		return `new File([], "file.txt")`
	}
	props, alias, ok := types(t)
	if !ok || depth > 3 {
		return "{}"
	}
	if alias != "" {
		return exampleValue(alias, types, indent, depth+1)
	}
	return exampleObject(props, types, indent, depth+1)
}

func exampleObject(props []exampleProp, types exampleTypes, indent string, depth int) string {
	required := 0
	for _, p := range props {
		if !strings.HasSuffix(p.typ, "?") {
			required++
		}
	}
	if required == 0 {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteString("{\n")
	for _, p := range props {
		if strings.HasSuffix(p.typ, "?") {
			continue // optional: listed in the field table, left out of the example
		}
		value := p.value
		if value == "" {
			value = exampleValue(p.typ, types, indent+"  ", depth)
		}
		sb.WriteString(fmt.Sprintf("%s  %s: %s,\n", indent, p.name, value))
	}
	sb.WriteString(indent + "}")
	return sb.String()
}

// exampleCall is the code of a Try it box. The result is kept in a const
// before it is returned: the editor does not type check what follows a
// top-level return.
func exampleCall(callPath string, args []string) string {
	return fmt.Sprintf("const result = await %s(%s)\nreturn result", callPath, strings.Join(args, ", "))
}

// exampleGroups is the request object of an API or SdkFunc: one property per
// group (params, payload, ...) that has a type.
func exampleGroups(groups []fieldGroup, types exampleTypes) string {
	var props []exampleProp
	for _, fg := range groups {
		if fg.typeName() != "" {
			props = append(props, exampleProp{name: fg.prefix, typ: fg.typeName()})
		}
	}
	return exampleObject(props, types, "", 0)
}

func (g *siteGenerator) buildApiExampleCode(fnName string, api *apispecz.Api) string {
	types := forgeExampleTypes(api.Typedefs)
	// An API with an authorization header takes the user's token as is (the
	// SDK adds it only for APIs labeled with token_types).
	headersType := api.Req.Headers.Head
	withAuth := func(name string) ([]exampleProp, string, bool) {
		props, alias, ok := types(name)
		if headersType != nil && name == headersType.Typ {
			for i := range props {
				if strings.EqualFold(props[i].name, "authorization") {
					props[i].value = `"Bearer " + (await sdk.auth().getValidAccessToken())`
				}
			}
		}
		return props, alias, ok
	}
	groups := []fieldGroup{
		{"params", api.Req.Params},
		{"query", api.Req.Query},
		{"headers", api.Req.Headers},
		{"payload", api.Req.Payload},
	}
	return exampleCall("sdk.apis()."+fnName, []string{exampleGroups(groups, withAuth)})
}

func (g *siteGenerator) buildSdkFuncExampleCode(callPath string, fn *moldsdkfunc.SdkFunc) string {
	groups := []fieldGroup{
		{"input", fn.Req.Input},
		{"partitions", fn.Req.Partitions},
	}
	return exampleCall(callPath, []string{exampleGroups(groups, forgeExampleTypes(fn.Typedefs))})
}

// buildHandwrittenExampleCode passes the required params, in order.
func (g *siteGenerator) buildHandwrittenExampleCode(m SdkMethod) string {
	var args []string
	for _, p := range m.Params {
		if p.Optional {
			break
		}
		args = append(args, exampleValue(p.Type, g.sdkExampleTypes, "", 0))
	}
	return exampleCall(m.CallPath, args)
}

func apiFuncName(apiCamelId string) string {
	name := strings.TrimPrefix(apiCamelId, "api")
	if len(name) == 0 {
		return name
	}
	return strings.ToLower(name[0:1]) + name[1:]
}

func sdkFuncName(funcCamelId, version string) string {
	name := strings.TrimPrefix(funcCamelId, "sdkFunc")
	if len(name) == 0 {
		return name
	}
	name = strings.ToLower(name[0:1]) + name[1:]
	if version != "" {
		name += camelz.P(version).P()
	}
	return name
}

func apiTokenType(api *apispecz.Api) string {
	if api.Spec.Desc == nil || api.Spec.Desc.Labels == nil {
		return ""
	}
	return api.Spec.Desc.Labels["token_types"]
}

func kebabCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteRune('-')
			}
			b.WriteRune(r + ('a' - 'A'))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
