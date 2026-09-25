package jsparser

import (
	"fmt"
	"reflect"

	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparsercustom"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseAny(v any) (jstypez.Typedef, jstypez.Typedefs) {
	val := reflect.ValueOf(v)

	for val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface {
		val = val.Elem()
	}

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

	case reflect.Map:
		var allTdfs jstypez.Typedefs

		// 1. Get the reflect.Type of the map
		mapType := val.Type()

		// 2. Generate a zero-value instance of the Key type (e.g., "" for string)
		keyType := mapType.Key()
		keyZeroVal := reflect.Zero(keyType).Interface()

		// 3. Generate a zero-value instance of the Element type (e.g., 0.0 for float64)
		elemType := mapType.Elem()
		elemZeroVal := reflect.Zero(elemType).Interface()

		// 4. Recursively parse the key and value types
		keyTdf, kTdfs := ParseAny(keyZeroVal)
		valTdf, vTdfs := ParseAny(elemZeroVal)

		// 5. Aggregate nested typedefs
		allTdfs.Add(kTdfs.All()...)
		allTdfs.Add(vTdfs.All()...)

		checker.StrEmpty(keyTdf.Head.Name, "keyTdf.Head.Name")
		checker.StrEmpty(valTdf.Head.Name, "valTdf.Head.Name")

		return jstypez.Typedef{
			Head: jstypez.Tuple{Typ: fmt.Sprintf("Record<%s, %s>", keyTdf.Head.Typ, valTdf.Head.Typ)},
		}, allTdfs

	default:
		if ret, tdfs, ok := jsparsercustom.Resolve(v); ok {
			return ret, tdfs
		}
		panic(fmt.Sprintf("unsupported type %T / %s: %#v", v, val.Kind(), v))
	}

}
