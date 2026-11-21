package targetingfinancedataz_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/forgez/forgejs/forgejsschema"
	"github.com/infinity6-ai/gox/schemaz/typeastz"
	"go.code.infinity6.ai/sdkgo/sdkgo/tabledataz"
	"go.code.infinity6.ai/sdkgo/sdkgo/tabledataz/targetingfinancedataz"
)

var _ tabledataz.TableData = (*targetingfinancedataz.ClientFeatures)(nil)

func validFeatures() targetingfinancedataz.ClientFeatures {
	return targetingfinancedataz.ClientFeatures{
		ID:        "bbb",
		StoreId:   "s1",
		CreatedAt: time.Date(2025, 6, 2, 0, 0, 0, 0, time.UTC),
		FeaturesFloat: []*targetingfinancedataz.FeatureFloat{
			{Key: "tempo_exercicio_meses", Value: 8},
			{Key: "quantidade_de_atendimento_mes", Value: 10},
		},
	}
}

func TestUnitClientFeaturesSerialization(t *testing.T) {
	jsonStr := `{
		"id": "bbb",
		"store_id": "s1",
		"created_at": "2025-06-02T00:00:00Z",
		"features_float": [
			{"key": "tempo_exercicio_meses", "value": 8},
			{"key": "quantidade_de_atendimento_mes", "value": 10}
		]
	}`

	var got targetingfinancedataz.ClientFeatures
	require.NoError(t, json.Unmarshal([]byte(jsonStr), &got))
	require.Equal(t, validFeatures(), got)

	b, err := json.Marshal(got)
	require.NoError(t, err)
	var roundTrip targetingfinancedataz.ClientFeatures
	require.NoError(t, json.Unmarshal(b, &roundTrip))
	require.Equal(t, got, roundTrip)
}

func TestUnitClientFeaturesValidation(t *testing.T) {
	type scenario struct {
		mutate func(c *targetingfinancedataz.ClientFeatures)
		errMsg string
	}
	check := func(t *testing.T, sc scenario) {
		t.Helper()
		c := validFeatures()
		if sc.mutate != nil {
			sc.mutate(&c)
		}
		err := c.Validate()
		if sc.errMsg == "" {
			require.NoError(t, err)
			return
		}
		require.ErrorContains(t, err, sc.errMsg)
	}

	t.Run("valid", func(t *testing.T) { check(t, scenario{}) })
	t.Run("any key name", func(t *testing.T) {
		check(t, scenario{mutate: func(c *targetingfinancedataz.ClientFeatures) {
			c.FeaturesFloat = append(c.FeaturesFloat, &targetingfinancedataz.FeatureFloat{Key: "Qualquer Nome-123", Value: -1.5})
		}})
	})
	t.Run("missing id", func(t *testing.T) {
		check(t, scenario{mutate: func(c *targetingfinancedataz.ClientFeatures) { c.ID = "" }, errMsg: "failed to validate id"})
	})
	t.Run("missing created_at", func(t *testing.T) {
		check(t, scenario{mutate: func(c *targetingfinancedataz.ClientFeatures) { c.CreatedAt = time.Time{} }, errMsg: "failed to validate created_at"})
	})
	t.Run("empty features", func(t *testing.T) {
		check(t, scenario{mutate: func(c *targetingfinancedataz.ClientFeatures) { c.FeaturesFloat = nil }, errMsg: "failed to validate features_float"})
	})
	t.Run("nil feature", func(t *testing.T) {
		check(t, scenario{mutate: func(c *targetingfinancedataz.ClientFeatures) { c.FeaturesFloat[1] = nil }, errMsg: "feature cannot be nil"})
	})
	t.Run("empty key", func(t *testing.T) {
		check(t, scenario{mutate: func(c *targetingfinancedataz.ClientFeatures) { c.FeaturesFloat[0].Key = "" }, errMsg: "failed to validate key"})
	})
	t.Run("duplicated key", func(t *testing.T) {
		check(t, scenario{mutate: func(c *targetingfinancedataz.ClientFeatures) { c.FeaturesFloat[1].Key = c.FeaturesFloat[0].Key }, errMsg: "duplicated key"})
	})
}

func TestUnitClientFeaturesSchema(t *testing.T) {
	ctx := t.Context()

	t.Run("FeatureFloat Schema", func(t *testing.T) {
		td, _ := forgejsschema.NewSdk().ParseSchema(ctx, "FeatureFloat", (&targetingfinancedataz.FeatureFloat{}).Schema())
		require.Equal(t, "FeatureFloat", td.Head.Name)
	})

	t.Run("ClientFeatures body", func(t *testing.T) {
		td, tds := forgejsschema.NewSdk().ParseSchema(ctx, "ClientFeatures", (&targetingfinancedataz.ClientFeatures{}).SchemaBody())
		require.Equal(t, "ClientFeatures", td.Head.Name)

		props := map[string]typeastz.Tuple{}
		for _, p := range td.Properties.All() {
			props[p.Name] = p
		}
		require.Equal(t, typeastz.TypString, props["id"].Typ)
		require.Equal(t, typeastz.TypTime, props["created_at"].Typ)
		require.Contains(t, props, "features_float")

		names := map[string]struct{}{}
		for _, d := range tds.All() {
			names[d.Head.Name] = struct{}{}
		}
		require.Contains(t, names, "ClientFeaturesFeaturesFloatElement")
	})

	t.Run("ClientFeatures query", func(t *testing.T) {
		td, _ := forgejsschema.NewSdk().ParseSchema(ctx, "ClientFeatures", (&targetingfinancedataz.ClientFeatures{}).SchemaQuery())
		require.Equal(t, "ClientFeatures", td.Head.Name)
	})
}
