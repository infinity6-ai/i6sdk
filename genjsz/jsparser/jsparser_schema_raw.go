package jsparser

import (
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseSchemaRaw(s any) (jstypez.Typedef, jstypez.Typedefs) {
	return ParseAny(s)
}
