package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/gox/httpz/httpzclient"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthlogout"
	"github.com/infinity6-ai/i6sdk/sdkgo/client"
	"github.com/infinity6-ai/i6sdk/sdkgo/config"
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
		opts            config.Options
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

		c, err := client.New(s.opts)
		if s.wantErr != "" {
			require.Error(t, err)
			require.Contains(t, err.Error(), s.wantErr)
			return
		}

		require.NoError(t, err)
		require.NotNil(t, c)
		require.NotNil(t, c.Apis())
		require.NotNil(t, c.Funcs())
		require.NotNil(t, c.DecsuiteClient())
		require.NotNil(t, c.TfinferClient())
		require.NotNil(t, c.UserzClient())
		require.NotNil(t, c.TokenStore())

		if s.wantDecsuite != "" {
			require.Equal(t, s.wantDecsuite, c.DecsuiteUrl().String())
			require.Equal(t, s.wantDecsuite, c.BaseUrl().String())
		}
		if s.wantTfinfer != "" {
			require.Equal(t, s.wantTfinfer, c.TfinferUrl().String())
		}
		if s.wantUserz != "" {
			require.Equal(t, s.wantUserz, c.UserzUrl().String())
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

	t.Run("Default client construction with production URLs", func(t *testing.T) {
		check(t, testScenario{
			wantDecsuite:    config.DefaultUrlOptions.Decsuite,
			wantTfinfer:     config.DefaultUrlOptions.Tfinfer,
			wantUserz:       config.DefaultUrlOptions.Userz,
			wantMemoryCache: false,
			wantCacheScheme: "file",
		})
	})

	t.Run("Client construction with SandboxUrlOptions preset", func(t *testing.T) {
		check(t, testScenario{
			opts: config.Options{
				Urls: config.SandboxUrlOptions,
			},
			wantDecsuite:    config.SandboxUrlOptions.Decsuite,
			wantTfinfer:     config.SandboxUrlOptions.Tfinfer,
			wantUserz:       config.SandboxUrlOptions.Userz,
			wantMemoryCache: false,
			wantCacheScheme: "file",
		})
	})

	t.Run("Client construction with StagingUrlOptions preset", func(t *testing.T) {
		check(t, testScenario{
			opts: config.Options{
				Urls: config.StagingUrlOptions,
			},
			wantDecsuite:    config.StagingUrlOptions.Decsuite,
			wantTfinfer:     config.StagingUrlOptions.Tfinfer,
			wantUserz:       config.StagingUrlOptions.Userz,
			wantMemoryCache: false,
			wantCacheScheme: "file",
		})
	})

	t.Run("Client construction with explicit UrlOptions", func(t *testing.T) {
		check(t, testScenario{
			opts: config.Options{
				Token: "my-token",
				Urls: config.UrlOptions{
					Decsuite: "https://custom.decsuite.ai",
					Tfinfer:  "https://custom.tfinfer.ai",
					Userz:    "https://custom.userz.ai",
				},
			},
			wantDecsuite:    "https://custom.decsuite.ai",
			wantTfinfer:     "https://custom.tfinfer.ai",
			wantUserz:       "https://custom.userz.ai",
			wantToken:       "my-token",
			wantMemoryCache: false,
			wantCacheScheme: "file",
		})
	})

	t.Run("In-memory cache option enables memoryCache flag", func(t *testing.T) {
		check(t, testScenario{
			opts: config.Options{
				Cache: config.CacheOptions{Memory: true},
			},
			wantDecsuite:    config.DefaultUrlOptions.Decsuite,
			wantTfinfer:     config.DefaultUrlOptions.Tfinfer,
			wantUserz:       config.DefaultUrlOptions.Userz,
			wantMemoryCache: true,
			wantCacheScheme: "",
		})
	})

	t.Run("Custom cache directory option sets file scheme url", func(t *testing.T) {
		check(t, testScenario{
			opts: config.Options{
				Cache: config.CacheOptions{Dir: "/tmp/custom-i6-cache"},
			},
			wantDecsuite:    config.DefaultUrlOptions.Decsuite,
			wantTfinfer:     config.DefaultUrlOptions.Tfinfer,
			wantUserz:       config.DefaultUrlOptions.Userz,
			wantMemoryCache: false,
			wantCacheScheme: "file",
		})
	})

	t.Run("Custom cache url option configures exact url", func(t *testing.T) {
		customUrl := urlz.MustParse("file:///var/custom/cache")
		check(t, testScenario{
			opts: config.Options{
				Cache: config.CacheOptions{Url: customUrl},
			},
			wantDecsuite:    config.DefaultUrlOptions.Decsuite,
			wantTfinfer:     config.DefaultUrlOptions.Tfinfer,
			wantUserz:       config.DefaultUrlOptions.Userz,
			wantMemoryCache: false,
			wantCacheScheme: "file",
		})
	})

	t.Run("Custom token store is successfully registered", func(t *testing.T) {
		myStore := &customTokenStore{val: "secret"}
		c, err := client.New(config.Options{TokenStore: myStore})
		require.NoError(t, err)
		require.Equal(t, myStore, c.TokenStore())

		val, err := c.TokenStore().Get(t.Context(), "key")
		require.NoError(t, err)
		require.Equal(t, "secret", val)
	})

	t.Run("Path routing dispatches to corresponding service clients", func(t *testing.T) {
		c, err := client.New(config.Options{})
		require.NoError(t, err)

		userClient := c.ClientForPath("/api/userz/auth/me")
		require.Equal(t, c.UserzClient(), userClient)

		tfClient1 := c.ClientForPath("/api/metadata/relevance-fashion/v1")
		require.Equal(t, c.TfinferClient(), tfClient1)

		tfClient2 := c.ClientForPath("/api/model/dataset/ds1/model/v1/predict")
		require.Equal(t, c.TfinferClient(), tfClient2)

		decClient := c.ClientForPath("/api/ds/dataset/ds1/store/list")
		require.Equal(t, c.DecsuiteClient(), decClient)
	})

	t.Run("Custom httpzclient Options are propagated", func(t *testing.T) {
		customCalled := false
		c, err := client.New(config.Options{
			Http: httpzclient.Options{
				GetClient: func(ctx context.Context) *http.Client {
					customCalled = true
					return http.DefaultClient
				},
			},
		})
		require.NoError(t, err)
		require.True(t, customCalled)
		require.NotNil(t, c.DecsuiteClient())
	})

	t.Run("Apis invoker routes and executes API via apiclientz with bearer auth", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/userz/auth/logout", r.URL.Path)
			require.Equal(t, "Bearer my-secret-token", r.Header.Get("Authorization"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		}))
		defer ts.Close()

		c, err := client.New(config.Options{
			Token: "my-secret-token",
			Urls: config.UrlOptions{
				Userz: ts.URL,
			},
		})
		require.NoError(t, err)

		api := &userapiauthlogout.Api{
			Req: &userapiauthlogout.Req{},
		}
		status, err := c.Apis().UserAuthLogout(t.Context(), api)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)
	})
}
