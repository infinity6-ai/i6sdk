package forgepyschema_test

import (
	"testing"

	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/infinity6-ai/gox/schemaz/schemamold"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/gox/schemaz/typeastz"
	"github.com/infinity6-ai/i6sdk/forgez/forgepy/forgepyschema"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz/dummydataz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz/relevancefashiondataz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz/targetingfinancedataz"
	"github.com/stretchr/testify/require"
)

type externalCustomType struct {
	Value string
}

func TestUnitParserInstanceAndAddKnownType(t *testing.T) {
	ctx := t.Context()

	parser := schemamold.New().AddKnownType(externalCustomType{}, "MyCustomPyString")

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
	require.Equal(t, "MyCustomPyString", td.Properties.All()[0].Typ)
}

func TestUnitParseSchema(t *testing.T) {
	ctx := t.Context()
	api := routezsamplefraction.NewService().Api()
	params := api.ApiDataRefs().PathParams
	require.NotNil(t, params)

	td, tds := forgepyschema.NewSdk().ParseSchema(ctx, "sampleParams", params)
	require.NotNil(t, td)
	require.NotNil(t, tds)
	require.Equal(t, "sampleParams", td.Head.Name)
}

func TestUnitToPythonType(t *testing.T) {
	type testScenario struct {
		typ   string
		goTyp string
		want  string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		got := forgepyschema.ToPythonType(s.typ, s.goTyp)
		require.Equal(t, s.want, got)
	}

	t.Run("String primitive", func(t *testing.T) {
		check(t, testScenario{typ: "string", want: "str"})
	})

	t.Run("Integer number", func(t *testing.T) {
		check(t, testScenario{typ: "number", goTyp: "int64", want: "int"})
	})

	t.Run("Float number", func(t *testing.T) {
		check(t, testScenario{typ: "number", goTyp: "float64", want: "float"})
	})

	t.Run("Generic number", func(t *testing.T) {
		check(t, testScenario{typ: "number", want: "Union[int, float]"})
	})

	t.Run("Boolean", func(t *testing.T) {
		check(t, testScenario{typ: "boolean", want: "bool"})
	})

	t.Run("Array of strings", func(t *testing.T) {
		check(t, testScenario{typ: "Array<string>", want: "List[str]"})
	})

	t.Run("Record mapping", func(t *testing.T) {
		check(t, testScenario{typ: "Record<string, number>", want: "Dict[str, Union[int, float]]"})
	})

	t.Run("Embedded JSON", func(t *testing.T) {
		check(t, testScenario{typ: "i6JsonEmbedded<MyItem>", want: "Union[str, MyItem]"})
	})

	t.Run("Numeric union / enum", func(t *testing.T) {
		check(t, testScenario{typ: "1 | 2 | 4 | 8 | number", want: "int"})
	})

	t.Run("Time type", func(t *testing.T) {
		check(t, testScenario{typ: "dszTime", want: "str"})
	})

	t.Run("FileStat type", func(t *testing.T) {
		check(t, testScenario{typ: "dszFileStat", want: "FileStat"})
	})
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

	td, tds := forgepyschema.NewSdk().ParseSchema(ctx, "in", schema)
	require.Equal(t, "inParams", td.Properties.All()[0].Typ, "nested object must be typed")

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

	td, tds := forgepyschema.NewSdk().ParseSchema(ctx, "testReq", schema)
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
		in          typeastz.Typedefs
		mustContain []string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		got := forgepyschema.WriteTypedefs(ctx, s.in)
		for _, substr := range s.mustContain {
			require.Contains(t, got, substr)
		}
	}

	t.Run("empty typedefs generates base header and classes", func(t *testing.T) {
		check(t, testScenario{
			in: typeastz.NewTypedefs(nil),
			mustContain: []string{
				"from __future__ import annotations",
				"class ServiceOptions(TypedDict, total=False):",
				"class FileStat(TypedDict, total=False):",
			},
		})
	})

	t.Run("typedefs with class syntax", func(t *testing.T) {
		td := &typeastz.Typedef{
			Head:       &typeastz.Tuple{Typ: typeastz.TypObject, GoTyp: "models.Item", Name: "Item", Summary: "Item definition"},
			Properties: &typeastz.Tuples{},
		}
		td.Properties.Add(
			typeastz.Tuple{Typ: typeastz.TypString, GoTyp: "string", Name: "id", Summary: "Item ID"},
		)
		tds := typeastz.NewTypedefs([]*typeastz.Typedef{td})

		check(t, testScenario{
			in: tds,
			mustContain: []string{
				"class Item(TypedDict, total=False):",
				"Item definition",
				"[Go: models.Item]",
				"id: str",
			},
		})
	})

	t.Run("typedefs with python keyword uses functional syntax", func(t *testing.T) {
		td := &typeastz.Typedef{
			Head:       &typeastz.Tuple{Typ: typeastz.TypObject, Name: "QueryReq"},
			Properties: &typeastz.Tuples{},
		}
		td.Properties.Add(
			typeastz.Tuple{Typ: typeastz.TypString, Name: "from", Summary: "From offset"},
		)
		tds := typeastz.NewTypedefs([]*typeastz.Typedef{td})

		check(t, testScenario{
			in: tds,
			mustContain: []string{
				"QueryReq = TypedDict(",
				"\"from\": str,",
			},
		})
	})
}

func TestUnitParseTableDataSchemas(t *testing.T) {
	ctx := t.Context()
	p := forgepyschema.NewSdk()

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
