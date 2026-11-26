package tokenstore_test

import (
	"testing"

	"github.com/infinity6-ai/i6sdk/sdkgo/tokenstore"
	"github.com/stretchr/testify/require"
)

func TestUnitMemoryTokenStore(t *testing.T) {
	ctx := t.Context()
	store := tokenstore.NewMemory()
	require.NotNil(t, store)

	t.Run("Get missing key returns empty string", func(t *testing.T) {
		val, err := store.Get(ctx, "nonexistent")
		require.NoError(t, err)
		require.Empty(t, val)
	})

	t.Run("Set and Get stores and retrieves value", func(t *testing.T) {
		err := store.Set(ctx, "access_token", "jwt-token-123")
		require.NoError(t, err)

		val, err := store.Get(ctx, "access_token")
		require.NoError(t, err)
		require.Equal(t, "jwt-token-123", val)
	})

	t.Run("Delete removes key from store", func(t *testing.T) {
		err := store.Set(ctx, "refresh_token", "ref-xyz")
		require.NoError(t, err)

		err = store.Delete(ctx, "refresh_token")
		require.NoError(t, err)

		val, err := store.Get(ctx, "refresh_token")
		require.NoError(t, err)
		require.Empty(t, val)
	})
}
