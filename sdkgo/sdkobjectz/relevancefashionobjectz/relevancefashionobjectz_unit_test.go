package relevancefashionobjectz_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkobjectz/relevancefashionobjectz"
)

func TestUnitEventSerialization(t *testing.T) {
	type testScenario struct {
		jsonStr string
		event   relevancefashionobjectz.Event
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()

		// Test Deserialization (JSON -> Struct)
		var got relevancefashionobjectz.Event
		err := json.Unmarshal([]byte(s.jsonStr), &got)
		require.NoError(t, err)

		// Assert that the timestamp matches within a second (to avoid minor format/timezone differences in string serialization)
		require.WithinDuration(t, s.event.CreatedAt, got.CreatedAt, time.Second)

		// Align CreatedAt exactly to perform exact field value checks
		got.CreatedAt = s.event.CreatedAt
		require.Equal(t, s.event, got)

		// Test Serialization (Struct -> JSON -> Struct)
		serialized, err := json.Marshal(s.event)
		require.NoError(t, err)

		var roundTrip relevancefashionobjectz.Event
		err = json.Unmarshal(serialized, &roundTrip)
		require.NoError(t, err)

		roundTrip.CreatedAt = s.event.CreatedAt
		require.Equal(t, s.event, roundTrip)
	}

	t.Run("Basic event fields", func(t *testing.T) {
		check(t, testScenario{
			jsonStr: `{
				"id": "evt_123",
				"store_id": "ABC",
				"session_id": "sess_456",
				"user_id": "usr_789",
				"created_at": "2026-10-08T12:00:00Z"
			}`,
			event: relevancefashionobjectz.Event{
				Id:        "evt_123",
				StoreId:   "ABC",
				SessionId: "sess_456",
				UserId:    "usr_789",
				CreatedAt: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC),
			},
		})
	})

	t.Run("Click event payload", func(t *testing.T) {
		check(t, testScenario{
			jsonStr: `{
				"id": "evt_123",
				"store_id": "ABC",
				"session_id": "sess_456",
				"user_id": "usr_789",
				"created_at": "2026-10-08T12:00:00Z",
				"click": {
					"price": 166.4,
					"sku_id": "146"
				}
			}`,
			event: relevancefashionobjectz.Event{
				Id:        "evt_123",
				StoreId:   "ABC",
				SessionId: "sess_456",
				UserId:    "usr_789",
				CreatedAt: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC),
				Click: &relevancefashionobjectz.Click{
					Price: 166.4,
					SkuId: "146",
				},
			},
		})
	})

	t.Run("AddToCart event payload", func(t *testing.T) {
		check(t, testScenario{
			jsonStr: `{
				"id": "evt_123",
				"store_id": "ABC",
				"session_id": "sess_456",
				"user_id": "usr_789",
				"created_at": "2026-10-08T12:00:00Z",
				"add_to_cart": {
					"price": 166.42,
					"sku_id": "121"
				}
			}`,
			event: relevancefashionobjectz.Event{
				Id:        "evt_123",
				StoreId:   "ABC",
				SessionId: "sess_456",
				UserId:    "usr_789",
				CreatedAt: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC),
				AddToCart: &relevancefashionobjectz.AddToCart{
					Price: 166.42,
					SkuId: "121",
				},
			},
		})
	})

	t.Run("RemoveFromCart event payload", func(t *testing.T) {
		check(t, testScenario{
			jsonStr: `{
				"id": "evt_123",
				"store_id": "ABC",
				"session_id": "sess_456",
				"user_id": "usr_789",
				"created_at": "2026-10-08T12:00:00Z",
				"remove_from_cart": {
					"price": 61.81,
					"sku_id": "113"
				}
			}`,
			event: relevancefashionobjectz.Event{
				Id:        "evt_123",
				StoreId:   "ABC",
				SessionId: "sess_456",
				UserId:    "usr_789",
				CreatedAt: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC),
				RemoveFromCart: &relevancefashionobjectz.RemoveFromCart{
					Price: 61.81,
					SkuId: "113",
				},
			},
		})
	})

	t.Run("Search event payload", func(t *testing.T) {
		check(t, testScenario{
			jsonStr: `{
				"id": "evt_123",
				"store_id": "ABC",
				"session_id": "sess_456",
				"user_id": "usr_789",
				"created_at": "2026-10-08T12:00:00Z",
				"search": {
					"text": "camisa da barbie"
				}
			}`,
			event: relevancefashionobjectz.Event{
				Id:        "evt_123",
				StoreId:   "ABC",
				SessionId: "sess_456",
				UserId:    "usr_789",
				CreatedAt: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC),
				Search: &relevancefashionobjectz.Search{
					Text: "camisa da barbie",
				},
			},
		})
	})

	t.Run("Purchase event payload", func(t *testing.T) {
		check(t, testScenario{
			jsonStr: `{
				"id": "evt_123",
				"store_id": "ABC",
				"session_id": "sess_456",
				"user_id": "usr_789",
				"created_at": "2026-10-08T12:00:00Z",
				"purchase": {
					"order_id": "8e37896a-7669-48b7-95fb-a661d6a1e597",
					"skus": [
						{
							"qtty": 1,
							"sku_id": "164",
							"unit_price": 259.29
						}
					]
				}
			}`,
			event: relevancefashionobjectz.Event{
				Id:        "evt_123",
				StoreId:   "ABC",
				SessionId: "sess_456",
				UserId:    "usr_789",
				CreatedAt: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC),
				Purchase: &relevancefashionobjectz.Purchase{
					OrderId: "8e37896a-7669-48b7-95fb-a661d6a1e597",
					Skus: []*relevancefashionobjectz.PurchaseSku{
						{
							Qtty:      1,
							SkuId:     "164",
							UnitPrice: 259.29,
						},
					},
				},
			},
		})
	})

	t.Run("Recraw event payload", func(t *testing.T) {
		check(t, testScenario{
			jsonStr: `{
				"id": "evt_123",
				"store_id": "ABC",
				"session_id": "sess_456",
				"user_id": "usr_789",
				"created_at": "2026-10-08T12:00:00Z",
				"recraw": {
					"rec_id": "xpto",
					"model_id": "relevance-fashion-fbt",
					"payload": "{\"results\":[{\"skus\":[{\"sku_id\":\"28\"}],\"rating\":0.0726400166749954}]}"
				}
			}`,
			event: relevancefashionobjectz.Event{
				Id:        "evt_123",
				StoreId:   "ABC",
				SessionId: "sess_456",
				UserId:    "usr_789",
				CreatedAt: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC),
				Recraw: &relevancefashionobjectz.Recraw{
					RecId:   "xpto",
					ModelId: "relevance-fashion-fbt",
					Payload: `{"results":[{"skus":[{"sku_id":"28"}],"rating":0.0726400166749954}]}`,
				},
			},
		})
	})

	t.Run("Recshow event payload", func(t *testing.T) {
		check(t, testScenario{
			jsonStr: `{
				"id": "evt_123",
				"store_id": "ABC",
				"session_id": "sess_456",
				"user_id": "usr_789",
				"created_at": "2026-10-08T12:00:00Z",
				"recshow": {
					"rec_id": "xpto",
					"payload": "{\"skus\":[{\"sku_id\":\"107\"}]}"
				}
			}`,
			event: relevancefashionobjectz.Event{
				Id:        "evt_123",
				StoreId:   "ABC",
				SessionId: "sess_456",
				UserId:    "usr_789",
				CreatedAt: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC),
				Recshow: &relevancefashionobjectz.Recshow{
					RecId:   "xpto",
					Payload: `{"skus":[{"sku_id":"107"}]}`,
				},
			},
		})
	})
}

func TestUnitEventSchema(t *testing.T) {
	t.Run("Click Schema", func(t *testing.T) {
		s := (&relevancefashionobjectz.Click{}).Schema()
		td, _ := jsparser.ParseSchema("Click", s)
		require.Equal(t, "Click", td.Head.Name)
	})
	t.Run("AddToCart Schema", func(t *testing.T) {
		s := (&relevancefashionobjectz.AddToCart{}).Schema()
		td, _ := jsparser.ParseSchema("AddToCart", s)
		require.Equal(t, "AddToCart", td.Head.Name)
	})
	t.Run("RemoveFromCart Schema", func(t *testing.T) {
		s := (&relevancefashionobjectz.RemoveFromCart{}).Schema()
		td, _ := jsparser.ParseSchema("RemoveFromCart", s)
		require.Equal(t, "RemoveFromCart", td.Head.Name)
	})
	t.Run("Search Schema", func(t *testing.T) {
		s := (&relevancefashionobjectz.Search{}).Schema()
		td, _ := jsparser.ParseSchema("Search", s)
		require.Equal(t, "Search", td.Head.Name)
	})
	t.Run("PurchaseSku Schema", func(t *testing.T) {
		s := (&relevancefashionobjectz.PurchaseSku{}).Schema()
		td, _ := jsparser.ParseSchema("PurchaseSku", s)
		require.Equal(t, "PurchaseSku", td.Head.Name)
	})
	t.Run("Purchase Schema", func(t *testing.T) {
		s := (&relevancefashionobjectz.Purchase{}).Schema()
		td, _ := jsparser.ParseSchema("Purchase", s)
		require.Equal(t, "Purchase", td.Head.Name)
	})
	t.Run("Recraw Schema", func(t *testing.T) {
		s := (&relevancefashionobjectz.Recraw{}).Schema()
		td, _ := jsparser.ParseSchema("Recraw", s)
		require.Equal(t, "Recraw", td.Head.Name)
	})
	t.Run("Recshow Schema", func(t *testing.T) {
		s := (&relevancefashionobjectz.Recshow{}).Schema()
		td, _ := jsparser.ParseSchema("Recshow", s)
		require.Equal(t, "Recshow", td.Head.Name)
	})
	t.Run("Event Schema", func(t *testing.T) {
		s := (&relevancefashionobjectz.Event{}).Schema()
		td, tds := jsparser.ParseSchema("Event", s)
		require.Equal(t, "Event", td.Head.Name)

		propMap := make(map[string]jstypez.Tuple)
		for _, prop := range td.Properties.All() {
			propMap[prop.Name] = prop
		}

		require.Contains(t, propMap, "id")
		require.Equal(t, "string", propMap["id"].Typ)

		require.Contains(t, propMap, "created_at")
		require.Equal(t, jstypez.TypTime, propMap["created_at"].Typ)

		require.Contains(t, propMap, "click")
		require.Equal(t, "object", propMap["click"].Typ)

		require.Contains(t, propMap, "purchase")
		require.Equal(t, "object", propMap["purchase"].Typ)

		// Check registered sub-typedefs
		tdfMap := make(map[string]*jstypez.Typedef)
		for _, tdf := range tds.All() {
			tdfMap[tdf.Head.Name] = tdf
		}

		require.Contains(t, tdfMap, "EventPurchaseSkusElement")
	})
}
