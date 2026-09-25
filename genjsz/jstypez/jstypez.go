package jstypez

import (
	"fmt"
	"maps"
	"reflect"
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
	Head       TupleV2
	Properties TuplesV2
}

type TypedefsV2 struct {
	m map[string]TypedefV2
}

func (t TypedefsV2) All() []TypedefV2 {
	return slices.Collect(maps.Values(t.m))
}

func (t *TypedefsV2) Add(typedefs ...TypedefV2) {
	// Prevent panics if the method is called on a nil pointer
	if t == nil {
		panic("cannot call Add on a nil *TypedefsV2")
	}

	// Lazy initialization of the map
	if t.m == nil {
		t.m = make(map[string]TypedefV2)
	}

	for _, n := range typedefs {
		// Using Head.Name as the unique key for the typedef
		old, ok := t.m[n.Head.Name]
		if ok {
			// TypedefV2 contains uncomparable fields (map and slice inside TuplesV2),
			// so we must use reflect.DeepEqual to check if they are identical.
			if !reflect.DeepEqual(old, n) {
				panic(fmt.Sprintf("collision, old: %#v, new: %#v", old, n))
			}
			// If it already exists and is completely identical, skip adding
			continue
		}

		t.m[n.Head.Name] = n
	}
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
