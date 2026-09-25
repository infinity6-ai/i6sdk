package jsparsercustom

import (
	"reflect"

	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func Resolve(v any) (*jstypez.Typedef, *jstypez.Typedefs) {
	if v, ok := v.(jstypez.CustomType); ok {
		return v.JSTypedef()
	}
	if v, ok := v.(jstypez.CustomTypeSimple); ok {
		return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: v.JSType()}}, nil
	}
	rv := reflect.ValueOf(v)
	if rv.IsValid() && rv.Kind() != reflect.Ptr {
		ptr := reflect.New(rv.Type())
		ptr.Elem().Set(rv)
		if ct, ok := ptr.Interface().(jstypez.CustomType); ok {
			return ct.JSTypedef()
		}
		if ct, ok := ptr.Interface().(jstypez.CustomTypeSimple); ok {
			return &jstypez.Typedef{Head: &jstypez.Tuple{Typ: ct.JSType()}}, nil
		}
	}

	return nil, nil
}
