package jsparser_test

import (
	"testing"

	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/sdkgo/grants"
)

func TestUnitParseSchema(t *testing.T) {
	ctx := t.Context()
	api := routezsamplefraction.NewService().Api()
	params := api.ApiDataRefs().PathParams
	require.NotNil(t, params)

	td, tds := jsparser.ParseSchema(ctx, "sampleParams", params)
	require.NotNil(t, td)
	require.NotNil(t, tds)
	require.Equal(t, "sampleParams", td.Head.Name)
}

func TestUnitDescHelpers(t *testing.T) {
	d := &schemaz.Desc{
		Name:    "TestName",
		Summary: "TestSummary",
	}
	require.Equal(t, "TestName", jsparser.DescName(d))
	require.Equal(t, "TestSummary", jsparser.DescSummary(d))
	require.Equal(t, "", jsparser.DescName(nil))
	require.Equal(t, "", jsparser.DescSummary(nil))
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

	td, tds := jsparser.ParseSchema(ctx, "in", schema)
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

	td, tds := jsparser.ParseSchema(ctx, "testReq", schema)
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

// func TestUnitParseSchemaAnyField(t *testing.T) {
// 	ctx := t.Context()
// 	var anyVal any // nil interface
// 	schema := &schemaz.Schema{
// 		Object: func(_ bool) map[string]*schemaz.Schema {
// 			return map[string]*schemaz.Schema{
// 				"inputs": {
// 					Field: func() any { return &anyVal },
// 				},
// 			}
// 		},
// 	}

// 	td, tds := jsparser.ParseSchema(ctx, "in", schema)
// 	require.NotNil(t, td)
// 	require.Empty(t, tds.All())
// 	require.Equal(t, "Record<string, any>", td.Properties.All()[0].Typ)
// }
