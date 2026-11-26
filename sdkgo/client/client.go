package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
)

const (
	DefaultBaseUrl = "https://api.infinity6.ai"
	EnvToken       = "I6_TOKEN"
	EnvBaseUrl     = "I6_BASE_URL"
)

// I6Client is the primary client for the Infinity6 SDK.
type I6Client struct {
	tokenStore   TokenStore
	httpClient   *http.Client
	baseUrl      *urlz.Url
	cacheUrl     *urlz.Url
	memoryCache  bool
	initialToken string

	apis  *apicollectionz.SdkApis
	funcs *apicollectionz.SdkFuncs
}

// New constructs an I6Client with the provided options. Alias of NewClient.
func New(opts ...Option) (*I6Client, error) {
	return NewClient(opts...)
}

// NewClient constructs an I6Client with the provided options.
func NewClient(opts ...Option) (*I6Client, error) {
	c := &I6Client{
		httpClient: &http.Client{Timeout: 60 * time.Second},
		tokenStore: newMemoryTokenStore(),
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
	}

	// Fallback to environment variables if not set via options
	if c.baseUrl == nil {
		rawBase := os.Getenv(EnvBaseUrl)
		if rawBase == "" {
			rawBase = DefaultBaseUrl
		}
		u, err := urlz.Parse(rawBase)
		if err != nil {
			return nil, fmt.Errorf("failed to parse base url %q: %w", rawBase, err)
		}
		c.baseUrl = u
	}

	if c.initialToken == "" {
		c.initialToken = os.Getenv(EnvToken)
	}

	// Default cache directory if not explicitly set and not in memory mode
	if !c.memoryCache && c.cacheUrl == nil {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = os.TempDir()
		}
		defaultCache := filepath.Join(homeDir, ".cache", "i6", "sdkgo", "fs")
		c.cacheUrl = urlz.MustParse("file://" + filepath.Clean(defaultCache))
	}

	// Initialize SdkApis with initial fallback invoker
	c.apis = &apicollectionz.SdkApis{
		Invoker: func(ctx context.Context, api apiz.Api) (int, error) {
			return http.StatusNotImplemented, fmt.Errorf("transport invoker not yet wired: step 4.3")
		},
	}

	// Initialize SdkFuncs with initial fallback invoker
	c.funcs = &apicollectionz.SdkFuncs{
		Invoker: func(ctx context.Context, fn sdkfuncz.Func) (int, error) {
			return http.StatusNotImplemented, fmt.Errorf("sdkfunc invoker not yet wired: step 4.7")
		},
	}

	return c, nil
}

// Apis returns the typed SdkApis client containing all API endpoints.
func (c *I6Client) Apis() *apicollectionz.SdkApis {
	return c.apis
}

// Funcs returns the typed SdkFuncs client containing all workflow functions.
func (c *I6Client) Funcs() *apicollectionz.SdkFuncs {
	return c.funcs
}

// BaseUrl returns the configured base URL.
func (c *I6Client) BaseUrl() *urlz.Url {
	return c.baseUrl
}

// HttpClient returns the configured HTTP client.
func (c *I6Client) HttpClient() *http.Client {
	return c.httpClient
}

// CacheUrl returns the configured fsz cache URL.
func (c *I6Client) CacheUrl() *urlz.Url {
	return c.cacheUrl
}

// IsMemoryCache returns true if in-memory cache is active.
func (c *I6Client) IsMemoryCache() bool {
	return c.memoryCache
}

// TokenStore returns the active token persistence store.
func (c *I6Client) TokenStore() TokenStore {
	return c.tokenStore
}

// InitialToken returns the initial token provided to the client.
func (c *I6Client) InitialToken() string {
	return c.initialToken
}
