package jsparser

import (
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseSchemaRaw(s any) (jstypez.Typedef, jstypez.TypedefsV2) {
	return ParseAny(s)
}
