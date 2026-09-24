package genjsz

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstyp"
)

func ParseSchemaObject(obj map[string]*schemaz.Schema) jstyp.Typ {
	ret := jstyp.Object{
		Fields: make(map[string]jstyp.Typ, len(obj)),
	}
	for field, schema := range obj {
		println(field)
		ParseSchema(schema)
	}
}
