package jstyp

import (
	"fmt"
)

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
	Schema
	Name string
	Typ  Typ
}

func (f Field) String() string {
	return fmt.Sprintf("%s:%s", f.Name, f.Typ)
}

type Schema interface {
	noop()
}

type Object struct {
	Schema
	Fields map[string]Schema
}

func (f Object) String() string {
	return fmt.Sprintf("%s", f.Fields)
}
