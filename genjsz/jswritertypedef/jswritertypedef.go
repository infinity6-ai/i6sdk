package jswritertypedef

import (
	"context"
	"strings"

	"go.code.infinity6.ai/sdkgo/genjsz/jstypedefwriter"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func Write(ctx context.Context, typedefs jstypez.Typedefs) string {
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

	if tdfs := jstypedefwriter.Write(&typedefs); tdfs != "" {
		b.WriteString(strings.TrimSuffix(tdfs, "\n"))
		b.WriteString("\n")
	}

	return b.String()
}
