package jsparser

import (
	"fmt"
	"reflect"

	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

func ParseAny(v any) (jstypez.Typedef, jstypez.TypedefsV2) {
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

	var tdfs jstypez.TypedefsV2

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return jstypez.Typedef{Head: jstypez.TupleV2{Typ: jstypez.TypeNumber}}, tdfs

	case reflect.String:
		return jstypez.Typedef{Head: jstypez.TupleV2{Typ: jstypez.TypeString}}, tdfs

	case reflect.Bool:
		return jstypez.Typedef{Head: jstypez.TupleV2{Typ: jstypez.TypeBool}}, tdfs

	// case reflect.Struct:
	// 	return jstyp.TypJson

	default:
		// if ret, ok := jsparsercustom.IsCheckId(v); ok {
		// 	return ret, tdfs
		// }
		panic(fmt.Sprintf("unsupported type %T / %s: %#v", v, val.Kind(), v))
	}

}

func bla(v any, val reflect.Value) {
	var targetType = reflect.TypeOf((*lakedsz.Dataset)(nil)).Elem()
	if val.Type().ConvertibleTo(targetType) {
		convertedValue := val.Convert(targetType)
		dataset := convertedValue.Interface().(lakedsz.Dataset)
		fmt.Printf("%s\n", dataset)

		dataset, ok := v.(lakedsz.Dataset)
		fmt.Printf("%s %t\n", dataset, ok)
	}
	// vType.ConvertibleTo(targetType)
	// val.CanConvert()
}

// func isJsonable(v any) *jstyp.Typ {
// 	b, err := jsonz.Format(v)
// 	if err != nil {
// 		return nil
// 	}
// 	s := b.String()
// 	_, err = jsonz.Parse(s, &v)
// 	if err != nil {
// 		errorz.Check(err)
// 	}
// 	if s[0] == '"' {
// 		return &jstyp.TypString
// 	}
// 	if s[0] >= '0' && s[0] <= '9' {
// 		return &jstyp.TypNumber
// 	}
// 	if s[0] == 'f' || s[0] == 't' {
// 		return &jstyp.TypBool
// 	}
// 	return nil
// }
