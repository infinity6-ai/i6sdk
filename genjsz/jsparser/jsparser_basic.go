package jsparser

import (
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

var TypTime = jstypez.Typedef{Head: &jstypez.Tuple{Typ: "string", Name: "dszTime", Summary: "Missing Summary"}}

var TypFileStat = jstypez.Typedef{
	Head: &jstypez.Tuple{Typ: "object", Name: "dszFileStat"},
	Properties: jstypez.NewTuples([]jstypez.Tuple{
		{Typ: "string", Name: "url", Summary: "missing summary"},
		{Typ: "string", Name: "content_type", Summary: "missing summary"},
		{Typ: "string", Name: "md5", Summary: "missing summary"},
		{Typ: "string", Name: "size", Summary: "missing summary"},
		{Typ: "string", Name: "etag", Summary: "missing summary"},
		{Typ: "dszTime", Name: "created_at", Summary: "missing summary"},
		{Typ: "dszTime", Name: "updated_at", Summary: "missing summary"},
	}),
}

func basicTypes() *jstypez.Typedefs {
	var tds jstypez.Typedefs
	tds.Add([]*jstypez.Typedef{
		&TypTime,
		&TypFileStat,
	}...)
	return &tds
}
