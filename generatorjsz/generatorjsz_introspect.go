package generatorjsz

import (
	"encoding"
	"encoding/json"
	"reflect"
	"sort"
	"strings"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

const maxDepth = 32

// jsonMarshalerType/textMarshalerType detectam tipos com serialização JSON customizada
// (time.Time, uuid.UUID, etc.) antes de tentar decompor os campos Go do struct, que dariam
// um resultado errado (ex.: time.Time só tem campos não exportados — viraria "{}" vazio,
// quando na real ele serializa como uma string).
var (
	jsonMarshalerType = reflect.TypeFor[json.Marshaler]()
	textMarshalerType = reflect.TypeFor[encoding.TextMarshaler]()
)

// implementsMarshaler checa t e *t (o método pode estar em qualquer um dos dois receivers).
func implementsMarshaler(t reflect.Type, iface reflect.Type) bool {
	return t.Implements(iface) || reflect.PointerTo(t).Implements(iface)
}

// Introspect percorre um *schemaz.Schema (de uma instância zero-valor de uma Api) e monta a IR.
func Introspect(schema *schemaz.Schema) *Node {
	return introspectSchema(schema, 0)
}

// introspectSchema decide qual dos 5 bindings (Object/Array/Raw/Str/Strs) está setado e delega.
func introspectSchema(s *schemaz.Schema, depth int) *Node {
	if s == nil {
		return &Node{Kind: KindUnknown}
	}
	var n *Node
	if depth > maxDepth {
		n = &Node{Kind: KindUnknown}
	} else {
		switch {
		case s.Object != nil:
			n = introspectObject(s.Object, depth)
		case s.Array != nil:
			n = introspectArray(s.Array, depth)
		case s.Raw != nil:
			n = introspectRaw(s.Raw, depth)
		case s.Str != nil:
			n = introspectStr(s.Str)
		case s.Strs != nil:
			n = introspectStrs(s.Strs)
		default:
			n = &Node{Kind: KindUnknown}
		}
	}
	n.Desc = s.Desc
	return n
}

// introspectObject chama Object(false) (direção unmarshal, aloca ponteiros aninhados) e ordena as chaves
func introspectObject(fn func(read bool) map[string]*schemaz.Schema, depth int) *Node {
	m := safely(func() map[string]*schemaz.Schema { return fn(false) })

	if m == nil {
		return &Node{Kind: KindObject, Fields: map[string]*Node{}}
	}

	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	fields := make(map[string]*Node, len(m))

	for _, k := range keys {
		fields[k] = introspectSchema(m[k], depth+1)
	}

	return &Node{Kind: KindObject, Fields: fields, FieldOrder: keys}
}

func introspectArray(fn func() *schemaz.Array, depth int) *Node {
	arr := safely(fn)

	if arr == nil {
		return &Node{Kind: KindArray, Element: &Node{Kind: KindUnknown}}
	}

	elemSchema := safely(func() *schemaz.Schema { return arr.Get(0, false) })

	if elemSchema == nil {
		return &Node{Kind: KindArray, Element: &Node{Kind: KindUnknown}}
	}

	return &Node{Kind: KindArray, Element: introspectSchema(elemSchema, depth+1)}
}

func introspectRaw(fn func() any, depth int) *Node {
	v := safely(fn)
	if v == nil {
		return &Node{Kind: KindUnknown}
	}
	return introspectType(reflect.TypeOf(v), map[reflect.Type]bool{}, depth)
}

func introspectStr(fn func() *schemaz.Parser[string]) *Node {
	native := nativeOf(fn)
	if native == nil {
		return &Node{Kind: KindString, StringArrayWire: true}
	}
	n := introspectType(reflect.TypeOf(native), map[reflect.Type]bool{}, 0)
	if n.Kind == KindUnknown {
		n.Kind = KindString
	}
	n.StringArrayWire = true
	return n
}

func introspectStrs(fn func() *schemaz.Parser[[]string]) *Node {
	native := nativeOf(fn)

	if native == nil {
		return &Node{Kind: KindUnknown, MultiValue: true, StringArrayWire: true}
	}

	base := introspectType(reflect.TypeOf(native), map[reflect.Type]bool{}, 0)
	base.MultiValue = true
	base.StringArrayWire = true
	return base
}

// introspectType mapeia um reflect.Type Go para o Kind JS/JSDoc correspondente, recursivamente.
func introspectType(t reflect.Type, visited map[reflect.Type]bool, depth int) *Node {
	if t == nil {
		return &Node{Kind: KindUnknown}
	}
	if depth > maxDepth || visited[t] {
		return &Node{Kind: KindUnknown}
	}
	switch t.Kind() {
	case reflect.String:
		return &Node{Kind: KindString}
	case reflect.Bool:
		return &Node{Kind: KindBoolean}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return &Node{Kind: KindNumber}
	case reflect.Pointer:
		return introspectType(t.Elem(), extendVisited(visited, t), depth+1)
	case reflect.Slice, reflect.Array:
		return &Node{Kind: KindArray, Element: introspectType(t.Elem(), extendVisited(visited, t), depth+1)}
	case reflect.Map:
		return &Node{Kind: KindObject, DynamicKeys: true, Element: introspectType(t.Elem(), extendVisited(visited, t), depth+1)}
	case reflect.Struct:
		if implementsMarshaler(t, jsonMarshalerType) {
			return &Node{Kind: KindUnknown}
		}
		if implementsMarshaler(t, textMarshalerType) {
			return &Node{Kind: KindString}
		}
		return introspectStructType(t, extendVisited(visited, t), depth+1)
	default:
		return &Node{Kind: KindUnknown}
	}
}

func introspectStructType(t reflect.Type, visited map[reflect.Type]bool, depth int) *Node {
	fields := map[string]*Node{}
	var order []string

	for f := range t.Fields() {
		if f.PkgPath != "" {
			continue
		}
		name, skip := jsonFieldName(f)
		if skip {
			continue
		}
		fields[name] = introspectType(f.Type, visited, depth+1)
		order = append(order, name)
	}

	sort.Strings(order)

	return &Node{Kind: KindObject, Fields: fields, FieldOrder: order}
}

func jsonFieldName(f reflect.StructField) (name string, skip bool) {
	tag := f.Tag.Get("json")
	if tag == "-" {
		return "", true
	}
	if tag == "" {
		return f.Name, false
	}

	parts := strings.SplitN(tag, ",", 2)

	if parts[0] == "" {
		return f.Name, false
	}

	return parts[0], false
}

func extendVisited(visited map[reflect.Type]bool, t reflect.Type) map[reflect.Type]bool {
	n := make(map[reflect.Type]bool, len(visited)+1)
	for k, v := range visited {
		n[k] = v
	}
	n[t] = true
	return n
}

func safely[T any](fn func() T) (out T) {
	defer func() {
		if recover() != nil {
			var zero T
			out = zero
		}
	}()
	return fn()
}

// nativeOf chama Str()/Strs().Format() e devolve só o valor nativo (segundo retorno), ou nil se
// o parser não tem Format ou entra em pânico.
func nativeOf[T any](fn func() *schemaz.Parser[T]) any {
	return safely(func() any {
		p := fn()
		if p == nil || p.Format == nil {
			return nil
		}
		_, native := p.Format()
		return native
	})
}
