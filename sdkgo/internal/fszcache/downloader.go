package fszcache

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/gox/fsz/fsz"
)

// Downloader manages downloading remote assets into the versioned Fs store with ETag caching.
type Downloader struct {
	fs         *Fs
	httpClient *http.Client
	mu         sync.Mutex
}

// NewDownloader creates a new Downloader using the provided Fs store and optional custom HTTP client.
func NewDownloader(fs *Fs, httpClient ...*http.Client) *Downloader {
	client := http.DefaultClient
	if len(httpClient) > 0 && httpClient[0] != nil {
		client = httpClient[0]
	}
	return &Downloader{
		fs:         fs,
		httpClient: client,
	}
}

// Fs returns the underlying versioned Fs store.
func (d *Downloader) Fs() *Fs {
	return d.fs
}

// Update downloads or revalidates a file by name, returning the local blob.bin path.
// If minPreloadTime (in seconds) has not elapsed since checked_at, network revalidation is bypassed.
func (d *Downloader) Update(
	ctx context.Context,
	name string,
	rawUrl string,
	minPreloadTime float64,
	timeout ...time.Duration,
) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	current, err := d.fs.Get(ctx, name)
	if err != nil {
		return "", fmt.Errorf("failed to get cached version for %s: %w", name, err)
	}

	// Check if cached version is within preload TTL
	if minPreloadTime > 0 && current != nil && current.CheckedAt != "" {
		if t, err := parseIsoTimestamp(current.CheckedAt); err == nil {
			if time.Since(t).Seconds() < minPreloadTime {
				return d.fs.Resolve(ctx, name, current.Version), nil
			}
		}
	}

	return d.fetchAndUpdate(ctx, name, rawUrl, current, timeout...)
}

func (d *Downloader) fetchAndUpdate(
	ctx context.Context,
	name string,
	rawUrl string,
	current *FileRef,
	timeout ...time.Duration,
) (string, error) {
	reqCtx := ctx
	if len(timeout) > 0 && timeout[0] > 0 {
		var cancel context.CancelFunc
		reqCtx, cancel = context.WithTimeout(ctx, timeout[0])
		defer cancel()
	}

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, rawUrl, nil)
	if err != nil {
		if current != nil {
			return d.fs.Resolve(ctx, name, current.Version), nil
		}
		return "", fmt.Errorf("failed to create http request for %s: %w", rawUrl, err)
	}

	req.Header.Set("User-Agent", "i6sdk-go/0.0.18")
	if current != nil && current.Etag != "" {
		req.Header.Set("If-None-Match", current.Etag)
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		// Network failure: fallback to stale cached version if available
		if current != nil {
			return d.fs.Resolve(ctx, name, current.Version), nil
		}
		return "", fmt.Errorf("download failed for %s (%s): %w", name, rawUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified && current != nil {
		d.fs.Touch(ctx, name)
		return d.fs.Resolve(ctx, name, current.Version), nil
	}

	if resp.StatusCode != http.StatusOK {
		if current != nil {
			return d.fs.Resolve(ctx, name, current.Version), nil
		}
		return "", fmt.Errorf("download failed for %s (%s): HTTP %d", name, rawUrl, resp.StatusCode)
	}

	respEtag := resp.Header.Get("ETag")
	newRef, err := d.fs.Create(ctx, name, respEtag)
	if err != nil {
		return "", fmt.Errorf("failed to create version for %s: %w", name, err)
	}

	targetUrl := d.fs.ResolveUrl(ctx, name, newRef.Version)
	tmpUrl, err := urlz.Parse(targetUrl.String() + ".tmp")
	if err != nil {
		_ = d.fs.Discard(ctx, name, newRef.Version)
		return "", fmt.Errorf("failed to create tmp url for %s: %w", targetUrl, err)
	}

	if err := fsz.Upload(ctx, tmpUrl, nil, resp.Body); err != nil {
		_ = d.fs.Discard(ctx, name, newRef.Version)
		return "", fmt.Errorf("failed to upload stream to %s: %w", tmpUrl, err)
	}

	if err := fsz.Move(ctx, tmpUrl, targetUrl); err != nil {
		_ = d.fs.Discard(ctx, name, newRef.Version)
		return "", fmt.Errorf("failed to atomically move %s to %s: %w", tmpUrl, targetUrl, err)
	}

	released, err := d.fs.Release(ctx, name, newRef.Version)
	if err != nil {
		return "", fmt.Errorf("failed to release version %s for %s: %w", newRef.Version, name, err)
	}

	return d.fs.Resolve(ctx, name, released.Version), nil
}

func parseIsoTimestamp(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02T15:04:05.000Z",
		time.RFC3339Nano,
		time.RFC3339,
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse ISO timestamp %q", s)
}
