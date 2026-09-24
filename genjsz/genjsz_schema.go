package genjsz

import (
	"fmt"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstyp"
)

func ParseSchema(schema *schemaz.Schema) jstyp.Schema {
	if schema.Object != nil {
		ParseSchemaObject(schema.Object(false))
	}
	if schema.Str != nil {
		return ParseSchemaStr(schema.Str())
	}
	panic(fmt.Sprintf("unsupported: %#v", schema))
}
