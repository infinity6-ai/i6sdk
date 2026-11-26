package fszcache

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
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

	current, _ := d.fs.Get(name)

	// Check if cached version is within preload TTL
	if minPreloadTime > 0 && current != nil && current.CheckedAt != "" {
		if t, err := parseIsoTimestamp(current.CheckedAt); err == nil {
			if time.Since(t).Seconds() < minPreloadTime {
				return d.fs.Resolve(name, current.Version), nil
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
			return d.fs.Resolve(name, current.Version), nil
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
			return d.fs.Resolve(name, current.Version), nil
		}
		return "", fmt.Errorf("download failed for %s (%s): %w", name, rawUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified && current != nil {
		d.fs.Touch(name)
		return d.fs.Resolve(name, current.Version), nil
	}

	if resp.StatusCode != http.StatusOK {
		if current != nil {
			return d.fs.Resolve(name, current.Version), nil
		}
		return "", fmt.Errorf("download failed for %s (%s): HTTP %d", name, rawUrl, resp.StatusCode)
	}

	respEtag := resp.Header.Get("ETag")
	newRef, err := d.fs.Create(name, respEtag)
	if err != nil {
		return "", fmt.Errorf("failed to create version for %s: %w", name, err)
	}

	targetPath := d.fs.Resolve(name, newRef.Version)
	tmpPath := targetPath + ".tmp"

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		_ = d.fs.Discard(name, newRef.Version)
		return "", fmt.Errorf("failed to create directory for %s: %w", targetPath, err)
	}

	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		_ = d.fs.Discard(name, newRef.Version)
		return "", fmt.Errorf("failed to create tmp file %s: %w", tmpPath, err)
	}

	buf := make([]byte, 64*1024)
	if _, err := io.CopyBuffer(tmpFile, resp.Body, buf); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
		_ = d.fs.Discard(name, newRef.Version)
		return "", fmt.Errorf("failed to stream download into %s: %w", tmpPath, err)
	}
	_ = tmpFile.Close()

	if err := os.Rename(tmpPath, targetPath); err != nil {
		_ = os.Remove(tmpPath)
		_ = d.fs.Discard(name, newRef.Version)
		return "", fmt.Errorf("failed to atomically rename %s to %s: %w", tmpPath, targetPath, err)
	}

	released, err := d.fs.Release(name, newRef.Version)
	if err != nil {
		return "", fmt.Errorf("failed to release version %s for %s: %w", newRef.Version, name, err)
	}

	return d.fs.Resolve(name, released.Version), nil
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
