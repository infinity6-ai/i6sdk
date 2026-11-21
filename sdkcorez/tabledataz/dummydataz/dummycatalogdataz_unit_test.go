package dummydataz_test

import (
	"encoding/json"
	"testing"

	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz/dummydataz"
	"github.com/stretchr/testify/require"
)

var _ tabledataz.TableData = (*dummydataz.Catalog)(nil)

func TestUnitDummyCatalogDataSerialization(t *testing.T) {
	type testScenario struct {
		jsonStr string
		catalog dummydataz.Catalog
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()

		// Test Deserialization (JSON -> Struct)
		var got dummydataz.Catalog
		err := json.Unmarshal([]byte(s.jsonStr), &got)
		require.NoError(t, err)
		require.Equal(t, s.catalog, got)

		// Test Serialization (Struct -> JSON -> Struct)
		serialized, err := json.Marshal(s.catalog)
		require.NoError(t, err)

		var roundTrip dummydataz.Catalog
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
				"store_id": "Dummy Store A"
			}`,
			catalog: dummydataz.Catalog{
				Id:          "dummy_123",
				Product:     "Dummy Shirt",
				Price:       29.99,
				QtInventory: 150,
				StoreId:     "Dummy Store A",
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
				"store_id": "Dummy Store B"
			}`,
			catalog: dummydataz.Catalog{
				Id:          "dummy_456",
				Product:     "Free Sticker",
				Price:       0.0,
				QtInventory: 0,
				StoreId:     "Dummy Store B",
			},
		})
	})
}

func TestUnitDummyCatalogDataSchema(t *testing.T) {
	t.Run("DummyCatalogData Schema Check", func(t *testing.T) {
		s := (&dummydataz.Catalog{}).SchemaQuery()
		require.NotNil(t, s)
		require.NotNil(t, s.Object)

		props := s.Object(true)
		require.Contains(t, props, "id")
		require.Contains(t, props, "product")
		require.Contains(t, props, "price")
		require.Contains(t, props, "qt_inventory")
		require.Contains(t, props, "store_id")
	})
}
