package jsparsercustom

import (
	"reflect"

	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func IsCheckId(v any) (jstypez.Typedef, jstypez.Typedefs, bool) {
	ret := func(a jstypez.Typedef, b jstypez.Typedefs) (jstypez.Typedef, jstypez.Typedefs, bool) {
		return a, b, true
	}

	var a jstypez.Typedef
	var b jstypez.Typedefs

	retSimple := func(s string) (jstypez.Typedef, jstypez.Typedefs, bool) {
		return ret(jstypez.Typedef{Head: jstypez.Tuple{Typ: s}}, b)
	}

	if v, ok := v.(jstypez.CustomType); ok {
		return ret(v.JSTypedef())
	}
	if v, ok := v.(jstypez.CustomTypeSimple); ok {
		return retSimple(v.JSType())
	}
	rv := reflect.ValueOf(v)
	if rv.IsValid() && rv.Kind() != reflect.Ptr {
		ptr := reflect.New(rv.Type())
		ptr.Elem().Set(rv)
		if ct, ok := ptr.Interface().(jstypez.CustomType); ok {
			return ret(ct.JSTypedef())
		}
		if ct, ok := ptr.Interface().(jstypez.CustomTypeSimple); ok {
			return retSimple(ct.JSType())
		}
	}

	return a, b, false
}
