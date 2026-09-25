package jsparsercustom

import (
	"reflect"

	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func IsCheckId(v any) (jstypez.Typedef, jstypez.Typedefs, bool) {
	// if _, ok := v.(lakedsz.Dataset); ok {
	// 	return jstypez.Typedef{Head: jstypez.Tuple{Typ: jstypez.TypeString}}, true
	// }

	ret := func(a jstypez.Typedef, b jstypez.Typedefs) (jstypez.Typedef, jstypez.Typedefs, bool) {
		return a, b, true
	}

	if v, ok := v.(jstypez.CustomType); ok {
		return ret(v.JSTypedef())
	}
	if v, ok := any(v).(jstypez.CustomType); ok {
		return ret(v.JSTypedef())
	}

	rv := reflect.ValueOf(v)
	if rv.IsValid() && rv.Kind() != reflect.Ptr {
		// reflect.New creates a pointer to the type (e.g., *lakedsz.Dataset)
		ptr := reflect.New(rv.Type())

		// Set the copied value into the new pointer's address
		ptr.Elem().Set(rv)

		// Check if this newly created pointer implements your interface
		if ct, ok := ptr.Interface().(jstypez.CustomType); ok {
			return ret(ct.JSTypedef())
		}
	}

	var a jstypez.Typedef
	var b jstypez.Typedefs
	return a, b, false
}
