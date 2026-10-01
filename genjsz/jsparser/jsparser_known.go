package jsparser

import (
	"reflect"

	"github.com/infinity6-ai/gox/fsz/fsz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func knownTypes(baseId string, val reflect.Value) (*jstypez.Typedef, *jstypez.Typedefs) {
	switch val.Interface().(type) {
	case fsz.FileStat:
		return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: TypFileStat.Head.Name}}, nil
	default:
		return nil, nil
	}
}
