package jsschema

import (
	"context"

	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func (p *Parser) ParseSchemaField(ctx context.Context, baseId string, s any) (*jstypez.Typedef, *jstypez.Typedefs) {
	return p.ParseAny(ctx, baseId, s)
}
