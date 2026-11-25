package schemagen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/infinity6-ai/gox/commonz/stringz/camelz"
)

// GenerateDir parses Go source files in the specified directory, identifies the Api struct,
// and writes out <package>.gen.go implementing func (d *Api) ApiDataRefs() *apiz.DataRefs.
func GenerateDir(dir string) (string, error) {
	code, pkgName, err := ParseAndGenerate(dir)
	if err != nil {
		return "", fmt.Errorf("failed to generate schema for %s: %w", dir, err)
	}

	outFile := filepath.Join(dir, pkgName+".gen.go")
	if err := os.WriteFile(outFile, code, 0644); err != nil {
		return "", fmt.Errorf("failed to write %s: %w", outFile, err)
	}
	return outFile, nil
}

// ParseAndGenerate parses the package in dir and generates the formatted Go source code.
func ParseAndGenerate(dir string) ([]byte, string, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi fs.FileInfo) bool {
		name := fi.Name()
		return !strings.HasSuffix(name, "_test.go") && !strings.HasSuffix(name, ".gen.go")
	}, parser.ParseComments)
	if err != nil {
		return nil, "", fmt.Errorf("error parsing directory %s: %w", dir, err)
	}

	if len(pkgs) == 0 {
		return nil, "", fmt.Errorf("no Go package found in %s", dir)
	}

	var astPkg *ast.Package
	var pkgName string
	for name, p := range pkgs {
		astPkg = p
		pkgName = name
		break
	}

	docPkg := doc.New(astPkg, "./", doc.AllDecls)

	typeInfo := collectTypeInfo(astPkg, docPkg)

	apiStruct, ok := typeInfo.structs["Api"]
	if !ok {
		return nil, "", fmt.Errorf("no 'type Api struct' found in package %s", pkgName)
	}

	code, err := generateApiDataRefs(pkgName, apiStruct, typeInfo)
	if err != nil {
		return nil, "", fmt.Errorf("error generating ApiDataRefs: %w", err)
	}

	formatted, err := format.Source(code)
	if err != nil {
		return nil, "", fmt.Errorf("error formatting generated source: %w\nSource:\n%s", err, string(code))
	}

	return formatted, pkgName, nil
}

type typeInfo struct {
	structs     map[string]*ast.StructType
	typeDocs    map[string]string
	rawTypeDocs map[string]string
	fieldDocs   map[string]map[string]string // structName -> fieldName -> doc
	fieldTypes  map[string]map[string]string // structName -> fieldName -> typeString
	imports       map[string]string            // alias -> importPath
	hasApiRoute   bool
	hasApiSpec    bool
	apiMethodName string
	pkgDoc        string
}

func collectTypeInfo(pkg *ast.Package, docPkg *doc.Package) *typeInfo {
	info := &typeInfo{
		structs:     make(map[string]*ast.StructType),
		typeDocs:    make(map[string]string),
		rawTypeDocs: make(map[string]string),
		fieldDocs:   make(map[string]map[string]string),
		fieldTypes:  make(map[string]map[string]string),
		imports:     make(map[string]string),
		pkgDoc:      docPkg.Doc,
	}

	for _, file := range pkg.Files {
		for _, imp := range file.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			alias := ""
			if imp.Name != nil {
				alias = imp.Name.Name
			} else {
				alias = filepath.Base(path)
			}
			info.imports[alias] = path
		}

		for _, decl := range file.Decls {
			if fd, ok := decl.(*ast.FuncDecl); ok && fd.Recv != nil && len(fd.Recv.List) > 0 {
				recvType, _ := resolveTypeName(fd.Recv.List[0].Type)
				if recvType == "Api" {
					if fd.Name.Name == "ApiRoute" {
						info.hasApiRoute = true
					}
					if fd.Name.Name == "ApiSpec" {
						info.hasApiSpec = true
					}
				}
				if strings.HasSuffix(recvType, "Id") && recvType != "Api" {
					info.apiMethodName = fd.Name.Name
				}
			}
		}
	}

	for _, t := range docPkg.Types {
		info.rawTypeDocs[t.Name] = t.Doc
		info.typeDocs[t.Name] = cleanDoc(t.Doc)
		info.fieldDocs[t.Name] = make(map[string]string)
		info.fieldTypes[t.Name] = make(map[string]string)

		if t.Decl != nil {
			for _, spec := range t.Decl.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				info.structs[ts.Name.Name] = st

				if st.Fields != nil {
					for _, f := range st.Fields.List {
						docStr := ""
						if f.Doc != nil {
							docStr = cleanDoc(f.Doc.Text())
						} else if f.Comment != nil {
							docStr = cleanDoc(f.Comment.Text())
						}
						typeStr := exprToString(f.Type)
						for _, ident := range f.Names {
							info.fieldDocs[ts.Name.Name][ident.Name] = docStr
							info.fieldTypes[ts.Name.Name][ident.Name] = typeStr
						}
					}
				}
			}
		}
	}

	return info
}

func generateApiDataRefs(pkgName string, apiStruct *ast.StructType, info *typeInfo) ([]byte, error) {
	var bodyBuf bytes.Buffer
	neededImports := make(map[string]string)

	reqField := findFieldInStruct(apiStruct, "Req")
	respField := findFieldInStruct(apiStruct, "Resp")

	if info.hasApiRoute && !info.hasApiSpec {
		apiDoc := info.rawTypeDocs["Api"]
		if apiDoc == "" {
			apiDoc = info.pkgDoc
		}
		title, summary := splitDocTitleAndSummary(apiDoc)
		if title == "" {
			title = pkgName
		}
		if summary == "" {
			summary = title
		}

		bodyBuf.WriteString("func (d *Api) ApiSpec() apiz.ApiSpec {\n")
		bodyBuf.WriteString("\tspec := d.ApiRoute()\n")
		if info.apiMethodName != "" {
			inferredId := camelz.P(info.apiMethodName).QL()
			bodyBuf.WriteString("\tif spec.Id == \"\" {\n")
			bodyBuf.WriteString(fmt.Sprintf("\t\tspec.Id = %q\n", inferredId))
			bodyBuf.WriteString("\t}\n")
		}
		bodyBuf.WriteString("\torigDesc := spec.Desc\n")
		bodyBuf.WriteString("\tspec.Desc = func() *schemaz.Desc {\n")
		bodyBuf.WriteString("\t\tvar desc *schemaz.Desc\n")
		bodyBuf.WriteString("\t\tif origDesc != nil {\n")
		bodyBuf.WriteString("\t\t\tdesc = origDesc()\n")
		bodyBuf.WriteString("\t\t}\n")
		bodyBuf.WriteString("\t\tif desc == nil {\n")
		bodyBuf.WriteString("\t\t\tdesc = &schemaz.Desc{}\n")
		bodyBuf.WriteString("\t\t}\n")
		bodyBuf.WriteString("\t\tif desc.Name == \"\" {\n")
		bodyBuf.WriteString(fmt.Sprintf("\t\t\tdesc.Name = %q\n", title))
		bodyBuf.WriteString("\t\t}\n")
		bodyBuf.WriteString("\t\tif desc.Summary == \"\" {\n")
		bodyBuf.WriteString(fmt.Sprintf("\t\t\tdesc.Summary = %q\n", summary))
		bodyBuf.WriteString("\t\t}\n")
		bodyBuf.WriteString("\t\treturn desc\n")
		bodyBuf.WriteString("\t}\n")
		bodyBuf.WriteString("\treturn spec\n")
		bodyBuf.WriteString("}\n\n")
	}

	bodyBuf.WriteString("func (d *Api) ApiDataRefs() *apiz.DataRefs {\n")
	bodyBuf.WriteString("\treturn &apiz.DataRefs{\n")

	// Determine Req slots
	var paramsStruct, queryStruct, reqHeadersStruct, reqPayloadStruct *ast.StructType
	var paramsTypeName, queryTypeName, reqHeadersTypeName, reqPayloadTypeName string
	var reqIsPtr bool
	var paramsIsPtr, queryIsPtr, reqHeadersIsPtr, reqPayloadIsPtr bool
	var reqPayloadFieldName string

	if reqField != nil {
		reqTypeName, isPtr := resolveTypeName(reqField.Type)
		reqIsPtr = isPtr
		if reqSt, ok := info.structs[reqTypeName]; ok {
			if f := findFieldInStruct(reqSt, "Params"); f != nil {
				paramsTypeName, paramsIsPtr = resolveTypeName(f.Type)
				paramsStruct = info.structs[paramsTypeName]
			}
			if f := findFieldInStruct(reqSt, "Query"); f != nil {
				queryTypeName, queryIsPtr = resolveTypeName(f.Type)
				queryStruct = info.structs[queryTypeName]
			}
			if f := findFieldInStruct(reqSt, "Headers"); f != nil {
				reqHeadersTypeName, reqHeadersIsPtr = resolveTypeName(f.Type)
				reqHeadersStruct = info.structs[reqHeadersTypeName]
			}
			if f := findFieldInStruct(reqSt, "Payload"); f != nil {
				reqPayloadFieldName = "Payload"
				reqPayloadTypeName, reqPayloadIsPtr = resolveTypeName(f.Type)
				reqPayloadStruct = info.structs[reqPayloadTypeName]
			} else if f := findFieldInStruct(reqSt, "Body"); f != nil {
				reqPayloadFieldName = "Body"
				reqPayloadTypeName, reqPayloadIsPtr = resolveTypeName(f.Type)
				reqPayloadStruct = info.structs[reqPayloadTypeName]
			}
		}
	}

	// Determine Resp slots
	var respHeadersStruct, respPayloadStruct *ast.StructType
	var respHeadersTypeName, respPayloadTypeName string
	var respIsPtr bool
	var respHeadersIsPtr, respPayloadIsPtr bool
	var respPayloadFieldName string

	if respField != nil {
		respTypeName, isPtr := resolveTypeName(respField.Type)
		respIsPtr = isPtr
		if respSt, ok := info.structs[respTypeName]; ok {
			if f := findFieldInStruct(respSt, "Headers"); f != nil {
				respHeadersTypeName, respHeadersIsPtr = resolveTypeName(f.Type)
				respHeadersStruct = info.structs[respHeadersTypeName]
			}
			if f := findFieldInStruct(respSt, "Payload"); f != nil {
				respPayloadFieldName = "Payload"
				respPayloadTypeName, respPayloadIsPtr = resolveTypeName(f.Type)
				respPayloadStruct = info.structs[respPayloadTypeName]
			} else if f := findFieldInStruct(respSt, "Body"); f != nil {
				respPayloadFieldName = "Body"
				respPayloadTypeName, respPayloadIsPtr = resolveTypeName(f.Type)
				respPayloadStruct = info.structs[respPayloadTypeName]
			}
		}
	}

	// 1. PathParams
	if paramsStruct != nil {
		writeUrlValueSlot(&bodyBuf, "PathParams", "Params", paramsTypeName, paramsStruct, reqIsPtr, paramsIsPtr, info, neededImports)
	} else {
		bodyBuf.WriteString("\t\tPathParams: nil,\n")
	}

	// 2. QueryParams
	if queryStruct != nil {
		writeUrlValueSlot(&bodyBuf, "QueryParams", "Query", queryTypeName, queryStruct, reqIsPtr, queryIsPtr, info, neededImports)
	} else {
		bodyBuf.WriteString("\t\tQueryParams: nil,\n")
	}

	// 3. ReqHeaders
	if reqHeadersStruct != nil {
		writeUrlValueSlot(&bodyBuf, "ReqHeaders", "Headers", reqHeadersTypeName, reqHeadersStruct, reqIsPtr, reqHeadersIsPtr, info, neededImports)
	} else {
		bodyBuf.WriteString("\t\tReqHeaders: nil,\n")
	}

	// 4. ReqBody
	if reqPayloadStruct != nil {
		writePayloadSlot(&bodyBuf, "ReqBody", "d.Req."+reqPayloadFieldName, reqPayloadFieldName, reqPayloadTypeName, reqPayloadStruct, reqIsPtr, reqPayloadIsPtr, info)
	} else {
		bodyBuf.WriteString("\t\tReqBody: nil,\n")
	}

	// 5. RespHeaders
	if respHeadersStruct != nil {
		writeUrlValueSlot(&bodyBuf, "RespHeaders", "Headers", respHeadersTypeName, respHeadersStruct, respIsPtr, respHeadersIsPtr, info, neededImports)
	} else {
		bodyBuf.WriteString("\t\tRespHeaders: nil,\n")
	}

	// 6. RespBody
	if respPayloadStruct != nil {
		writePayloadSlot(&bodyBuf, "RespBody", "d.Resp."+respPayloadFieldName, respPayloadFieldName, respPayloadTypeName, respPayloadStruct, respIsPtr, respPayloadIsPtr, info)
	} else {
		bodyBuf.WriteString("\t\tRespBody: nil,\n")
	}

	bodyBuf.WriteString("\t}\n")
	bodyBuf.WriteString("}\n")

	var b bytes.Buffer
	b.WriteString("// Code generated by schemagen. DO NOT EDIT.\n\n")
	b.WriteString(fmt.Sprintf("package %s\n\n", pkgName))
	b.WriteString("import (\n")
	b.WriteString("\t\"github.com/infinity6-ai/gox/routez/apiz\"\n")
	b.WriteString("\t\"github.com/infinity6-ai/gox/schemaz/schemaz\"\n")
	for _, impPath := range neededImports {
		b.WriteString(fmt.Sprintf("\t%q\n", impPath))
	}
	b.WriteString(")\n\n")
	b.Write(bodyBuf.Bytes())

	return b.Bytes(), nil
}

func writeUrlValueSlot(b *bytes.Buffer, slotName, slotField, typeName string, st *ast.StructType, reqIsPtr, slotIsPtr bool, info *typeInfo, neededImports map[string]string) {
	b.WriteString(fmt.Sprintf("\t\t%s: &schemaz.Schema{\n", slotName))
	b.WriteString("\t\t\tObject: func(read bool) map[string]*schemaz.Schema {\n")

	if reqIsPtr {
		b.WriteString("\t\t\t\tif read && d.Req == nil {\n\t\t\t\t\treturn nil\n\t\t\t\t}\n")
		b.WriteString("\t\t\t\tif d.Req == nil {\n\t\t\t\t\td.Req = &Req{}\n\t\t\t\t}\n")
	}
	if slotIsPtr {
		b.WriteString(fmt.Sprintf("\t\t\t\tif read && d.Req.%s == nil {\n\t\t\t\t\treturn nil\n\t\t\t\t}\n", slotField))
		b.WriteString(fmt.Sprintf("\t\t\t\tif d.Req.%s == nil {\n\t\t\t\t\td.Req.%s = &%s{}\n\t\t\t\t}\n", slotField, slotField, typeName))
	}

	b.WriteString("\t\t\t\treturn map[string]*schemaz.Schema{\n")
	if st.Fields != nil {
		for _, f := range st.Fields.List {
			jsonTag := extractJsonName(f.Tag)
			if jsonTag == "-" || jsonTag == "" {
				continue
			}
			for _, ident := range f.Names {
				docStr := info.fieldDocs[typeName][ident.Name]
				urlParser := resolveUrlParser(f.Type, fmt.Sprintf("&d.Req.%s.%s", slotField, ident.Name), neededImports, info.imports)
				b.WriteString(fmt.Sprintf("\t\t\t\t\t%q: {\n", jsonTag))
				if docStr != "" {
					b.WriteString(fmt.Sprintf("\t\t\t\t\t\tDesc: &schemaz.Desc{Summary: %q},\n", docStr))
				}
				b.WriteString(fmt.Sprintf("\t\t\t\t\t\tUrlValue: %s,\n", urlParser))
				b.WriteString("\t\t\t\t\t},\n")
			}
		}
	}
	b.WriteString("\t\t\t\t}\n")
	b.WriteString("\t\t\t},\n")
	b.WriteString("\t\t},\n")
}

func writePayloadSlot(b *bytes.Buffer, slotName, accessor, slotField, typeName string, st *ast.StructType, parentIsPtr, slotIsPtr bool, info *typeInfo) {
	b.WriteString(fmt.Sprintf("\t\t%s: &schemaz.Schema{\n", slotName))
	docSummary := info.typeDocs[typeName]
	if docSummary != "" {
		b.WriteString(fmt.Sprintf("\t\t\tDesc: &schemaz.Desc{Summary: %q},\n", docSummary))
	}
	b.WriteString("\t\t\tObject: func(read bool) map[string]*schemaz.Schema {\n")

	parentPrefix := "d.Resp"
	parentType := "Resp"
	if strings.HasPrefix(accessor, "d.Req") {
		parentPrefix = "d.Req"
		parentType = "Req"
	}

	if parentIsPtr {
		b.WriteString(fmt.Sprintf("\t\t\t\tif read && %s == nil {\n\t\t\t\t\treturn nil\n\t\t\t\t}\n", parentPrefix))
		b.WriteString(fmt.Sprintf("\t\t\t\tif %s == nil {\n\t\t\t\t\t%s = &%s{}\n\t\t\t\t}\n", parentPrefix, parentPrefix, parentType))
	}
	if slotIsPtr {
		b.WriteString(fmt.Sprintf("\t\t\t\tif read && %s == nil {\n\t\t\t\t\treturn nil\n\t\t\t\t}\n", accessor))
		b.WriteString(fmt.Sprintf("\t\t\t\tif %s == nil {\n\t\t\t\t\t%s = &%s{}\n\t\t\t\t}\n", accessor, accessor, typeName))
	}

	b.WriteString("\t\t\t\treturn map[string]*schemaz.Schema{\n")
	if st.Fields != nil {
		for _, f := range st.Fields.List {
			jsonTag := extractJsonName(f.Tag)
			if jsonTag == "-" || jsonTag == "" {
				continue
			}
			for _, ident := range f.Names {
				fieldDoc := info.fieldDocs[typeName][ident.Name]
				fieldAccessor := fmt.Sprintf("%s.%s", accessor, ident.Name)
				writeFieldSchema(b, "\t\t\t\t\t", jsonTag, fieldDoc, f.Type, fieldAccessor, info)
			}
		}
	}
	b.WriteString("\t\t\t\t}\n")
	b.WriteString("\t\t\t},\n")
	b.WriteString("\t\t},\n")
}

func writeFieldSchema(b *bytes.Buffer, indent, jsonTag, docStr string, expr ast.Expr, accessor string, info *typeInfo) {
	b.WriteString(fmt.Sprintf("%s%q: {\n", indent, jsonTag))
	if docStr != "" {
		b.WriteString(fmt.Sprintf("%s\tDesc: &schemaz.Desc{Summary: %q},\n", indent, docStr))
	}

	// Check if slice/array
	if arrType, ok := expr.(*ast.ArrayType); ok {
		elemTypeStr, isPtr := resolveTypeName(arrType.Elt)
		sliceType := fmt.Sprintf("[]*%s", elemTypeStr)
		elemType := fmt.Sprintf("*%s", elemTypeStr)
		if !isPtr {
			sliceType = fmt.Sprintf("[]%s", elemTypeStr)
			elemType = elemTypeStr
		}

		if elemStruct, ok := info.structs[elemTypeStr]; ok {
			elemDoc := info.typeDocs[elemTypeStr]
			b.WriteString(fmt.Sprintf("%s\tArray: schemaz.NewArray(schemaz.ArrayOptions[%s, %s]{\n", indent, sliceType, elemType))
			b.WriteString(fmt.Sprintf("%s\t\tSlice: func() *%s { return &%s },\n", indent, sliceType, accessor))
			if isPtr {
				b.WriteString(fmt.Sprintf("%s\t\tNewElement: func(idx int) %s { return &%s{} },\n", indent, elemType, elemTypeStr))
				b.WriteString(fmt.Sprintf("%s\t\tGet: func(idx int, element *%s) *schemaz.Schema {\n", indent, elemType))
			} else {
				b.WriteString(fmt.Sprintf("%s\t\tNewElement: func(idx int) %s { return %s{} },\n", indent, elemType, elemTypeStr))
				b.WriteString(fmt.Sprintf("%s\t\tGet: func(idx int, element *%s) *schemaz.Schema {\n", indent, elemType))
			}
			b.WriteString(fmt.Sprintf("%s\t\t\treturn &schemaz.Schema{\n", indent))
			if elemDoc != "" {
				b.WriteString(fmt.Sprintf("%s\t\t\t\tDesc: &schemaz.Desc{Summary: %q},\n", indent, elemDoc))
			}
			b.WriteString(fmt.Sprintf("%s\t\t\t\tObject: func(_ bool) map[string]*schemaz.Schema {\n", indent))
			b.WriteString(fmt.Sprintf("%s\t\t\t\t\treturn map[string]*schemaz.Schema{\n", indent))

			for _, ef := range elemStruct.Fields.List {
				eTag := extractJsonName(ef.Tag)
				if eTag == "-" || eTag == "" {
					continue
				}
				for _, eIdent := range ef.Names {
					eDoc := info.fieldDocs[elemTypeStr][eIdent.Name]
					var eAccessor string
					if isPtr {
						eAccessor = fmt.Sprintf("(*element).%s", eIdent.Name)
					} else {
						eAccessor = fmt.Sprintf("element.%s", eIdent.Name)
					}
					b.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t%q: {\n", indent, eTag))
					if eDoc != "" {
						b.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\tDesc: &schemaz.Desc{Summary: %q},\n", indent, eDoc))
					}
					b.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t\tField: func() any { return &%s },\n", indent, eAccessor))
					b.WriteString(fmt.Sprintf("%s\t\t\t\t\t\t},\n", indent))
				}
			}
			b.WriteString(fmt.Sprintf("%s\t\t\t\t\t}\n", indent))
			b.WriteString(fmt.Sprintf("%s\t\t\t\t},\n", indent))
			b.WriteString(fmt.Sprintf("%s\t\t\t}\n", indent))
			b.WriteString(fmt.Sprintf("%s\t\t},\n", indent))
			b.WriteString(fmt.Sprintf("%s\t}),\n", indent))
			b.WriteString(fmt.Sprintf("%s},\n", indent))
			return
		}
	}

	// Normal Field binding
	b.WriteString(fmt.Sprintf("%s\tField: func() any { return &%s },\n", indent, accessor))
	b.WriteString(fmt.Sprintf("%s},\n", indent))
}

func resolveUrlParser(expr ast.Expr, ptrAccessor string, neededImports map[string]string, availableImports map[string]string) string {
	typeName, _ := resolveTypeName(expr)
	if strings.Contains(typeName, ".") {
		parts := strings.Split(typeName, ".")
		pkgAlias := parts[0]
		if pkgAlias == "lakedsz" {
			if impPath, ok := availableImports[pkgAlias]; ok {
				neededImports[pkgAlias] = impPath
			}
			return fmt.Sprintf("schemaz.ParseStrGet(%s, lakedsz.MustParse)", ptrAccessor)
		}
		if pkgAlias == "solutiondomain" {
			if impPath, ok := availableImports[pkgAlias]; ok {
				neededImports[pkgAlias] = impPath
			}
			return fmt.Sprintf("schemaz.ParseStrGet(%s, solutiondomain.MustParse)", ptrAccessor)
		}
	}
	switch typeName {
	case "int", "int8", "int16", "int32", "int64":
		return fmt.Sprintf("schemaz.ParseInt(%s)", ptrAccessor)
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return fmt.Sprintf("schemaz.ParseUint(%s)", ptrAccessor)
	case "bool":
		return fmt.Sprintf("schemaz.ParseBool(%s)", ptrAccessor)
	case "float32", "float64":
		return fmt.Sprintf("schemaz.ParseFloat(%s)", ptrAccessor)
	default:
		return fmt.Sprintf("schemaz.ParseStr(%s)", ptrAccessor)
	}
}

func resolveTypeName(expr ast.Expr) (string, bool) {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name, false
	case *ast.StarExpr:
		name, _ := resolveTypeName(t.X)
		return name, true
	case *ast.SelectorExpr:
		pkgName, _ := resolveTypeName(t.X)
		return fmt.Sprintf("%s.%s", pkgName, t.Sel.Name), false
	default:
		return exprToString(expr), false
	}
}

func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	default:
		return ""
	}
}

func findFieldInStruct(st *ast.StructType, names ...string) *ast.Field {
	if st == nil || st.Fields == nil {
		return nil
	}
	for _, f := range st.Fields.List {
		for _, ident := range f.Names {
			for _, n := range names {
				if strings.EqualFold(ident.Name, n) {
					return f
				}
			}
		}
	}
	return nil
}

func extractJsonName(tagLit *ast.BasicLit) string {
	if tagLit == nil {
		return ""
	}
	tagVal := strings.Trim(tagLit.Value, "`")
	stTag := reflect.StructTag(tagVal)
	jsonTag := stTag.Get("json")
	if jsonTag == "" {
		return ""
	}
	parts := strings.Split(jsonTag, ",")
	return parts[0]
}

func cleanDoc(text string) string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	var cleaned []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return strings.Join(cleaned, " ")
}

func splitDocTitleAndSummary(rawDoc string) (string, string) {
	trimmed := strings.TrimSpace(rawDoc)
	if trimmed == "" {
		return "", ""
	}
	lines := strings.Split(trimmed, "\n")
	var cleanLines []string
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if t != "" {
			cleanLines = append(cleanLines, t)
		}
	}
	if len(cleanLines) == 0 {
		return "", ""
	}
	if len(cleanLines) == 1 {
		line := cleanLines[0]
		if idx := strings.Index(line, ". "); idx > 0 {
			return strings.TrimSuffix(line[:idx], "."), strings.TrimSpace(line[idx+2:])
		}
		return line, line
	}
	title := strings.TrimSuffix(cleanLines[0], ".")
	summary := strings.Join(cleanLines[1:], " ")
	return title, summary
}
