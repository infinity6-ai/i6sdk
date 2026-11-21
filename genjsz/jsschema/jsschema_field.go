package jsschema

import (
	"context"

	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseSchemaField(ctx context.Context, baseId string, s any) (*jstypez.Typedef, *jstypez.Typedefs) {
	return ParseAny(ctx, baseId, s)
}
