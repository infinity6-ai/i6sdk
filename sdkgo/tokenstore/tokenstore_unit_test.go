package tokenstore_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokenstore"
	"github.com/stretchr/testify/require"
)

func TestUnitMemoryTokenStore(t *testing.T) {
	ctx := t.Context()
	store := tokenstore.NewMemory()
	require.NotNil(t, store)
	defer func() {
		_ = store.Close()
	}()

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

	t.Run("Prefix listing and deletion", func(t *testing.T) {
		err := store.Set(ctx, "i6_ds:alpha", "val_a")
		require.NoError(t, err)
		err = store.Set(ctx, "i6_ds:beta", "val_b")
		require.NoError(t, err)
		err = store.Set(ctx, "i6_rf:gamma", "val_c")
		require.NoError(t, err)

		list, err := store.ListPrefix(ctx, "i6_ds:")
		require.NoError(t, err)
		require.Len(t, list, 2)
		require.Equal(t, "val_a", list["i6_ds:alpha"])
		require.Equal(t, "val_b", list["i6_ds:beta"])

		err = store.DeletePrefix(ctx, "i6_ds:")
		require.NoError(t, err)

		listAfter, err := store.ListPrefix(ctx, "i6_ds:")
		require.NoError(t, err)
		require.Empty(t, listAfter)

		// Non-matching prefix remains intact
		rem, err := store.Get(ctx, "i6_rf:gamma")
		require.NoError(t, err)
		require.Equal(t, "val_c", rem)
	})
}

func TestUnitSqliteTokenStore(t *testing.T) {
	ctx := t.Context()

	t.Run("Default URL resolves successfully", func(t *testing.T) {
		u, err := tokenstore.DefaultSqliteUrl()
		require.NoError(t, err)
		require.NotNil(t, u)
		require.Equal(t, "file", u.Scheme)
		require.Contains(t, u.String(), "tokens.db")
	})

	t.Run("Missing DB file is automatically created with 0600 permissions", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "tokenstore-sqlite-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		dbPath := filepath.Join(tmpDir, "sub", "tokens.db")
		u := urlz.MustParse("file://" + dbPath)
		store, err := tokenstore.NewSqlite(u)
		require.NoError(t, err)
		defer func() { _ = store.Close() }()

		// Verify file exists on disk
		info, err := os.Stat(dbPath)
		require.NoError(t, err)
		// Check file permissions (0600)
		require.Equal(t, os.FileMode(0600), info.Mode().Perm())

		// Missing key returns empty string without error
		val, err := store.Get(ctx, "missing_key")
		require.NoError(t, err)
		require.Empty(t, val)
	})

	t.Run("Set and Get persists and retrieves values", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "tokenstore-sqlite-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		dbPath := filepath.Join(tmpDir, "tokens.db")
		u := urlz.MustParse("file://" + dbPath)
		store, err := tokenstore.NewSqlite(u)
		require.NoError(t, err)
		defer func() { _ = store.Close() }()

		err = store.Set(ctx, "user_token", `{"access_token":"token-abc"}`)
		require.NoError(t, err)

		val, err := store.Get(ctx, "user_token")
		require.NoError(t, err)
		require.Equal(t, `{"access_token":"token-abc"}`, val)

		// Updating existing key (upsert)
		err = store.Set(ctx, "user_token", `{"access_token":"token-updated"}`)
		require.NoError(t, err)

		valUpdated, err := store.Get(ctx, "user_token")
		require.NoError(t, err)
		require.Equal(t, `{"access_token":"token-updated"}`, valUpdated)
	})

	t.Run("Delete removes key from store", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "tokenstore-sqlite-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		u := urlz.MustParse("file://" + filepath.Join(tmpDir, "tokens.db"))
		store, err := tokenstore.NewSqlite(u)
		require.NoError(t, err)
		defer func() { _ = store.Close() }()

		err = store.Set(ctx, "temp_key", "temp_value")
		require.NoError(t, err)

		err = store.Delete(ctx, "temp_key")
		require.NoError(t, err)

		val, err := store.Get(ctx, "temp_key")
		require.NoError(t, err)
		require.Empty(t, val)

		// Deleting non-existent key succeeds cleanly
		err = store.Delete(ctx, "nonexistent")
		require.NoError(t, err)
	})

	t.Run("Prefix listing and deletion", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "tokenstore-sqlite-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		u := urlz.MustParse("file://" + filepath.Join(tmpDir, "tokens.db"))
		store, err := tokenstore.NewSqlite(u)
		require.NoError(t, err)
		defer func() { _ = store.Close() }()

		err = store.Set(ctx, "i6_ds:orders", `{"dataset":"orders"}`)
		require.NoError(t, err)
		err = store.Set(ctx, "i6_ds:users", `{"dataset":"users"}`)
		require.NoError(t, err)
		err = store.Set(ctx, "i6_rf:orders", "robot-ed25519-key")
		require.NoError(t, err)
		err = store.Set(ctx, "i6_user_user_token", `{"user":"123"}`)
		require.NoError(t, err)

		// List dataset access tokens
		dsMap, err := store.ListPrefix(ctx, "i6_ds:")
		require.NoError(t, err)
		require.Len(t, dsMap, 2)
		require.Equal(t, `{"dataset":"orders"}`, dsMap["i6_ds:orders"])
		require.Equal(t, `{"dataset":"users"}`, dsMap["i6_ds:users"])

		// Delete dataset tokens by prefix
		err = store.DeletePrefix(ctx, "i6_ds:")
		require.NoError(t, err)

		dsMapAfter, err := store.ListPrefix(ctx, "i6_ds:")
		require.NoError(t, err)
		require.Empty(t, dsMapAfter)

		// Other keys are unaffected
		rfVal, err := store.Get(ctx, "i6_rf:orders")
		require.NoError(t, err)
		require.Equal(t, "robot-ed25519-key", rfVal)

		uVal, err := store.Get(ctx, "i6_user_user_token")
		require.NoError(t, err)
		require.Equal(t, `{"user":"123"}`, uVal)
	})

	t.Run("Concurrent Set and Get across goroutines", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "tokenstore-sqlite-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		u := urlz.MustParse("file://" + filepath.Join(tmpDir, "tokens.db"))
		store, err := tokenstore.NewSqlite(u)
		require.NoError(t, err)
		defer func() { _ = store.Close() }()

		var wg sync.WaitGroup
		concurrency := 20
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				k := fmt.Sprintf("i6_ds:dataset_%d", idx)
				v := fmt.Sprintf("val_%d", idx)
				err := store.Set(context.Background(), k, v)
				require.NoError(t, err)

				got, err := store.Get(context.Background(), k)
				require.NoError(t, err)
				require.Equal(t, v, got)
			}(i)
		}
		wg.Wait()
	})

	t.Run("In-memory SQLite mode works", func(t *testing.T) {
		u := urlz.MustParse("file:///:memory:")
		store, err := tokenstore.NewSqlite(u)
		require.NoError(t, err)
		defer func() { _ = store.Close() }()

		err = store.Set(ctx, "k1", "v1")
		require.NoError(t, err)

		got, err := store.Get(ctx, "k1")
		require.NoError(t, err)
		require.Equal(t, "v1", got)
	})
}
