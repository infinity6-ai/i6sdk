package jstyp

import "go.code.infinity6.ai/sdkgo/genjsz/jstyp"

type Typ struct {
	Js string
}

func (t Typ) String() string {
	return t.Js
}

var TypInt = Typ{
	Js: "int",
}

var TypFloat = Typ{
	Js: "float",
}

var TypString = Typ{
	Js: "string",
}

var TypBool = Typ{
	Js: "bool",
}

type Field struct {
	Name string
	Typ  jstyp.Typ
}
