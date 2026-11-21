package forgejsschema_test

import (
	"testing"

	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/infinity6-ai/gox/schemaz/schemamold"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/gox/schemaz/typeastz"
	"github.com/stretchr/testify/require"
	"github.com/infinity6-ai/i6sdk/forgez/forgejs/forgejsschema"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz/dummydataz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz/relevancefashiondataz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz/targetingfinancedataz"
)

type externalCustomType struct {
	Value string
}

func TestUnitParserInstanceAndAddKnownType(t *testing.T) {
	ctx := t.Context()

	parser := schemamold.New().AddKnownType(externalCustomType{}, "MyCustomJsString")

	var val externalCustomType
	schema := &schemaz.Schema{
		Object: func(_ bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"custom_field": {
					Field: func() any { return &val },
				},
			}
		},
	}

	td, _ := parser.ParseSchema(ctx, "testCustom", schema)
	require.NotNil(t, td)
	require.Len(t, td.Properties.All(), 1)
	require.Equal(t, "custom_field", td.Properties.All()[0].Name)
	require.Equal(t, "MyCustomJsString", td.Properties.All()[0].Typ)
}

func TestUnitParseSchema(t *testing.T) {
	ctx := t.Context()
	api := routezsamplefraction.NewService().Api()
	params := api.ApiDataRefs().PathParams
	require.NotNil(t, params)

	td, tds := forgejsschema.NewSdk().ParseSchema(ctx, "sampleParams", params)
	require.NotNil(t, td)
	require.NotNil(t, tds)
	require.Equal(t, "sampleParams", td.Head.Name)
}

func TestUnitDescHelpers(t *testing.T) {
	d := &schemaz.Desc{
		Name:    "TestName",
		Summary: "TestSummary",
	}
	require.Equal(t, "TestName", schemamold.DescName(d))
	require.Equal(t, "TestSummary", schemamold.DescSummary(d))
	require.Equal(t, "", schemamold.DescName(nil))
	require.Equal(t, "", schemamold.DescSummary(nil))
}

func TestUnitParseSchemaNestedObject(t *testing.T) {
	ctx := t.Context()
	var s string
	schema := &schemaz.Schema{
		Object: func(_ bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"params": {
					Object: func(_ bool) map[string]*schemaz.Schema {
						return map[string]*schemaz.Schema{
							"min": {Field: func() any { return &s }},
						}
					},
				},
			}
		},
	}

	td, tds := forgejsschema.NewSdk().ParseSchema(ctx, "in", schema)
	require.Equal(t, "inParams", td.Properties.All()[0].Typ, "nested object must be typed, not a bare object")

	var found bool
	for _, d := range tds.All() {
		if d.Head.Name == "inParams" {
			found = true
			require.Equal(t, "string", d.Properties.All()[0].Typ)
		}
	}
	require.True(t, found)
}

func TestUnitParseSchemaCustomTypeEnum(t *testing.T) {
	ctx := t.Context()
	var g grants.Grant
	schema := &schemaz.Schema{
		Object: func(_ bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"grants": {
					Field: func() any { return &g },
				},
			}
		},
	}

	td, tds := forgejsschema.NewSdk().ParseSchema(ctx, "testReq", schema)
	require.NotNil(t, td)
	require.Equal(t, "Grant", td.Properties.All()[0].Typ)

	var found bool
	for _, d := range tds.All() {
		if d.Head.Name == "Grant" {
			found = true
			require.Equal(t, "1 | 2 | 4 | 8 | number", d.Head.Typ)
			require.Equal(t, "grants.Grant", d.Head.GoTyp)
		}
	}
	require.True(t, found)
}

func TestUnitWriteTypedefs(t *testing.T) {
	ctx := t.Context()
	type testScenario struct {
		in   typeastz.Typedefs
		want string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		got := forgejsschema.WriteTypedefs(ctx, s.in)
		require.Equal(t, s.want, got)
	}

	t.Run("empty typedefs", func(t *testing.T) {
		check(t, testScenario{
			in:   typeastz.NewTypedefs(nil),
			want: "/**\n * @typedef {import(\"../services/services.js\").Services} Services\n */\n\n/**\n * @typedef {import(\"../invoker/invoker.js\").ServiceOptions} ServiceOptions\n */\n\n/**\n * @template T\n * @typedef {T} i6JsonEmbedded\n */\n\n",
		})
	})

	t.Run("typedefs with gotyp comment", func(t *testing.T) {
		td := &typeastz.Typedef{
			Head:       &typeastz.Tuple{Typ: typeastz.TypObject, GoTyp: "models.Item", Name: "Item", Summary: "Item definition"},
			Properties: &typeastz.Tuples{},
		}
		td.Properties.Add(
			typeastz.Tuple{Typ: typeastz.TypString, GoTyp: "string", Name: "id", Summary: "Item ID"},
		)
		tds := typeastz.NewTypedefs([]*typeastz.Typedef{td})

		check(t, testScenario{
			in:   tds,
			want: "/**\n * @typedef {import(\"../services/services.js\").Services} Services\n */\n\n/**\n * @typedef {import(\"../invoker/invoker.js\").ServiceOptions} ServiceOptions\n */\n\n/**\n * @template T\n * @typedef {T} i6JsonEmbedded\n */\n\n/**\n * [Go: models.Item] Item definition\n * @typedef {object} Item\n * @property {string} id - [Go: string] Item ID\n */\n\n",
		})
	})
}

func TestUnitParseTableDataSchemas(t *testing.T) {
	ctx := t.Context()
	p := forgejsschema.NewSdk()

	t.Run("dummy catalog data", func(t *testing.T) {
		td, _ := p.ParseSchema(ctx, "DummyCatalogData", (&dummydataz.Catalog{}).SchemaQuery())
		require.NotNil(t, td)
		require.Equal(t, "DummyCatalogData", td.Head.Name)
	})

	t.Run("relevance fashion event", func(t *testing.T) {
		td, tds := p.ParseSchema(ctx, "Event", (&relevancefashiondataz.Event{}).SchemaQuery())
		require.NotNil(t, td)
		require.Equal(t, "Event", td.Head.Name)
		require.NotNil(t, tds)
	})

	t.Run("targeting finance client status", func(t *testing.T) {
		td, _ := p.ParseSchema(ctx, "ClientStatus", (&targetingfinancedataz.ClientStatus{}).SchemaBody())
		require.NotNil(t, td)
		require.Equal(t, "ClientStatus", td.Head.Name)
	})
}

