package jsparser

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseSchemaObject(obj map[string]*schemaz.Schema) jstypez.Typedef {
	ret := jstypez.Typedef{
		Typ: "object",
	}
	for field, schema := range obj {
		s := ParseSchema(schema)
		ret.SetProperty(jstypez.Property{
			Name: field,
			Typ:  s.Name,
			Desc: descSummary(schema.Desc),
		})
	}
	return ret
}
