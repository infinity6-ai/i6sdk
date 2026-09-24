package genjsz

import (
	"fmt"
	"reflect"

	"go.code.infinity6.ai/sdkgo/genjsz/jstyp"
)

func ParseAny(v any) jstyp.Typ {
	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return jstyp.TypInt

	case reflect.Float32, reflect.Float64:
		return jstyp.TypFloat

	case reflect.String:
		return jstyp.TypString

	case reflect.Bool:
		return jstyp.TypBool

	default:
		panic(fmt.Sprintf("unsupported type %T / %s: %#v", v, val.Kind(), v))
	}
}
