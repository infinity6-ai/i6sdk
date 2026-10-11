package dszstore

import (
	"reflect"

	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

type DszStore struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (d *DszStore) JSTypedef() (*jstypez.Typedef, *jstypez.Typedefs) {
	t := reflect.TypeOf(*d).String()
	td := &jstypez.Typedef{
		Head:       &jstypez.Tuple{Typ: jstypez.TypObject, GoTyp: t, Name: ""},
		Properties: &jstypez.Tuples{},
	}
	td.Properties.Add(
		jstypez.Tuple{Typ: jstypez.TypString, GoTyp: reflect.TypeOf(d.Id).String(), Name: "id"},
		jstypez.Tuple{Typ: jstypez.TypString, GoTyp: reflect.TypeOf(d.Name).String(), Name: "name"},
	)
	return td, nil
}
