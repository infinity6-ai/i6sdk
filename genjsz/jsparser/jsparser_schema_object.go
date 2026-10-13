package jsparser

import (
	"maps"
	"slices"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/camelstring"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseSchemaObject(baseId string, obj map[string]*schemaz.Schema) (*jstypez.Typedef, *jstypez.Typedefs) {
	ret := jstypez.Typedef{
		Head:       &jstypez.Tuple{Typ: jstypez.TypObject},
		Properties: &jstypez.Tuples{},
	}
	tdfs := &jstypez.Typedefs{}
	for _, field := range slices.Sorted(maps.Keys(obj)) {
		schema := obj[field]
		fieldId := baseId + camelstring.ToCamelUpper(field)
		s, schemaTdfs := ParseSchema(fieldId, schema)
		typ := s.Head.Typ
		if s.Head.Typ == jstypez.TypObject && s.Properties != nil {
			tdfs.Add(s)
			typ = s.Head.Name
		}
		ret.Properties.Add(jstypez.Tuple{
			Name:    field,
			Typ:     typ,
			GoTyp:   s.Head.GoTyp,
			Summary: DescSummary(schema.Desc),
		})
		tdfs.ExtractFrom(schemaTdfs)
	}
	return &ret, tdfs
}
