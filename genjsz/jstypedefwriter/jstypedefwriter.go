package jstypedefwriter

import (
	"fmt"
	"strings"

	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func Write(pipedefs *jstypez.Typedefs) string {
	var b strings.Builder
	for _, pipedef := range pipedefs.All() {
		b.WriteString("/**\n")
		if pipedef.Head.Summary != "" {
			fmt.Fprintf(&b, " * %s\n", pipedef.Head.Summary)
		}

		fmt.Fprintf(&b, " * @typedef {%s} %s\n", pipedef.Head.Typ, pipedef.Head.Name)

		if pipedef.Properties != nil {
			for _, prop := range pipedef.Properties.All() {
				if prop.Summary != "" {
					fmt.Fprintf(&b, " * @property {%s} %s - %s\n", prop.Typ, prop.Name, prop.Summary)
				} else {
					fmt.Fprintf(&b, " * @property {%s} %s\n", prop.Typ, prop.Name)
				}
			}
		}
		b.WriteString(" */\n\n")
	}
	return b.String()
}
