package genjsz

import "go.code.infinity6.ai/sdkgo/genjsz/jstyp"

func ParseSchemaRaw(s any) jstyp.Typ {
	return ParseAny(s)
}
