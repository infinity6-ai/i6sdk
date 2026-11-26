package client_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/i6sdk/sdkgo/client"
	"github.com/stretchr/testify/require"
)

type customTokenStore struct {
	val string
}

func (c *customTokenStore) Get(_ context.Context, key string) (string, error) {
	return c.val, nil
}

func (c *customTokenStore) Set(_ context.Context, key string, value string) error {
	c.val = value
	return nil
}

func (c *customTokenStore) Delete(_ context.Context, key string) error {
	c.val = ""
	return nil
}

func TestUnitClientConstruction(t *testing.T) {
	type testScenario struct {
		opts            []client.Option
		envVars         map[string]string
		wantBaseUrl     string
		wantToken       string
		wantMemoryCache bool
		wantCacheScheme string
		wantErr         string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()

		for k, v := range s.envVars {
			t.Setenv(k, v)
		}

		c, err := client.NewClient(s.opts...)
		if s.wantErr != "" {
			require.Error(t, err)
			require.Contains(t, err.Error(), s.wantErr)
			return
		}

		require.NoError(t, err)
		require.NotNil(t, c)
		require.NotNil(t, c.Apis())
		require.NotNil(t, c.Funcs())
		require.NotNil(t, c.HttpClient())
		require.NotNil(t, c.TokenStore())

		if s.wantBaseUrl != "" {
			require.Equal(t, s.wantBaseUrl, c.BaseUrl().String())
		}
		if s.wantToken != "" {
			require.Equal(t, s.wantToken, c.InitialToken())
		}
		require.Equal(t, s.wantMemoryCache, c.IsMemoryCache())
		if s.wantCacheScheme != "" {
			require.NotNil(t, c.CacheUrl())
			require.Equal(t, s.wantCacheScheme, c.CacheUrl().Scheme)
		} else if s.wantMemoryCache {
			require.Nil(t, c.CacheUrl())
		}
	}

	t.Run("Default client construction without options", func(t *testing.T) {
		check(t, testScenario{
			wantBaseUrl:     client.DefaultBaseUrl,
			wantMemoryCache: false,
			wantCacheScheme: "file",
		})
	})

	t.Run("Client construction with explicit token and base url options", func(t *testing.T) {
		check(t, testScenario{
			opts: []client.Option{
				client.WithToken("explicit-token-123"),
				client.WithBaseUrl("https://custom.endpoint.com"),
			},
			wantBaseUrl:     "https://custom.endpoint.com",
			wantToken:       "explicit-token-123",
			wantMemoryCache: false,
			wantCacheScheme: "file",
		})
	})

	t.Run("Client construction with environment variable fallbacks", func(t *testing.T) {
		check(t, testScenario{
			envVars: map[string]string{
				client.EnvToken:   "env-token-456",
				client.EnvBaseUrl: "https://env.endpoint.com",
			},
			wantBaseUrl:     "https://env.endpoint.com",
			wantToken:       "env-token-456",
			wantMemoryCache: false,
			wantCacheScheme: "file",
		})
	})

	t.Run("Options take precedence over environment variables", func(t *testing.T) {
		check(t, testScenario{
			envVars: map[string]string{
				client.EnvToken:   "env-token-ignore",
				client.EnvBaseUrl: "https://env-ignore.com",
			},
			opts: []client.Option{
				client.WithToken("priority-token"),
				client.WithBaseUrl("https://priority-endpoint.com"),
			},
			wantBaseUrl:     "https://priority-endpoint.com",
			wantToken:       "priority-token",
			wantMemoryCache: false,
			wantCacheScheme: "file",
		})
	})

	t.Run("In-memory cache option enables memoryCache flag", func(t *testing.T) {
		check(t, testScenario{
			opts: []client.Option{
				client.WithMemoryCache(),
			},
			wantBaseUrl:     client.DefaultBaseUrl,
			wantMemoryCache: true,
			wantCacheScheme: "",
		})
	})

	t.Run("Custom cache directory option sets file scheme url", func(t *testing.T) {
		check(t, testScenario{
			opts: []client.Option{
				client.WithCacheDir("/tmp/custom-i6-cache"),
			},
			wantBaseUrl:     client.DefaultBaseUrl,
			wantMemoryCache: false,
			wantCacheScheme: "file",
		})
	})

	t.Run("Custom cache url option configures exact url", func(t *testing.T) {
		customUrl := urlz.MustParse("file:///var/custom/cache")
		check(t, testScenario{
			opts: []client.Option{
				client.WithCacheUrl(customUrl),
			},
			wantBaseUrl:     client.DefaultBaseUrl,
			wantMemoryCache: false,
			wantCacheScheme: "file",
		})
	})

	t.Run("Nil http client returns error", func(t *testing.T) {
		check(t, testScenario{
			opts: []client.Option{
				client.WithHttpClient(nil),
			},
			wantErr: "http client cannot be nil",
		})
	})

	t.Run("Nil token store returns error", func(t *testing.T) {
		check(t, testScenario{
			opts: []client.Option{
				client.WithTokenStore(nil),
			},
			wantErr: "token store cannot be nil",
		})
	})

	t.Run("Custom token store is successfully registered", func(t *testing.T) {
		myStore := &customTokenStore{val: "secret"}
		c, err := client.New(client.WithTokenStore(myStore))
		require.NoError(t, err)
		require.Equal(t, myStore, c.TokenStore())

		val, err := c.TokenStore().Get(t.Context(), "key")
		require.NoError(t, err)
		require.Equal(t, "secret", val)
	})

	t.Run("Custom http client is successfully registered", func(t *testing.T) {
		customHttp := &http.Client{Timeout: 10 * time.Second}
		c, err := client.New(client.WithHttpClient(customHttp))
		require.NoError(t, err)
		require.Equal(t, customHttp, c.HttpClient())
	})
}
