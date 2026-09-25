package jstypez

import (
	"fmt"
	"slices"
)

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

func (t *Typedef) SetProperty(p Property) {
	exists := slices.ContainsFunc(t.Properties, func(existingProperty Property) bool {
		return existingProperty.Name == p.Name
	})
	if exists {
		panic(fmt.Sprintf("property collision: %#v", p))
	}
	t.Properties = append(t.Properties, p)
}

var TypeString = Typedef{Name: "string"}
var TypeNumber = Typedef{Name: "number"}
var TypeBool = Typedef{Name: "bool"}
