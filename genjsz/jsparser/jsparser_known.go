package jsparser

import (
	"reflect"
	"time"

	"github.com/infinity6-ai/gox/fsz/fsz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func knownTypes(val reflect.Value) (*jstypez.Typedef, *jstypez.Typedefs) {
	switch val.Interface().(type) {
	case fsz.FileStat:
		return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: jstypez.TypFileStat}}, nil
	case time.Time:
		return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: jstypez.TypTime}}, nil
	default:
		return nil, nil
	}
}
