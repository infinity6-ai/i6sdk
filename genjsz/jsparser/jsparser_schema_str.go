package jsparser

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseSchemaStr(parser *schemaz.Parser[string]) jstypez.Typedef {
	_, v := parser.Format()
	typ := ParseAny(v)
	return typ
}
