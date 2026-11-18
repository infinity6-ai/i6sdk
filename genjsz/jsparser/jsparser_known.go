package jsparser

import (
	"context"
	"reflect"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/infinity6-ai/gox/fsz/fsz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func knownTypes(ctx context.Context, baseId string, val reflect.Value) (*jstypez.Typedef, *jstypez.Typedefs) {
	t := reflect.TypeOf(val.Interface()).String()
	switch v := val.Interface().(type) {
	case fsz.FileStat:
		return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: jstypez.TypFileStat, GoTyp: t}}, nil
	case time.Time:
		return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: jstypez.TypTime, GoTyp: t}}, nil
	case bigquery.NullString:
		return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: jstypez.TypString, GoTyp: t}}, nil
	case bigquery.NullFloat64, bigquery.NullInt64:
		return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: jstypez.TypNumber, GoTyp: t}}, nil
	case schemaz.Schema:
		a, b := ParseSchema(ctx, baseId, &v)
		a.Head.Name = ""
		return a, b
	default:
		return nil, nil
	}
}
