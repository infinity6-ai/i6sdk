package jsparser

import (
	"fmt"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseSchema(schema *schemaz.Schema) (jstypez.Typedef, jstypez.TypedefsV2) {
	if schema.Object != nil {
		s := schema.Object(false)
		return ParseSchemaObject(s)
	}
	if schema.Str != nil {
		s := schema.Str()
		return ParseSchemaStr(s)
	}
	if schema.Raw != nil {
		s := schema.Raw()
		return ParseSchemaRaw(s)
	}
	panic(fmt.Sprintf("unsupported: %#v", schema))
}
