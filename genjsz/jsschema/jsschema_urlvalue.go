package jsschema

import (
	"context"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func (p *Parser) ParseSchemaUrlValue(ctx context.Context, baseId string, parser *schemaz.Parser[string]) (*jstypez.Typedef, *jstypez.Typedefs) {
	_, v := parser.Format()
	return p.ParseAny(ctx, baseId, v)
}
