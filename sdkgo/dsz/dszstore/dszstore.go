package dszstore

import (
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

type DszStore struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (*DszStore) JSTypedef() (*jstypez.Typedef, *jstypez.Typedefs) {
	td := &jstypez.Typedef{
		Head:       &jstypez.Tuple{Typ: jstypez.TypObject, Name: ""},
		Properties: &jstypez.Tuples{},
	}
	td.Properties.Add(
		jstypez.Tuple{Typ: jstypez.TypString, Name: "id"},
		jstypez.Tuple{Typ: jstypez.TypString, Name: "name"},
	)
	return td, nil
}
