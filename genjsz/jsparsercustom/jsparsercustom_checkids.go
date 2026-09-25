package jsparsercustom

import (
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

func IsCheckId(v any) (jstypez.Typedef, bool) {
	if _, ok := v.(lakedsz.Dataset); ok {
		return jstypez.TypeString, true
	}
	var zero jstypez.Typedef
	return zero, false
}
