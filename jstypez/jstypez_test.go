package jstypez_test

import (
	"testing"

	"go.code.infinity6.ai/sdkgo/jstypez"
)

type Labels map[string]string

func (Labels) JSTypedef() jstypez.Typedef {
	return jstypez.Typedef{Typ: "Record<string, string>", Name: "Person"}
}

type Person struct {
	Name   string
	Age    int
	Labels Labels
}

func (Person) JSTypedef() jstypez.Typedef {
	return jstypez.Typedef{
		Typ:  "Object",
		Name: "Person",
		Properties: []jstypez.Property{
			{Typ: "string", Name: "dataset", Desc: "Dataset name"},
			{Typ: "number", Name: "amount", Desc: "amount"},
			{Typ: "Record<string, string>}", Name: "partitions", Desc: "paritions"},
		},
	}
}

func TestUnitSample(t *testing.T) {

}
