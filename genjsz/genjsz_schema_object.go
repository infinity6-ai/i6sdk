package genjsz

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

func ParseSchemaObject(obj map[string]*schemaz.Schema) {
	for field, schema := range obj {
		println(field)
		ParseSchema(schema)
	}
}
