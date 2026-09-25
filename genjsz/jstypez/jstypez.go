package jstypez

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
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

func (p *TuplesV2) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.All())
}

func (p *TuplesV2) UnmarshalJSON(data []byte) error {
	var s []TupleV2
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("cannot unmarshal path from json: %w", err)
	}
	p.Add(s...)
	return nil
}

func (t TuplesV2) All() []TupleV2 {
	ret := make([]TupleV2, len(t.s))
	for idx, name := range t.s {
		ret[idx] = t.m[name]
	}
	return ret
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

type Typedef struct {
	Head       TupleV2  `json:"head"`
	Properties TuplesV2 `json:"properties"`
}

type TypedefsV2 struct {
	m map[string]Typedef
}

func (p *TypedefsV2) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.m)
}

func (p *TypedefsV2) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &p.m); err != nil {
		return fmt.Errorf("cannot unmarshal path from json: %w", err)
	}
	return nil
}

func (t TypedefsV2) All() []Typedef {
	ret := slices.Collect(maps.Values(t.m))
	slices.SortFunc(ret, func(a, b Typedef) int {
		return strings.Compare(a.Head.Name, b.Head.Name)
	})
	return ret
}

func (t *TypedefsV2) Add(typedefs ...Typedef) {
	// Prevent panics if the method is called on a nil pointer
	if t == nil {
		panic("cannot call Add on a nil *TypedefsV2")
	}

	// Lazy initialization of the map
	if t.m == nil {
		t.m = make(map[string]Typedef)
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

const (
	TypeString = "string"
	TypeNumber = "number"
	TypeBool   = "bool"
	TypeObject = "object"
)

//////

// type Property struct {
// 	Typ  string
// 	Name string
// 	Desc string
// }

// type Typedef struct {
// 	Typ        string
// 	Name       string
// 	Properties []Property
// }

// func (t *Typedef) SetProperty(p Property) {
// 	exists := slices.ContainsFunc(t.Properties, func(existingProperty Property) bool {
// 		return existingProperty.Name == p.Name
// 	})
// 	if exists {
// 		panic(fmt.Sprintf("property collision: %#v", p))
// 	}
// 	t.Properties = append(t.Properties, p)
// }

// var TypeString = Typedef{Name: "string"}
// var TypeNumber = Typedef{Name: "number"}
// var TypeBool = Typedef{Name: "bool"}
