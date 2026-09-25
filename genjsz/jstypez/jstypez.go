package jstypez

type Property struct {
	Typ  string
	Name string
	Desc string
}

type Typedef struct {
	Typ        string
	Name       string
	Properties []Property
}

type JSTypeSchema interface {
	JSTypedef() Typedef
}
