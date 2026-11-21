package jsschema

import (
	"context"
	"strings"

	"github.com/infinity6-ai/gox/schemaz/jsschemaz"
	"github.com/infinity6-ai/gox/schemaz/jstypez"
)

func WriteModuleTypedefs(ctx context.Context, typedefs jstypez.Typedefs) string {
	var b strings.Builder
	b.WriteString("/**\n")
	b.WriteString(" * @typedef {import(\"../services/services.js\").Services} Services\n")
	b.WriteString(" */\n")
	b.WriteString("\n")
	b.WriteString("/**\n")
	b.WriteString(" * @typedef {import(\"../invoker/invoker.js\").ServiceOptions} ServiceOptions\n")
	b.WriteString(" */\n")
	b.WriteString("\n")
	b.WriteString("/**\n")
	b.WriteString(" * @template T\n")
	b.WriteString(" * @typedef {T} i6JsonEmbedded\n")
	b.WriteString(" */\n")
	b.WriteString("\n")

	if tdfs := jsschemaz.FormatTypedefs(&typedefs); tdfs != "" {
		b.WriteString(strings.TrimSuffix(tdfs, "\n"))
		b.WriteString("\n")
	}

	return b.String()
}

func WriteTypedefs(ctx context.Context, typedefs jstypez.Typedefs) string {
	return WriteModuleTypedefs(ctx, typedefs)
}
