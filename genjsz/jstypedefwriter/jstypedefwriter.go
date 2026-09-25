package jstypedefwriter

import (
	"fmt"
	"regexp"
	"strings"

	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

var (
	boolRe   = regexp.MustCompile(`\bbool\b`)
	objectRe = regexp.MustCompile(`\bobject\b`)
)

func mapType(typ string) string {
	typ = boolRe.ReplaceAllString(typ, "boolean")
	typ = objectRe.ReplaceAllString(typ, "Object")
	return typ
}

func WritePipedefs(pipedefs jstypez.Typedefs) string {
	var b strings.Builder
	for _, pipedef := range pipedefs.All() {
		b.WriteString("/**\n")
		if pipedef.Head.Summary != "" {
			fmt.Fprintf(&b, " * %s\n", pipedef.Head.Summary)
		}
		
		typ := mapType(pipedef.Head.Typ)

		fmt.Fprintf(&b, " * @typedef {%s} %s\n", typ, pipedef.Head.Name)

		if pipedef.Properties != nil {
			for _, prop := range pipedef.Properties.All() {
				pTyp := mapType(prop.Typ)
				if prop.Summary != "" {
					fmt.Fprintf(&b, " * @property {%s} %s - %s\n", pTyp, prop.Name, prop.Summary)
				} else {
					fmt.Fprintf(&b, " * @property {%s} %s\n", pTyp, prop.Name)
				}
			}
		}
		b.WriteString(" */\n\n")
	}
	return b.String()
}
