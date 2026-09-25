package jstypez

import (
	"fmt"
	"slices"
)

type TupleV2 struct {
	Typ     string `json:"typ"`
	Name    string `json:"name"`
	Summary string `json:"summary"`
}

type TuplesV2 struct {
	m map[string]TupleV2
	s []string
}

func (t *TuplesV2) Add(tuples ...TupleV2) {
	// Prevent panics if the method is called on a nil pointer
	if t == nil {
		panic("cannot call Add on a nil *TuplesV2")
	}

	// Lazy initialization of the map
	if t.m == nil {
		t.m = make(map[string]TupleV2)
	}

	for _, n := range tuples {
		old, ok := t.m[n.Name]
		if ok {
			if old != n {
				panic(fmt.Sprintf("collision, old: %#v, new: %#v", old, n))
			}
			// If it already exists and is identical, skip appending to avoid slice duplicates
			continue
		}

		t.m[n.Name] = n
		t.s = append(t.s, n.Name)
	}
}

type TypedefV2 struct {
	Typ        string
	Name       string
	Properties TuplesV2
}

//////

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
