package fszcache_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/i6sdk/sdkgo/internal/fszcache"
	"github.com/stretchr/testify/require"
)

func TestUnitFsLifecycle(t *testing.T) {
	ctx := context.Background()
	tmpDir, err := os.MkdirTemp("", "i6-fszcache-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	baseUrl := urlz.MustParse("file://" + filepath.Clean(tmpDir))
	fs, err := fszcache.NewFs(baseUrl)
	require.NoError(t, err)
	require.Equal(t, baseUrl.String(), fs.BaseUrl().String())

	t.Run("DefaultBaseUrl resolves successfully", func(t *testing.T) {
		def, err := fszcache.DefaultBaseUrl()
		require.NoError(t, err)
		require.NotNil(t, def)
		require.Equal(t, "file", def.Scheme)
		require.Contains(t, def.String(), "fs")
	})

	t.Run("Create version and write etag", func(t *testing.T) {
		ref, err := fs.Create(ctx, "model/catalog", `"etag-1"`)
		require.NoError(t, err)
		require.Equal(t, "model/catalog", ref.Name)
		require.NotEmpty(t, ref.Version)
		require.Equal(t, `"etag-1"`, ref.Etag)

		// Create again with same etag reuses version
		ref2, err := fs.Create(ctx, "model/catalog", `"etag-1"`)
		require.NoError(t, err)
		require.Equal(t, ref.Version, ref2.Version)
	})

	t.Run("Resolve, Write blob, and Release version", func(t *testing.T) {
		ref, err := fs.Create(ctx, "dash/funnel", `"etag-v1"`)
		require.NoError(t, err)

		blobPath := fs.Resolve("dash/funnel", ref.Version)
		require.Contains(t, blobPath, filepath.Join("v"+ref.Version, "blob.bin"))

		// Write blob
		err = os.WriteFile(blobPath, []byte("sqlite-binary-data-v1"), 0644)
		require.NoError(t, err)

		released, err := fs.Release(ctx, "dash/funnel", ref.Version)
		require.NoError(t, err)
		require.Equal(t, ref.Version, released.Version)
		require.NotEmpty(t, released.CheckedAt)

		// Exists is now true
		require.True(t, fs.Exists(ctx, "dash/funnel"))

		// Get without version returns released version
		current, err := fs.Get(ctx, "dash/funnel")
		require.NoError(t, err)
		require.NotNil(t, current)
		require.Equal(t, ref.Version, current.Version)
		require.Equal(t, `"etag-v1"`, current.Etag)
		require.NotEmpty(t, current.CheckedAt)

		// Resolve without version points to released blob
		resolvedPath := fs.Resolve("dash/funnel")
		require.Equal(t, blobPath, resolvedPath)
	})

	t.Run("CleanFile prunes older versions on release", func(t *testing.T) {
		asset := "dash/retention"
		ref1, err := fs.Create(ctx, asset, `"v1"`)
		require.NoError(t, err)
		blob1 := fs.Resolve(asset, ref1.Version)
		err = os.WriteFile(blob1, []byte("v1"), 0644)
		require.NoError(t, err)
		_, err = fs.Release(ctx, asset, ref1.Version)
		require.NoError(t, err)

		time.Sleep(10 * time.Millisecond) // ensure distinct ISO timestamp

		ref2, err := fs.Create(ctx, asset, `"v2"`)
		require.NoError(t, err)
		blob2 := fs.Resolve(asset, ref2.Version)
		err = os.WriteFile(blob2, []byte("v2"), 0644)
		require.NoError(t, err)
		_, err = fs.Release(ctx, asset, ref2.Version)
		require.NoError(t, err)

		// v1 directory should have been cleaned
		v1Dir := filepath.Dir(blob1)
		_, err = os.Stat(v1Dir)
		require.True(t, os.IsNotExist(err), "v1 directory should be removed after v2 release")

		// v2 directory exists
		v2Dir := filepath.Dir(blob2)
		_, err = os.Stat(v2Dir)
		require.NoError(t, err)
	})

	t.Run("Touch updates checkedat timestamp", func(t *testing.T) {
		ts1 := fs.Touch(ctx, "asset1")
		require.NotEmpty(t, ts1)

		custom := "2026-11-26T15:00:00.000Z"
		ts2 := fs.Touch(ctx, "asset1", custom)
		require.Equal(t, custom, ts2)
	})

	t.Run("Delete removes asset directory", func(t *testing.T) {
		ref, err := fs.Create(ctx, "to_delete", "etag")
		require.NoError(t, err)
		blob := fs.Resolve("to_delete", ref.Version)
		_ = os.WriteFile(blob, []byte("data"), 0644)
		_, _ = fs.Release(ctx, "to_delete", ref.Version)
		require.True(t, fs.Exists(ctx, "to_delete"))

		err = fs.Delete(ctx, "to_delete")
		require.NoError(t, err)
		require.False(t, fs.Exists(ctx, "to_delete"))
	})
}

func TestUnitDownloader(t *testing.T) {
	ctx := context.Background()
	tmpDir, err := os.MkdirTemp("", "i6-downloader-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	baseUrl := urlz.MustParse("file://" + filepath.Clean(tmpDir))
	fs, err := fszcache.NewFs(baseUrl)
	require.NoError(t, err)
	downloader := fszcache.NewDownloader(fs)

	t.Run("First download streams content and releases version", func(t *testing.T) {
		var serverHits int32
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&serverHits, 1)
			w.Header().Set("ETag", `"etag-initial"`)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("sqlite-initial-binary"))
		}))
		defer ts.Close()

		filePath, err := downloader.Update(ctx, "dash/metrics", ts.URL, 300.0)
		require.NoError(t, err)
		require.Equal(t, int32(1), atomic.LoadInt32(&serverHits))

		content, err := os.ReadFile(filePath)
		require.NoError(t, err)
		require.Equal(t, "sqlite-initial-binary", string(content))

		ref, err := fs.Get(ctx, "dash/metrics")
		require.NoError(t, err)
		require.NotNil(t, ref)
		require.Equal(t, `"etag-initial"`, ref.Etag)
		require.NotEmpty(t, ref.CheckedAt)
	})

	t.Run("TTL bypass returns cached path with zero network calls", func(t *testing.T) {
		var serverHits int32
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&serverHits, 1)
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		// CheckedAt is fresh (< 300 seconds), so TTL bypass triggers
		filePath, err := downloader.Update(ctx, "dash/metrics", ts.URL, 300.0)
		require.NoError(t, err)
		require.Equal(t, int32(0), atomic.LoadInt32(&serverHits), "Network should not be touched within preload TTL")

		content, err := os.ReadFile(filePath)
		require.NoError(t, err)
		require.Equal(t, "sqlite-initial-binary", string(content))
	})

	t.Run("HTTP 304 Not Modified revalidates cache and updates checkedat", func(t *testing.T) {
		var serverHits int32
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&serverHits, 1)
			require.Equal(t, `"etag-initial"`, r.Header.Get("If-None-Match"))
			w.WriteHeader(http.StatusNotModified)
		}))
		defer ts.Close()

		// Force minPreloadTime = 0 to trigger remote check
		beforeRef, _ := fs.Get(ctx, "dash/metrics")
		time.Sleep(10 * time.Millisecond)

		filePath, err := downloader.Update(ctx, "dash/metrics", ts.URL, 0)
		require.NoError(t, err)
		require.Equal(t, int32(1), atomic.LoadInt32(&serverHits))

		afterRef, _ := fs.Get(ctx, "dash/metrics")
		require.Equal(t, beforeRef.Version, afterRef.Version)
		require.GreaterOrEqual(t, afterRef.CheckedAt, beforeRef.CheckedAt)

		content, err := os.ReadFile(filePath)
		require.NoError(t, err)
		require.Equal(t, "sqlite-initial-binary", string(content))
	})

	t.Run("HTTP 200 with new ETag downloads new content and deletes older version", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("ETag", `"etag-updated"`)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("sqlite-updated-binary"))
		}))
		defer ts.Close()

		oldRef, _ := fs.Get(ctx, "dash/metrics")

		time.Sleep(10 * time.Millisecond)
		filePath, err := downloader.Update(ctx, "dash/metrics", ts.URL, 0)
		require.NoError(t, err)

		content, err := os.ReadFile(filePath)
		require.NoError(t, err)
		require.Equal(t, "sqlite-updated-binary", string(content))

		newRef, _ := fs.Get(ctx, "dash/metrics")
		require.NotEqual(t, oldRef.Version, newRef.Version)
		require.Equal(t, `"etag-updated"`, newRef.Etag)
	})

	t.Run("Server error falls back to existing cached version", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		filePath, err := downloader.Update(ctx, "dash/metrics", ts.URL, 0)
		require.NoError(t, err)

		content, err := os.ReadFile(filePath)
		require.NoError(t, err)
		require.Equal(t, "sqlite-updated-binary", string(content))
	})
}
