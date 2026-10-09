package dummycatalogdataz_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/tabledataz/dummycatalogdataz"
)

func TestUnitDummyCatalogDataSerialization(t *testing.T) {
	type testScenario struct {
		jsonStr string
		catalog dummycatalogdataz.DummyCatalogData
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()

		// Test Deserialization (JSON -> Struct)
		var got dummycatalogdataz.DummyCatalogData
		err := json.Unmarshal([]byte(s.jsonStr), &got)
		require.NoError(t, err)
		require.Equal(t, s.catalog, got)

		// Test Serialization (Struct -> JSON -> Struct)
		serialized, err := json.Marshal(s.catalog)
		require.NoError(t, err)

		var roundTrip dummycatalogdataz.DummyCatalogData
		err = json.Unmarshal(serialized, &roundTrip)
		require.NoError(t, err)
		require.Equal(t, s.catalog, roundTrip)
	}

	t.Run("Basic dummy catalog fields", func(t *testing.T) {
		check(t, testScenario{
			jsonStr: `{
				"id": "dummy_123",
				"product": "Dummy Shirt",
				"price": 29.99,
				"qt_inventory": 150,
				"store": "Dummy Store A"
			}`,
			catalog: dummycatalogdataz.DummyCatalogData{
				Id:          "dummy_123",
				Product:     "Dummy Shirt",
				Price:       29.99,
				QtInventory: 150,
				Store:       "Dummy Store A",
			},
		})
	})

	t.Run("Zero price and empty inventory", func(t *testing.T) {
		check(t, testScenario{
			jsonStr: `{
				"id": "dummy_456",
				"product": "Free Sticker",
				"price": 0.0,
				"qt_inventory": 0,
				"store": "Dummy Store B"
			}`,
			catalog: dummycatalogdataz.DummyCatalogData{
				Id:          "dummy_456",
				Product:     "Free Sticker",
				Price:       0.0,
				QtInventory: 0,
				Store:       "Dummy Store B",
			},
		})
	})
}

func TestUnitDummyCatalogDataSchema(t *testing.T) {
	t.Run("DummyCatalogData Schema Check", func(t *testing.T) {
		s := (&dummycatalogdataz.DummyCatalogData{}).Schema()
		td, _ := jsparser.ParseSchema("DummyCatalogData", s)
		require.Equal(t, "DummyCatalogData", td.Head.Name)

		propMap := make(map[string]jstypez.Tuple)
		for _, prop := range td.Properties.All() {
			propMap[prop.Name] = prop
		}

		require.Contains(t, propMap, "id")
		require.Equal(t, jstypez.TypString, propMap["id"].Typ)

		require.Contains(t, propMap, "product")
		require.Equal(t, jstypez.TypString, propMap["product"].Typ)

		require.Contains(t, propMap, "price")
		require.Equal(t, jstypez.TypNumber, propMap["price"].Typ)

		require.Contains(t, propMap, "qt_inventory")
		require.Equal(t, jstypez.TypNumber, propMap["qt_inventory"].Typ)

		require.Contains(t, propMap, "store")
		require.Equal(t, jstypez.TypString, propMap["store"].Typ)
	})
}
