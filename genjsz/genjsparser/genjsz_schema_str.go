package genjsparser

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstyp"
)

func ParseSchemaStr(parser *schemaz.Parser[string]) jstyp.Typ {
	_, v := parser.Format()
	typ := ParseAny(v)
	return typ
}
