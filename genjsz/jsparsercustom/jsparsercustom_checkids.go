package jsparsercustom

import (
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
	var a jstypez.Typedef
	var b jstypez.Typedefs
	return a, b, false
}
