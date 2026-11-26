package client

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/infinity6-ai/gox/commonz/urlz"
)

// Option configures an I6Client.
type Option func(*I6Client) error

// WithToken sets the initial bearer access token.
func WithToken(token string) Option {
	return func(c *I6Client) error {
		c.initialToken = token
		return nil
	}
}

// WithBaseUrl sets the base API URL.
func WithBaseUrl(rawUrl string) Option {
	return func(c *I6Client) error {
		u, err := urlz.Parse(rawUrl)
		if err != nil {
			return fmt.Errorf("failed to parse base url: %w", err)
		}
		c.baseUrl = u
		return nil
	}
}

// WithHttpClient overrides the underlying HTTP client.
func WithHttpClient(httpClient *http.Client) Option {
	return func(c *I6Client) error {
		if httpClient == nil {
			return fmt.Errorf("http client cannot be nil")
		}
		c.httpClient = httpClient
		return nil
	}
}

// WithCacheDir configures the local cache directory for SQLite bundles as a file:// URL.
func WithCacheDir(dir string) Option {
	return func(c *I6Client) error {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return fmt.Errorf("failed to resolve absolute cache directory: %w", err)
		}
		u, err := urlz.Parse("file://" + absDir)
		if err != nil {
			return fmt.Errorf("failed to parse cache directory as url: %w", err)
		}
		c.cacheUrl = u
		c.memoryCache = false
		return nil
	}
}

// WithCacheUrl sets a custom fsz URL for the bundle cache.
func WithCacheUrl(u *urlz.Url) Option {
	return func(c *I6Client) error {
		if u == nil {
			return fmt.Errorf("cache url cannot be nil")
		}
		c.cacheUrl = u
		c.memoryCache = false
		return nil
	}
}

// WithMemoryCache configures in-memory execution for SQLite bundles.
func WithMemoryCache() Option {
	return func(c *I6Client) error {
		c.memoryCache = true
		c.cacheUrl = nil
		return nil
	}
}

// WithTokenStore configures a custom TokenStore implementation.
func WithTokenStore(store TokenStore) Option {
	return func(c *I6Client) error {
		if store == nil {
			return fmt.Errorf("token store cannot be nil")
		}
		c.tokenStore = store
		return nil
	}
}
