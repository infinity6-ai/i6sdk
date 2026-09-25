package jsparser

import (
	"maps"
	"slices"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseSchemaObject(obj map[string]*schemaz.Schema) (*jstypez.Typedef, *jstypez.Typedefs) {
	ret := jstypez.Typedef{
		Head:       &jstypez.Tuple{Typ: jstypez.TypeObject},
		Properties: &jstypez.Tuples{},
	}
	tdfs := &jstypez.Typedefs{}
	for _, field := range slices.Sorted(maps.Keys(obj)) {
		schema := obj[field]
		s, schemaTdfs := ParseSchema(schema)
		ret.Properties.Add(jstypez.Tuple{
			Name:    field,
			Typ:     s.Head.Typ,
			Summary: descSummary(schema.Desc),
		})
		tdfs.ExtractFrom(schemaTdfs)
	}
	return &ret, tdfs
}
