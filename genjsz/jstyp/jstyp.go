package jstyp

import (
	"github.com/infinity6-ai/gox/commonz/jsonz"
)

type Typ struct {
	Js     string         `json:"js"`
	Fields map[string]Typ `json:"fields,omitempty"`
}

func (t Typ) String() string {
	return jsonz.MustFormat(t).String()
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
