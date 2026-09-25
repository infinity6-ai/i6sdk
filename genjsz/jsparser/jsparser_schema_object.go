package jsparser

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseSchemaObject(obj map[string]*schemaz.Schema) (jstypez.Typedef, jstypez.TypedefsV2) {
	ret := jstypez.Typedef{
		Head: jstypez.TupleV2{Typ: jstypez.TypeObject},
	}
	var tdfs jstypez.TypedefsV2
	for field, schema := range obj {
		s, schemaTdfs := ParseSchema(schema)
		ret.Properties.Add(jstypez.TupleV2{
			Name:    field,
			Typ:     s.Head.Typ,
			Summary: descSummary(schema.Desc),
		})
		tdfs.Add(schemaTdfs.All()...)
	}
	return ret, tdfs
}
