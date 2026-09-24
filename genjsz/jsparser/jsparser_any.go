package jsparser

import (
	"fmt"
	"reflect"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/commonz/jsonz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstyp"
)

func ParseAny(v any) jstyp.Typ {
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

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return jstyp.TypNumber

	case reflect.String:
		return jstyp.TypString

	case reflect.Bool:
		return jstyp.TypBool

	// case reflect.Struct:
	// 	return jstyp.TypJson

	default:
		ret := isJsonable(v)
		if ret != nil {
			return *ret
		}
		panic(fmt.Sprintf("unsupported type %T / %s: %#v", v, val.Kind(), v))
	}

}

func isJsonable(v any) *jstyp.Typ {
	b, err := jsonz.Format(v)
	if err != nil {
		return nil
	}
	s := b.String()
	_, err = jsonz.Parse(s, &v)
	if err != nil {
		// return nil
		errorz.Check(err)
	}
	if s[0] == '"' {
		return &jstyp.TypString
	}
	if s[0] >= '0' && s[0] <= '9' {
		return &jstyp.TypNumber
	}
	if s[0] == 'f' || s[0] == 't' {
		return &jstyp.TypBool
	}
	// if s[0] == 'n' || s[0] == '{' {
	panic(fmt.Sprintf("unsupported type %#v", v))
	// }
}
