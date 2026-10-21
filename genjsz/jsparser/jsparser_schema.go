package jsparser

import (
	"context"
	"fmt"

	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseSchema(ctx context.Context, baseId string, schema *schemaz.Schema) (*jstypez.Typedef, *jstypez.Typedefs) {
	var typedef *jstypez.Typedef
	var typedefs *jstypez.Typedefs
	if schema.Object != nil {
		s := schema.Object(false)
		typedef, typedefs = ParseSchemaObject(ctx, baseId, s)
	}
	if schema.UrlValue != nil {
		s := schema.UrlValue()
		typedef, typedefs = ParseSchemaUrlValue(ctx, baseId, s)
	}
	if schema.Field != nil {
		s := schema.Field()
		typedef, typedefs = ParseSchemaField(ctx, baseId, s)
	}
	if schema.Array != nil {
		s := schema.Array()
		typedef, typedefs = ParseSchemaArray(ctx, baseId, s)
	}
	if typedef == nil {
		panic(fmt.Sprintf("unsupported: %#v", schema))
	}
	checker.StrEmpty(typedef.Head.Name, "typedef.Head.Name")
	typedef.Head.Name = baseId
	typedef.Head.Summary = DescSummary(schema.Desc)
	return typedef, typedefs
}
