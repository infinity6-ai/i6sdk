package jstyp

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
