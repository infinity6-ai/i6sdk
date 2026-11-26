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
	// Default Sandbox URLs
	DefaultDecsuiteUrl = "https://decsuite.sandbox.rs.infinity6.ai"
	DefaultTfinferUrl  = "https://tfinfer.sandbox.lb.infinity6.ai"
	DefaultUserzUrl    = "https://userz.sandbox.rs.infinity6.ai"

	// Staging URLs
	StagingDecsuiteUrl = "https://decsuite.staging.rs.infinity6.ai"
	StagingTfinferUrl  = "https://tfinfer.staging.lb.infinity6.ai"
	StagingUserzUrl    = "https://userz.staging.rs.infinity6.ai"
)

// I6Client is the primary client for the Infinity6 SDK.
type I6Client struct {
	tokenStore   TokenStore
	httpClient   *http.Client
	decsuiteUrl  *urlz.Url
	tfinferUrl   *urlz.Url
	userzUrl     *urlz.Url
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

	// 1. Resolve decsuite URL
	if c.decsuiteUrl == nil {
		c.decsuiteUrl = urlz.MustParse(DefaultDecsuiteUrl)
	}

	// 2. Resolve tfinfer URL
	if c.tfinferUrl == nil {
		c.tfinferUrl = urlz.MustParse(DefaultTfinferUrl)
	}

	// 3. Resolve userz URL
	if c.userzUrl == nil {
		c.userzUrl = urlz.MustParse(DefaultUserzUrl)
	}

	// 4. Default cache directory if not explicitly set and not in memory mode
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

// DecsuiteUrl returns the configured DecSuite base URL.
func (c *I6Client) DecsuiteUrl() *urlz.Url {
	return c.decsuiteUrl
}

// TfinferUrl returns the configured TensorFlow Inference base URL.
func (c *I6Client) TfinferUrl() *urlz.Url {
	return c.tfinferUrl
}

// UserzUrl returns the configured Userz authentication base URL.
func (c *I6Client) UserzUrl() *urlz.Url {
	return c.userzUrl
}

// BaseUrl returns the base API URL (alias of DecsuiteUrl).
func (c *I6Client) BaseUrl() *urlz.Url {
	return c.decsuiteUrl
}

// ServiceUrl returns the target URL for a given service name ("decsuite", "tfinfer", "userz").
func (c *I6Client) ServiceUrl(service string) *urlz.Url {
	switch service {
	case "tfinfer":
		if c.tfinferUrl != nil {
			return c.tfinferUrl
		}
	case "userz":
		if c.userzUrl != nil {
			return c.userzUrl
		}
	}
	return c.decsuiteUrl
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
