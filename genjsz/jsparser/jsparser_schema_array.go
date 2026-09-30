package jsparser

import (
	"fmt"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseSchemaArray(baseId string, s *schemaz.Array) (*jstypez.Typedef, *jstypez.Typedefs) {
	element := s.Get(0, false)
	eTypedef, eTypedefs := ParseSchema(baseId+"Element", element)
	// eTypedef.Head.Name = element
	typedef := &jstypez.Typedef{
		Head: &jstypez.Tuple{Typ: fmt.Sprintf("Array<%s>", eTypedef.Head.Name)},
	}
	typedefs := eTypedefs
	typedefs.Add(eTypedef)
	return typedef, typedefs
}
