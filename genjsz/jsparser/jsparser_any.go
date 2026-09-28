package jsparser

import (
	"fmt"
	"reflect"

	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparsercustom"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func ParseAny(v any) (*jstypez.Typedef, *jstypez.Typedefs) {
	// if a, b := checkCustom(v); a != nil {
	// 	return a, b
	// }

	val := reflect.ValueOf(v)

	for val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface {
		if val.IsNil() {
			val = reflect.Zero(val.Type().Elem())
			break
		}
		val = val.Elem()
	}

	if a, b := checkCustom(val); a != nil {
		return a, b
	}

	if !val.IsValid() {
		panic(fmt.Sprintf("invalid value or nil pointer: %#v", v))
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: jstypez.TypeNumber}}, nil

	case reflect.String:
		return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: jstypez.TypeString}}, nil

	case reflect.Bool:
		return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: jstypez.TypeBool}}, nil

	case reflect.Slice, reflect.Array:
		var allTdfs jstypez.Typedefs

		// 1. Get the reflect.Type of the slice/array element
		elemType := val.Type().Elem()

		// Optional: Special handling for []byte (often mapped to string or Uint8Array in JS)
		// if elemType.Kind() == reflect.Uint8 {
		// 	return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: "string"}}, allTdfs
		// }

		// 2. Generate a zero-value instance of the Element type
		elemZeroVal := reflect.Zero(elemType).Interface()

		// 3. Recursively parse the element type
		elemTdf, eTdfs := ParseAny(elemZeroVal)

		// 4. Aggregate nested typedefs
		allTdfs.ExtractFrom(eTdfs)

		// 5. Format as a TypeScript array.
		// You can use "%s[]" (e.g., "string[]") or "Array<%s>" (e.g., "Array<string>")
		return &jstypez.Typedef{
			Head: &jstypez.Tuple{Typ: fmt.Sprintf("Array<%s>", elemTdf.Head.Typ)},
		}, &allTdfs

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
		allTdfs.ExtractFrom(kTdfs)
		allTdfs.ExtractFrom(vTdfs)

		checker.StrEmpty(keyTdf.Head.Name, "keyTdf.Head.Name")
		checker.StrEmpty(valTdf.Head.Name, "valTdf.Head.Name")

		return &jstypez.Typedef{
			Head: &jstypez.Tuple{Typ: fmt.Sprintf("Record<%s, %s>", keyTdf.Head.Typ, valTdf.Head.Typ)},
		}, &allTdfs

	default:
		panic(fmt.Sprintf("unsupported type %T / %s: %#v", v, val.Kind(), v))
	}

}

func checkCustom(val reflect.Value) (*jstypez.Typedef, *jstypez.Typedefs) {
	// if val.IsNil() {
	// 	val = reflect.Zero(val.Type().Elem())
	// }
	v := val.Interface()
	if ret, tdfs := jsparsercustom.Resolve(v); ret != nil {
		if tdfs != nil {
			checker.Empty(tdfs.All(), "must be empty")
		}
		return ret, tdfs
	}
	return nil, nil
}
