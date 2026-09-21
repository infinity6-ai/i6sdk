package generatorjsz_test

import (
	"testing"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/generatorjsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

// Str lê e escreve strings: mesmo quando o tipo nativo é um struct com MarshalJSON (o walker não
// consegue saber a forma JSON dele, KindUnknown), o lado JS enxerga uma string.
func TestUnitIntrospectStrWithMarshalerNativeIsString(t *testing.T) {
	var dataset lakedsz.Dataset
	s := &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"dataset": {Str: lakedsz.ParseStr(&dataset)},
			}
		},
	}

	n := generatorjsz.Introspect(s)
	require.Equal(t, generatorjsz.KindString, n.Fields["dataset"].Kind)
	require.True(t, n.Fields["dataset"].StringArrayWire)
}

// Raw com o mesmo tipo continua desconhecido: aqui a forma JSON real é a do MarshalJSON, que o
// walker não tem como inferir — melhor "*" do que chutar.
func TestUnitIntrospectRawWithMarshalerNativeIsUnknown(t *testing.T) {
	var dataset lakedsz.Dataset
	s := &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"dataset": {Raw: func() any { return &dataset }},
			}
		},
	}

	n := generatorjsz.Introspect(s)
	require.Equal(t, generatorjsz.KindUnknown, n.Fields["dataset"].Kind)
}
