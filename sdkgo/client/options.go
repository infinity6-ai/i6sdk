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

// WithDecsuiteUrl sets the DecSuite service URL.
func WithDecsuiteUrl(rawUrl string) Option {
	return func(c *I6Client) error {
		u, err := urlz.Parse(rawUrl)
		if err != nil {
			return fmt.Errorf("failed to parse decsuite url: %w", err)
		}
		c.decsuiteUrl = u
		return nil
	}
}

// WithTfinferUrl sets the TensorFlow Inference service URL.
func WithTfinferUrl(rawUrl string) Option {
	return func(c *I6Client) error {
		u, err := urlz.Parse(rawUrl)
		if err != nil {
			return fmt.Errorf("failed to parse tfinfer url: %w", err)
		}
		c.tfinferUrl = u
		return nil
	}
}

// WithUserzUrl sets the Userz authentication service URL.
func WithUserzUrl(rawUrl string) Option {
	return func(c *I6Client) error {
		u, err := urlz.Parse(rawUrl)
		if err != nil {
			return fmt.Errorf("failed to parse userz url: %w", err)
		}
		c.userzUrl = u
		return nil
	}
}

// WithBaseUrl sets the base DecSuite service URL. Alias of WithDecsuiteUrl.
func WithBaseUrl(rawUrl string) Option {
	return WithDecsuiteUrl(rawUrl)
}

// WithServiceUrls configures service URLs by service name ("decsuite", "tfinfer", "userz").
func WithServiceUrls(urls map[string]string) Option {
	return func(c *I6Client) error {
		for name, rawUrl := range urls {
			switch name {
			case "decsuite":
				if err := WithDecsuiteUrl(rawUrl)(c); err != nil {
					return err
				}
			case "tfinfer":
				if err := WithTfinferUrl(rawUrl)(c); err != nil {
					return err
				}
			case "userz":
				if err := WithUserzUrl(rawUrl)(c); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unknown service name %q in service urls", name)
			}
		}
		return nil
	}
}

// WithStaging configures all service URLs to point to Infinity6 Staging.
func WithStaging() Option {
	return WithServiceUrls(map[string]string{
		"decsuite": StagingDecsuiteUrl,
		"tfinfer":  StagingTfinferUrl,
		"userz":    StagingUserzUrl,
	})
}

// WithSandbox configures all service URLs to point to Infinity6 Sandbox.
func WithSandbox() Option {
	return WithServiceUrls(map[string]string{
		"decsuite": DefaultDecsuiteUrl,
		"tfinfer":  DefaultTfinferUrl,
		"userz":    DefaultUserzUrl,
	})
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
