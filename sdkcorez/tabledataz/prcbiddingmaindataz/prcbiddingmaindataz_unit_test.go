package prcbiddingmaindataz_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz/prcbiddingmaindataz"
)

func TestUnitCatalog(t *testing.T) {
	var _ tabledataz.TableData = (*prcbiddingmaindataz.Catalog)(nil)

	t.Run("Valid catalog", func(t *testing.T) {
		cat := prcbiddingmaindataz.Catalog{
			ID:            "cat_1",
			Desc:          "description",
			Supplier:      "supp_1",
			Status:        "active",
			ExclusiveSale: "yes",
		}
		require.NoError(t, cat.Validate())
	})

	type scenario struct {
		mutate func(c *prcbiddingmaindataz.Catalog)
		errMsg string
	}

	check := func(t *testing.T, s scenario) {
		t.Helper()
		cat := prcbiddingmaindataz.Catalog{
			ID:            "cat_1",
			Desc:          "description",
			Supplier:      "supp_1",
			Status:        "active",
			ExclusiveSale: "yes",
		}
		s.mutate(&cat)
		err := cat.Validate()
		if s.errMsg != "" {
			require.Error(t, err)
			require.Contains(t, err.Error(), s.errMsg)
			return
		}
		require.NoError(t, err)
	}

	t.Run("Invalid: Missing ID", func(t *testing.T) {
		check(t, scenario{mutate: func(c *prcbiddingmaindataz.Catalog) { c.ID = "" }, errMsg: "failed to validate id"})
	})
	t.Run("Invalid: Missing Desc", func(t *testing.T) {
		check(t, scenario{mutate: func(c *prcbiddingmaindataz.Catalog) { c.Desc = "" }, errMsg: "failed to validate desc"})
	})
	t.Run("Invalid: Missing Supplier", func(t *testing.T) {
		check(t, scenario{mutate: func(c *prcbiddingmaindataz.Catalog) { c.Supplier = "" }, errMsg: "failed to validate supplier"})
	})
	t.Run("Invalid: Missing Status", func(t *testing.T) {
		check(t, scenario{mutate: func(c *prcbiddingmaindataz.Catalog) { c.Status = "" }, errMsg: "failed to validate status"})
	})
	t.Run("Invalid: Missing ExclusiveSale", func(t *testing.T) {
		check(t, scenario{mutate: func(c *prcbiddingmaindataz.Catalog) { c.ExclusiveSale = "" }, errMsg: "failed to validate exclusive_sale"})
	})
	t.Run("Invalid: ID too long", func(t *testing.T) {
		longID := string(make([]byte, 65))
		check(t, scenario{mutate: func(c *prcbiddingmaindataz.Catalog) { c.ID = longID }, errMsg: "failed to validate id"})
	})
}

func TestUnitCatalogSerialization(t *testing.T) {
	cat := prcbiddingmaindataz.Catalog{
		ID:            "cat_123",
		Desc:          "A great product",
		Supplier:      "Supplier Co",
		Status:        "available",
		ExclusiveSale: "true",
	}

	serialized, err := json.Marshal(cat)
	require.NoError(t, err)

	var roundTrip prcbiddingmaindataz.Catalog
	err = json.Unmarshal(serialized, &roundTrip)
	require.NoError(t, err)
	require.Equal(t, cat, roundTrip)
}

func TestUnitCatalogSchema(t *testing.T) {
	t.Run("Catalog Schema Query", func(t *testing.T) {
		s := (&prcbiddingmaindataz.Catalog{}).SchemaQuery()
		require.NotNil(t, s)
		require.NotNil(t, s.Object)

		props := s.Object(true)
		require.Contains(t, props, "id")
		require.Contains(t, props, "desc")
		require.Contains(t, props, "supplier")
		require.Contains(t, props, "status")
		require.Contains(t, props, "exclusive_sale")
	})

	t.Run("Catalog Schema Body", func(t *testing.T) {
		s := (&prcbiddingmaindataz.Catalog{}).SchemaBody()
		require.NotNil(t, s)
		require.NotNil(t, s.Object)

		props := s.Object(true)
		require.Contains(t, props, "id")
		require.Contains(t, props, "desc")
		require.Contains(t, props, "supplier")
		require.Contains(t, props, "status")
		require.Contains(t, props, "exclusive_sale")
	})
}
