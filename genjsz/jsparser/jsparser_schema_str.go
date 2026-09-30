package jsparser

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseSchemaStr(baseId string, parser *schemaz.Parser[string]) (*jstypez.Typedef, *jstypez.Typedefs) {
	_, v := parser.Format()
	return ParseAny(baseId, v)
}
