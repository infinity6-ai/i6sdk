package genjsz

import (
	"fmt"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

func ParseSchema(schema *schemaz.Schema) {
	if schema.Object == nil {
		ParseSchemaObject(schema.Object(false))
	}
	panic(fmt.Sprintf("unsupported: %#v", schema))
}

func ParseSchemaObject(obj map[string]*schemaz.Schema) {
	for field, schema := range obj {
		println(field)
		ParseSchema(schema)
	}
}
