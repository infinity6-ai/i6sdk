package jswritertypedef_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/genjsz/jswritertypedef"
)

func TestUnitWrite(t *testing.T) {
	ctx := t.Context()
	type testScenario struct {
		in   jstypez.Typedefs
		want string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		got := jswritertypedef.Write(ctx, s.in)
		require.Equal(t, s.want, got)
	}

	t.Run("empty typedefs", func(t *testing.T) {
		check(t, testScenario{
			in:   jstypez.NewTypedefs(nil),
			want: "/**\n * @typedef {import(\"../services/services.js\").Services} Services\n */\n\n",
		})
	})

	t.Run("typedefs with gotyp comment", func(t *testing.T) {
		td := &jstypez.Typedef{
			Head:       &jstypez.Tuple{Typ: jstypez.TypObject, GoTyp: "models.Item", Name: "Item", Summary: "Item definition"},
			Properties: &jstypez.Tuples{},
		}
		td.Properties.Add(
			jstypez.Tuple{Typ: jstypez.TypString, GoTyp: "string", Name: "id", Summary: "Item ID"},
		)
		tds := jstypez.NewTypedefs([]*jstypez.Typedef{td})

		check(t, testScenario{
			in:   tds,
			want: "/**\n * @typedef {import(\"../services/services.js\").Services} Services\n */\n\n/**\n * [Go: models.Item] Item definition\n * @typedef {object} Item\n * @property {string} id - [Go: string] Item ID\n */\n\n",
		})
	})
}
