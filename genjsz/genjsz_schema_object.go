package genjsz

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstyp"
)

func ParseSchemaObject(obj map[string]*schemaz.Schema) jstyp.Schema {
	ret := jstyp.Object{
		Fields: make(map[string]jstyp.Schema, len(obj)),
	}
	for field, schema := range obj {
		s := ParseSchema(schema)
		ret.Fields[field] = s
	}
	return ret
}
