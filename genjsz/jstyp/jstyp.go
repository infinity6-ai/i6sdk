package jstyp

import (
	"fmt"
	"strings"
)

type Typ struct {
	Schema
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
	Schema
	Name string
	Typ  Typ
}

func (f Field) String() string {
	return fmt.Sprintf("%s:%s", f.Name, f.Typ)
}

type Schema interface {
	noop()
	String() string
}

type Object struct {
	Schema
	Fields map[string]Schema
}

func (f Object) String() string {
	var ret strings.Builder
	ret.WriteString("{")
	c := 0
	for k, v := range f.Fields {
		ret.WriteString(k)
		ret.WriteString(":")
		ret.WriteString(v.String())
		if c < len(f.Fields)-1 {
			ret.WriteString(",")
		}
		c++
	}
	ret.WriteString("}")
	return ret.String()
}
