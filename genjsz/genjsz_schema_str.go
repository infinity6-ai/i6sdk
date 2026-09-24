package genjsz

import (
	"fmt"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

func ParseSchemaStr(parser *schemaz.Parser[string]) {
	s, v := parser.Format()
	println(s)
	typ := ParseAny(v)
	fmt.Printf("%s\n", typ)
}
