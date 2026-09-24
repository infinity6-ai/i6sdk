package genjsz

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstyp"
)

func ParseSchemaStr(parser *schemaz.Parser[string]) jstyp.Field {
	s, v := parser.Format()
	typ := ParseAny(v)
	return jstyp.Field{
		Name: s,
		Typ:  typ,
	}
}
