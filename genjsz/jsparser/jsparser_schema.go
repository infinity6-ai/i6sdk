package jsparser

import (
	"fmt"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseSchema(schema *schemaz.Schema) (*jstypez.Typedef, *jstypez.Typedefs) {
	var typedef *jstypez.Typedef
	var typedefs *jstypez.Typedefs
	if schema.Object != nil {
		s := schema.Object(false)
		typedef, typedefs = ParseSchemaObject(s)
	}
	if schema.Str != nil {
		s := schema.Str()
		typedef, typedefs = ParseSchemaStr(s)
	}
	if schema.Raw != nil {
		s := schema.Raw()
		typedef, typedefs = ParseSchemaRaw(s)
	}
	if typedef == nil {
		panic(fmt.Sprintf("unsupported: %#v", schema))
	}
	return typedef, typedefs
}
