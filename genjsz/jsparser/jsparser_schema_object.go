package jsparser

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseSchemaObject(obj map[string]*schemaz.Schema) (jstypez.Typedef, jstypez.TypedefsV2) {
	ret := jstypez.Typedef{
		Typ: "object",
	}
	var tdfs jstypez.TypedefsV2
	for field, schema := range obj {
		s, schemaTdfs := ParseSchema(schema)
		ret.SetProperty(jstypez.Property{
			Name: field,
			Typ:  s.Name,
			Desc: descSummary(schema.Desc),
		})
		tdfs.Add(schemaTdfs.All()...)
	}
	return ret, tdfs
}
