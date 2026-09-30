package jsparser

import (
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseSchemaRaw(baseId string, s any) (*jstypez.Typedef, *jstypez.Typedefs) {
	return ParseAny(baseId, s)
}
