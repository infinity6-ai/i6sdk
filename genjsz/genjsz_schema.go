package genjsz

import (
	"fmt"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

func ParseSchema(schema *schemaz.Schema) {
	if schema.Object != nil {
		ParseSchemaObject(schema.Object(false))
	}
	if schema.Str != nil {
		ParseSchemaStr(schema.Str())
	}
	panic(fmt.Sprintf("unsupported: %#v", schema))
}
