package jsparser

import (
	"reflect"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/infinity6-ai/gox/fsz/fsz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func knownTypes(val reflect.Value) (*jstypez.Typedef, *jstypez.Typedefs) {
	t := reflect.TypeOf(val.Interface()).String()
	switch val.Interface().(type) {
	case fsz.FileStat:
		return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: jstypez.TypFileStat, GoTyp: t}}, nil
	case time.Time:
		return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: jstypez.TypTime, GoTyp: t}}, nil
	case bigquery.NullString:
		return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: jstypez.TypString, GoTyp: t}}, nil
	default:
		return nil, nil
	}
}
