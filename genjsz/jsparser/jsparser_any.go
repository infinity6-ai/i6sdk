package jsparser

import (
	"fmt"
	"reflect"

	"go.code.infinity6.ai/sdkgo/genjsz/jsparsercustom"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseAny(v any) (jstypez.Typedef, jstypez.Typedefs) {
	val := reflect.ValueOf(v)

	for val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	// typ := reflect.TypeOf(v)
	// ptrToV := reflect.New(typ).Interface()
	// x, ok := ptrToV.(json.Unmarshaler)
	// if ok {
	// 	println(x)
	// }

	if !val.IsValid() {
		panic(fmt.Sprintf("invalid value or nil pointer: %#v", v))
	}

	var zeroTdfs jstypez.Typedefs

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return jstypez.Typedef{Head: jstypez.Tuple{Typ: jstypez.TypeNumber}}, zeroTdfs

	case reflect.String:
		return jstypez.Typedef{Head: jstypez.Tuple{Typ: jstypez.TypeString}}, zeroTdfs

	case reflect.Bool:
		return jstypez.Typedef{Head: jstypez.Tuple{Typ: jstypez.TypeBool}}, zeroTdfs

	// case reflect.Struct:
	// 	return jstyp.TypJson

	default:
		if ret, tdfs, ok := jsparsercustom.IsCheckId(v); ok {
			return ret, tdfs
		}
		panic(fmt.Sprintf("unsupported type %T / %s: %#v", v, val.Kind(), v))
	}

}
