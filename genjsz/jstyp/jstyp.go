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

var TypNumber = Typ{
	Js: "number",
}

var TypString = Typ{
	Js: "string",
}

var TypBool = Typ{
	Js: "bool",
}

var TypJson = Typ{
	Js: "json",
}
