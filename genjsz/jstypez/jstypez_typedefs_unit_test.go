package jstypez_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func TestUnitTuples(t *testing.T) {
	type testScenario struct {
		name       string
		tuples     []jstypez.Tuple
		jsonInput  string
		wantTuples []jstypez.Tuple
		wantJSON   string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()

		if s.tuples != nil {
			tups := jstypez.NewTuples(s.tuples)
			require.Equal(t, s.wantTuples, tups.All())

			data, err := json.Marshal(tups)
			require.NoError(t, err)
			if s.wantJSON != "" {
				require.JSONEq(t, s.wantJSON, string(data))
			}

			var unmarshaled jstypez.Tuples
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)
			require.Equal(t, s.wantTuples, unmarshaled.All())
		}

		if s.jsonInput != "" {
			var unmarshaled jstypez.Tuples
			err := json.Unmarshal([]byte(s.jsonInput), &unmarshaled)
			require.NoError(t, err)
			require.Equal(t, s.wantTuples, unmarshaled.All())
		}
	}

	t.Run("tuples with GoTyp preserves field through creation and serialization", func(t *testing.T) {
		check(t, testScenario{
			tuples: []jstypez.Tuple{
				{Typ: jstypez.TypString, GoTyp: "string", Name: "id", Summary: "Identifier"},
				{Typ: jstypez.TypNumber, GoTyp: "int64", Name: "count", Summary: "Counter"},
			},
			wantTuples: []jstypez.Tuple{
				{Typ: jstypez.TypString, GoTyp: "string", Name: "id", Summary: "Identifier"},
				{Typ: jstypez.TypNumber, GoTyp: "int64", Name: "count", Summary: "Counter"},
			},
			wantJSON: `[
				{"typ":"string","go_typ":"string","name":"id","summary":"Identifier"},
				{"typ":"number","go_typ":"int64","name":"count","summary":"Counter"}
			]`,
		})
	})

	t.Run("tuples unmarshal from json with go_typ", func(t *testing.T) {
		check(t, testScenario{
			jsonInput: `[
				{"typ":"object","go_typ":"dszFileStat","name":"file","summary":"File stat"}
			]`,
			wantTuples: []jstypez.Tuple{
				{Typ: "object", GoTyp: "dszFileStat", Name: "file", Summary: "File stat"},
			},
		})
	})

	t.Run("tuples skip identical duplicate tuple", func(t *testing.T) {
		tups := jstypez.NewTuples([]jstypez.Tuple{
			{Typ: jstypez.TypString, GoTyp: "string", Name: "foo", Summary: "summary"},
		})
		tups.Add(jstypez.Tuple{Typ: jstypez.TypString, GoTyp: "string", Name: "foo", Summary: "summary"})
		require.Len(t, tups.All(), 1)
	})

	t.Run("tuples panic on collision with different GoTyp", func(t *testing.T) {
		tups := jstypez.NewTuples([]jstypez.Tuple{
			{Typ: jstypez.TypNumber, GoTyp: "int32", Name: "field", Summary: "desc"},
		})
		require.Panics(t, func() {
			tups.Add(jstypez.Tuple{Typ: jstypez.TypNumber, GoTyp: "int64", Name: "field", Summary: "desc"})
		})
	})
}

func TestUnitTypedefs(t *testing.T) {
	t.Run("typedef with GoTyp in Head and Properties", func(t *testing.T) {
		td := &jstypez.Typedef{
			Head: &jstypez.Tuple{
				Typ:     jstypez.TypObject,
				GoTyp:   "MyStruct",
				Name:    "MyStruct",
				Summary: "Sample struct",
			},
			Properties: jstypez.NewTuples([]jstypez.Tuple{
				{Typ: jstypez.TypString, GoTyp: "string", Name: "name", Summary: "Name"},
				{Typ: jstypez.TypNumber, GoTyp: "float64", Name: "ratio", Summary: "Ratio"},
			}),
		}

		tds := jstypez.NewTypedefs([]*jstypez.Typedef{td})
		all := tds.All()
		require.Len(t, all, 1)
		require.Equal(t, "MyStruct", all[0].Head.GoTyp)

		props := all[0].Properties.All()
		require.Len(t, props, 2)
		require.Equal(t, "string", props[0].GoTyp)
		require.Equal(t, "float64", props[1].GoTyp)

		data, err := json.Marshal(&tds)
		require.NoError(t, err)

		var unmarshaled jstypez.Typedefs
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		unmarshaledAll := unmarshaled.All()
		require.Len(t, unmarshaledAll, 1)
		require.Equal(t, "MyStruct", unmarshaledAll[0].Head.GoTyp)
		unmarshaledProps := unmarshaledAll[0].Properties.All()
		require.Len(t, unmarshaledProps, 2)
		require.Equal(t, "string", unmarshaledProps[0].GoTyp)
		require.Equal(t, "float64", unmarshaledProps[1].GoTyp)
	})

	t.Run("typedefs panic on collision with different GoTyp", func(t *testing.T) {
		td1 := &jstypez.Typedef{
			Head: &jstypez.Tuple{Typ: jstypez.TypObject, GoTyp: "TypeA", Name: "Item"},
		}
		td2 := &jstypez.Typedef{
			Head: &jstypez.Tuple{Typ: jstypez.TypObject, GoTyp: "TypeB", Name: "Item"},
		}
		tds := jstypez.NewTypedefs([]*jstypez.Typedef{td1})
		require.Panics(t, func() {
			tds.Add(td2)
		})
	})
}
