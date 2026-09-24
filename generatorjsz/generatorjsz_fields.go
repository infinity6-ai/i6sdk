package generatorjsz

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

const maxDepth = 8

// field é um campo de uma seção do DataRefs (path, query, headers, body) ou de um struct Go.
type field struct {
	name    string
	typ     string // tipo JSDoc
	summary string
	// str é true para Str/Strs: no corpo viajam como array JSON de strings.
	str bool
}

func (f field) list() bool { return strings.HasSuffix(f.typ, "[]") }

// sectionFields lê uma seção de nível superior do DataRefs. Só Object com filhos Raw/Str/Strs
// é suportado; nil vira lista vazia (a Api não tem essa seção).
func sectionFields(section string, s *schemaz.Schema) ([]field, error) {
	if s == nil {
		return nil, nil
	}

	if s.Object == nil {
		return nil, fmt.Errorf("generatorjsz: %s: only Object is implemented", section)
	}

	m := s.Object(false)
	names := make([]string, 0, len(m))

	for name := range m {
		names = append(names, name)
	}

	sort.Strings(names)
	fields := make([]field, 0, len(names))

	for _, name := range names {
		f, err := schemaField(name, m[name])

		if err != nil {
			return nil, fmt.Errorf("generatorjsz: %s.%s: %w", section, name, err)
		}

		fields = append(fields, f)
	}

	return fields, nil
}

// schemaField converte um Schema folha em field. Object e Array aninhados ainda não são
// implementados.
func schemaField(name string, s *schemaz.Schema) (field, error) {
	f := field{name: name}

	if s.Desc != nil {
		f.summary = s.Desc.Summary
	}

	switch {
	case s.Raw != nil:
		f.typ = jsType(reflect.TypeOf(s.Raw()), 0)
	case s.Str != nil:
		_, native := s.Str().Format()
		f.typ, f.str = jsType(reflect.TypeOf(native), 0), true
	case s.Strs != nil:
		_, native := s.Strs().Format()
		f.typ, f.str = jsType(reflect.TypeOf(native), 0), true
	default:
		return f, fmt.Errorf("not implemented: only Raw, Str and Strs")
	}

	return f, nil
}

var (
	jsonMarshalerType = reflect.TypeFor[json.Marshaler]()
	textMarshalerType = reflect.TypeFor[encoding.TextMarshaler]()
)

func implements(t, iface reflect.Type) bool {
	return t.Implements(iface) || reflect.PointerTo(t).Implements(iface)
}

// jsType é o único uso de reflection: mapeia um tipo Go para o tipo JSDoc. "*" é o desconhecido.
func jsType(t reflect.Type, depth int) string {
	if t == nil || depth > maxDepth {
		return "*"
	}

	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Pointer:
		return jsType(t.Elem(), depth+1)
	case reflect.Slice, reflect.Array:
		return jsType(t.Elem(), depth+1) + "[]"
	case reflect.Map:
		return "Record<string, " + jsType(t.Elem(), depth+1) + ">"
	case reflect.Struct:
		// time.Time, uuid.UUID etc. não serializam pelos campos Go.
		if implements(t, jsonMarshalerType) {
			return "*"
		}

		if implements(t, textMarshalerType) {
			return "string"
		}

		return structType(t, depth)
	}

	return "*"
}

func structType(t reflect.Type, depth int) string {
	var fields []field

	for f := range t.Fields() {
		name := f.Name
		tag, _, _ := strings.Cut(f.Tag.Get("json"), ",")

		if f.PkgPath != "" || tag == "-" {
			continue
		}

		if tag != "" {
			name = tag
		}

		fields = append(fields, field{name: name, typ: jsType(f.Type, depth+1)})
	}

	sort.Slice(fields, func(i, j int) bool { return fields[i].name < fields[j].name })

	return objectType(fields)
}

// objectType renderiza "{a: T, b: U}", ou "Object" quando não há campos.
func objectType(fields []field) string {
	if len(fields) == 0 {
		return "Object"
	}

	parts := make([]string, 0, len(fields))

	for _, f := range fields {
		parts = append(parts, f.name+": "+f.typ)
	}

	return "{" + strings.Join(parts, ", ") + "}"
}
