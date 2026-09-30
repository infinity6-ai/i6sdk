package jsparser

import (
	"fmt"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseSchemaArray(s *schemaz.Array) (*jstypez.Typedef, *jstypez.Typedefs) {
	element := s.Get(0, false)
	eTypedef, eTypedefs := ParseSchema(element)
	// eTypedef.Head.Name = element
	typedef := &jstypez.Typedef{
		Head: &jstypez.Tuple{Typ: fmt.Sprintf("Array<%s>", eTypedef.Head.Typ)},
	}
	typedefs := eTypedefs
	typedefs.Add(eTypedef)
	return typedef, typedefs
}
