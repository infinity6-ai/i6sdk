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

func TestUnitFszTokenStore(t *testing.T) {
	ctx := t.Context()

	t.Run("Default URL resolves successfully", func(t *testing.T) {
		u, err := tokenstore.DefaultFszUrl()
		require.NoError(t, err)
		require.NotNil(t, u)
		require.Equal(t, "file", u.Scheme)
		require.Contains(t, u.String(), "tokens.json")
	})

	t.Run("Missing file returns empty string without error", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "tokenstore-fsz-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		u := urlz.MustParse("file://" + filepath.Join(tmpDir, "tokens.json"))
		store, err := tokenstore.NewFsz(u)
		require.NoError(t, err)

		val, err := store.Get(ctx, "missing_key")
		require.NoError(t, err)
		require.Empty(t, val)
	})

	t.Run("Set and Get persists and retrieves values", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "tokenstore-fsz-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		filePath := filepath.Join(tmpDir, "sub", "tokens.json")
		u := urlz.MustParse("file://" + filePath)
		store, err := tokenstore.NewFsz(u)
		require.NoError(t, err)

		err = store.Set(ctx, "user_token", `{"access_token":"token-abc"}`)
		require.NoError(t, err)

		// Verify file exists on disk
		info, err := os.Stat(filePath)
		require.NoError(t, err)
		require.Greater(t, info.Size(), int64(0))

		val, err := store.Get(ctx, "user_token")
		require.NoError(t, err)
		require.Equal(t, `{"access_token":"token-abc"}`, val)

		// Another key
		err = store.Set(ctx, "dataset_token", `{"dataset":"ds1"}`)
		require.NoError(t, err)

		val1, err := store.Get(ctx, "user_token")
		require.NoError(t, err)
		require.Equal(t, `{"access_token":"token-abc"}`, val1)

		val2, err := store.Get(ctx, "dataset_token")
		require.NoError(t, err)
		require.Equal(t, `{"dataset":"ds1"}`, val2)
	})

	t.Run("Delete removes key from file", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "tokenstore-fsz-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		u := urlz.MustParse("file://" + filepath.Join(tmpDir, "tokens.json"))
		store, err := tokenstore.NewFsz(u)
		require.NoError(t, err)

		err = store.Set(ctx, "temp_key", "temp_value")
		require.NoError(t, err)

		err = store.Delete(ctx, "temp_key")
		require.NoError(t, err)

		val, err := store.Get(ctx, "temp_key")
		require.NoError(t, err)
		require.Empty(t, val)

		// Deleting non-existent key does not error
		err = store.Delete(ctx, "nonexistent")
		require.NoError(t, err)
	})

	t.Run("Corrupted JSON recovers gracefully with empty map", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "tokenstore-fsz-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		filePath := filepath.Join(tmpDir, "tokens.json")
		err = os.WriteFile(filePath, []byte("{not-valid-json"), 0600)
		require.NoError(t, err)

		u := urlz.MustParse("file://" + filePath)
		store, err := tokenstore.NewFsz(u)
		require.NoError(t, err)

		val, err := store.Get(ctx, "any_key")
		require.NoError(t, err)
		require.Empty(t, val)

		// Can write new value over corrupted file
		err = store.Set(ctx, "recovered", "ok")
		require.NoError(t, err)

		val, err = store.Get(ctx, "recovered")
		require.NoError(t, err)
		require.Equal(t, "ok", val)
	})

	t.Run("Concurrent Set and Get across goroutines", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "tokenstore-fsz-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		u := urlz.MustParse("file://" + filepath.Join(tmpDir, "tokens.json"))
		store, err := tokenstore.NewFsz(u)
		require.NoError(t, err)

		var wg sync.WaitGroup
		concurrency := 10
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				k := fmt.Sprintf("key_%d", idx)
				v := fmt.Sprintf("val_%d", idx)
				_ = store.Set(context.Background(), k, v)
				got, _ := store.Get(context.Background(), k)
				if got != "" && got != v {
					t.Errorf("expected %s or empty, got %s", v, got)
				}
			}(i)
		}
		wg.Wait()
	})
}
