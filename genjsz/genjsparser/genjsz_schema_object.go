package genjsparser

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstyp"
)

func ParseSchemaObject(obj map[string]*schemaz.Schema) jstyp.Typ {
	ret := jstyp.Typ{
		Js:     "object",
		Fields: make(map[string]jstyp.Typ, len(obj)),
	}
	for field, schema := range obj {
		s := ParseSchema(schema)
		ret.Fields[field] = s
	}
	return ret
}
