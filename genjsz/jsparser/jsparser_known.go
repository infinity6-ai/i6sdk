package jsparser

import (
	"reflect"

	"github.com/infinity6-ai/gox/fsz/fsz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func knownTypes(val reflect.Value) (*jstypez.Typedef, *jstypez.Typedefs) {
	switch val.Interface().(type) {
	case fsz.FileStat:
		return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: jstypez.TypFileStat}}, nil
	default:
		return nil, nil
	}
}
