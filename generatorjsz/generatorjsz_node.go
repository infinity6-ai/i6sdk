package generatorjsz

import "github.com/infinity6-ai/gox/schemaz/schemaz"

// Kind é o tipo JS/JSDoc
type Kind string

const (
	KindString  Kind = "string"
	KindNumber  Kind = "number"
	KindBoolean Kind = "boolean"
	KindObject  Kind = "object"
	KindArray   Kind = "array"
	// KindUnknown não é erro: o campo é tratado como JSON.
	KindUnknown Kind = "unknown"
)

// Node é a representação intermediária (IR) de um nó do schemaz.Schema, produzida pelo walker
// e consumida pelo emissor de JS.
type Node struct {
	Kind Kind
	// Desc é opcional.
	Desc       *schemaz.Desc
	MultiValue bool
	// StringArrayWire é true para nós vindos de Str ou Strs: ambos serializam como array JSON
	// de strings (Str = array de 1 elemento). Só importa para montar/ler ReqBody manualmente;
	StringArrayWire bool
	// DynamicKeys indica Raw -> map[K]V (chaves não fixas).
	DynamicKeys bool
	Fields      map[string]*Node
	FieldOrder  []string
	// Element é o tipo do elemento de um array, ou o tipo de valor quando DynamicKeys.
	Element *Node
}
