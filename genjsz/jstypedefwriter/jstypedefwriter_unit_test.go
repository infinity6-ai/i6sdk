package jstypedefwriter_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypedefwriter"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func TestUnitWritePipedefs(t *testing.T) {
	type testScenario struct {
		name string
		in   jstypez.Typedefs
		want string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		got := jstypedefwriter.WritePipedefs(s.in)
		require.Equal(t, s.want, got)
	}

	t.Run("basic object", func(t *testing.T) {
		var tds jstypez.Typedefs
		td := &jstypez.Typedef{
			Head: &jstypez.Tuple{Typ: "object", Name: "MyType", Summary: "My summary"},
			Properties: &jstypez.Tuples{},
		}
		td.Properties.Add(
			jstypez.Tuple{Typ: "string", Name: "foo"},
			jstypez.Tuple{Typ: "number", Name: "bar", Summary: "bar summary"},
		)
		tds.Add(td)
		
		check(t, testScenario{
			name: "basic object",
			in:   tds,
			want: "/**\n * My summary\n * @typedef {Object} MyType\n * @property {string} foo\n * @property {number} bar - bar summary\n */\n\n",
		})
	})

	t.Run("multiple typedefs", func(t *testing.T) {
		var tds jstypez.Typedefs
		td1 := &jstypez.Typedef{
			Head: &jstypez.Tuple{Typ: "string", Name: "AliasType"},
		}
		td2 := &jstypez.Typedef{
			Head: &jstypez.Tuple{Typ: "object", Name: "BType"},
			Properties: &jstypez.Tuples{},
		}
		td2.Properties.Add(jstypez.Tuple{Typ: "number", Name: "x"})
		
		tds.Add(td1, td2)
		
		check(t, testScenario{
			name: "multiple typedefs",
			in:   tds,
			want: "/**\n * @typedef {string} AliasType\n */\n\n/**\n * @typedef {Object} BType\n * @property {number} x\n */\n\n",
		})
	})
}
