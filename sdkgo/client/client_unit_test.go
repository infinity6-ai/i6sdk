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
		wantDecsuite    string
		wantTfinfer     string
		wantUserz       string
		wantToken       string
		wantMemoryCache bool
		wantCacheScheme string
		wantErr         string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()

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

		if s.wantDecsuite != "" {
			require.Equal(t, s.wantDecsuite, c.DecsuiteUrl().String())
			require.Equal(t, s.wantDecsuite, c.BaseUrl().String())
			require.Equal(t, s.wantDecsuite, c.ServiceUrl("decsuite").String())
			require.Equal(t, s.wantDecsuite, c.ServiceUrl("other").String())
		}
		if s.wantTfinfer != "" {
			require.Equal(t, s.wantTfinfer, c.TfinferUrl().String())
			require.Equal(t, s.wantTfinfer, c.ServiceUrl("tfinfer").String())
		}
		if s.wantUserz != "" {
			require.Equal(t, s.wantUserz, c.UserzUrl().String())
			require.Equal(t, s.wantUserz, c.ServiceUrl("userz").String())
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

	t.Run("Default client construction with sandbox URLs", func(t *testing.T) {
		check(t, testScenario{
			wantDecsuite:    client.DefaultDecsuiteUrl,
			wantTfinfer:     client.DefaultTfinferUrl,
			wantUserz:       client.DefaultUserzUrl,
			wantMemoryCache: false,
			wantCacheScheme: "file",
		})
	})

	t.Run("Client construction with explicit 3 URLs", func(t *testing.T) {
		check(t, testScenario{
			opts: []client.Option{
				client.WithToken("my-token"),
				client.WithDecsuiteUrl("https://custom.decsuite.ai"),
				client.WithTfinferUrl("https://custom.tfinfer.ai"),
				client.WithUserzUrl("https://custom.userz.ai"),
			},
			wantDecsuite:    "https://custom.decsuite.ai",
			wantTfinfer:     "https://custom.tfinfer.ai",
			wantUserz:       "https://custom.userz.ai",
			wantToken:       "my-token",
			wantMemoryCache: false,
			wantCacheScheme: "file",
		})
	})

	t.Run("Client construction with WithServiceUrls map", func(t *testing.T) {
		check(t, testScenario{
			opts: []client.Option{
				client.WithServiceUrls(map[string]string{
					"decsuite": "https://map.decsuite.ai",
					"tfinfer":  "https://map.tfinfer.ai",
					"userz":    "https://map.userz.ai",
				}),
			},
			wantDecsuite:    "https://map.decsuite.ai",
			wantTfinfer:     "https://map.tfinfer.ai",
			wantUserz:       "https://map.userz.ai",
			wantMemoryCache: false,
			wantCacheScheme: "file",
		})
	})

	t.Run("Client construction with WithStaging preset", func(t *testing.T) {
		check(t, testScenario{
			opts: []client.Option{
				client.WithStaging(),
			},
			wantDecsuite:    client.StagingDecsuiteUrl,
			wantTfinfer:     client.StagingTfinferUrl,
			wantUserz:       client.StagingUserzUrl,
			wantMemoryCache: false,
			wantCacheScheme: "file",
		})
	})

	t.Run("In-memory cache option enables memoryCache flag", func(t *testing.T) {
		check(t, testScenario{
			opts: []client.Option{
				client.WithMemoryCache(),
			},
			wantDecsuite:    client.DefaultDecsuiteUrl,
			wantMemoryCache: true,
			wantCacheScheme: "",
		})
	})

	t.Run("Custom cache directory option sets file scheme url", func(t *testing.T) {
		check(t, testScenario{
			opts: []client.Option{
				client.WithCacheDir("/tmp/custom-i6-cache"),
			},
			wantDecsuite:    client.DefaultDecsuiteUrl,
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
			wantDecsuite:    client.DefaultDecsuiteUrl,
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

	t.Run("Unknown service name in WithServiceUrls returns error", func(t *testing.T) {
		check(t, testScenario{
			opts: []client.Option{
				client.WithServiceUrls(map[string]string{
					"unknown_srv": "https://api.test",
				}),
			},
			wantErr: "unknown service name",
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
