package jsschema_test

import (
	"testing"

	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsschema"
	"go.code.infinity6.ai/sdkgo/sdkgo/grants"
)

func TestUnitParseSchema(t *testing.T) {
	ctx := t.Context()
	api := routezsamplefraction.NewService().Api()
	params := api.ApiDataRefs().PathParams
	require.NotNil(t, params)

	td, tds := jsschema.ParseSchema(ctx, "sampleParams", params)
	require.NotNil(t, td)
	require.NotNil(t, tds)
	require.Equal(t, "sampleParams", td.Head.Name)
}

func TestUnitDescHelpers(t *testing.T) {
	d := &schemaz.Desc{
		Name:    "TestName",
		Summary: "TestSummary",
	}
	require.Equal(t, "TestName", jsschema.DescName(d))
	require.Equal(t, "TestSummary", jsschema.DescSummary(d))
	require.Equal(t, "", jsschema.DescName(nil))
	require.Equal(t, "", jsschema.DescSummary(nil))
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

	td, tds := jsschema.ParseSchema(ctx, "in", schema)
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

	td, tds := jsschema.ParseSchema(ctx, "testReq", schema)
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
