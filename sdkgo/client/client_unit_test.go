package client_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/gox/httpz/httpzclient"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapistorelist"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthlogout"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashfunnel"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/inferz/inferzrelevancefashionfbt"
	"github.com/infinity6-ai/i6sdk/sdkgo/client"
	"github.com/infinity6-ai/i6sdk/sdkgo/config"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokens"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
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

	t.Run("Apis automatically uses dynamic user token from Tokens manager", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/userz/auth/logout", r.URL.Path)
			require.Equal(t, "Bearer dynamic-user-jwt", r.Header.Get("Authorization"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		}))
		defer ts.Close()

		c, err := client.New(config.Options{
			Urls: config.UrlOptions{
				Userz: ts.URL,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, c.Tokens())

		// Store dynamic user token
		err = c.Tokens().SetUserToken(t.Context(), &tokens.UserToken{
			UserAccessToken: "dynamic-user-jwt",
		})
		require.NoError(t, err)

		api := &userapiauthlogout.Api{
			Req: &userapiauthlogout.Req{},
		}
		status, err := c.Apis().UserAuthLogout(t.Context(), api)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)
	})

	t.Run("Apis automatically uses dynamic dataset token for dataset routes", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/ds/dataset/test-ds/store/list", r.URL.Path)
			require.Equal(t, "Bearer dynamic-dataset-token", r.Header.Get("Authorization"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"payload":{"dataset":"test-ds","stores":[]}}`))
		}))
		defer ts.Close()

		c, err := client.New(config.Options{
			Urls: config.UrlOptions{
				Decsuite: ts.URL,
			},
		})
		require.NoError(t, err)

		// Store dynamic dataset token
		err = c.Tokens().SetDatasetToken(t.Context(), "test-ds", &tokens.DatasetToken{
			Dataset:     "test-ds",
			AccessToken: "dynamic-dataset-token",
		})
		require.NoError(t, err)

		api := &dszapistorelist.Api{
			Req: &dszapistorelist.Req{
				Params: &dszapistorelist.ReqParams{
					Dataset: lakedsz.MustParse("test-ds"),
				},
			},
		}
		status, err := c.Apis().DszListStores(t.Context(), api)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)
	})

	t.Run("Funcs executes analytical query via Dash runner", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "client-dash-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		dbPath := filepath.Join(tmpDir, "dash.sqlite")
		db, err := sql.Open("sqlite", dbPath)
		require.NoError(t, err)
		_, err = db.Exec("CREATE TABLE namedqueries (name TEXT PRIMARY KEY, query TEXT NOT NULL)")
		require.NoError(t, err)
		_, err = db.Exec(`
			INSERT INTO namedqueries (name, query) VALUES 
			('product-relevance-fashion-dash-funnel', 
			'SELECT json_object("results", json_object("rows", json_array(json_object("event_type", "view", "from_rec", "true", "qt_events", "999")))) AS ret')
		`)
		require.NoError(t, err)
		db.Close()

		dbBytes, err := os.ReadFile(dbPath)
		require.NoError(t, err)

		tsStorage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("ETag", `"etag-db"`)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(dbBytes)
		}))
		defer tsStorage.Close()

		tsApi := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/ds/dataset/fashion-ds/signed-url", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"urls":[{"url":%q}]}`, tsStorage.URL)))
		}))
		defer tsApi.Close()

		c, err := client.New(config.Options{
			Urls: config.UrlOptions{
				Decsuite: tsApi.URL,
			},
			Cache: config.CacheOptions{
				Dir: filepath.Join(tmpDir, "cache"),
			},
		})
		require.NoError(t, err)
		require.NotNil(t, c.Dash())
		require.NotNil(t, c.Funcs())

		fn := &relevancefashiondashfunnel.Func{
			Req: &relevancefashiondashfunnel.Req{
				Input: &relevancefashiondashfunnel.ReqInput{
					Dataset: lakedsz.MustParse("fashion-ds"),
				},
				Partitions: &relevancefashiondashfunnel.ReqPartitions{
					StoreId: "1",
				},
			},
		}

		status, err := c.Funcs().ProductRelevanceFashionDashFunnelV1(t.Context(), fn)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)
		require.NotNil(t, fn.Resp)
		require.NotNil(t, fn.Resp.Output)
		require.Len(t, fn.Resp.Output.Results.Rows, 1)
		require.Equal(t, "view", fn.Resp.Output.Results.Rows[0].EventType)
		require.Equal(t, "999", fn.Resp.Output.Results.Rows[0].QtEvents)
	})

	t.Run("Funcs executes inference workflow via Infer runner", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "client-infer-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		dbPath := filepath.Join(tmpDir, "model.sqlite")
		db, err := sql.Open("sqlite", dbPath)
		require.NoError(t, err)
		_, err = db.Exec("CREATE TABLE namedqueries (name TEXT PRIMARY KEY, query TEXT NOT NULL)")
		require.NoError(t, err)
		_, err = db.Exec(`
			INSERT INTO namedqueries (name, query) VALUES 
			('model-relevance-fashion-fbt-input', 'SELECT json_object("tensors", ?1) AS ret'),
			('model-relevance-fashion-fbt-output', 'SELECT json_object("results", json_array(json_object("sku_id", "rec-item-1", "rating", 0.99))) AS ret')
		`)
		require.NoError(t, err)
		db.Close()

		dbBytes, err := os.ReadFile(dbPath)
		require.NoError(t, err)

		tsStorage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("ETag", `"etag-model"`)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(dbBytes)
		}))
		defer tsStorage.Close()

		// Server handling decsuite (dszapimodelget, ingestz) and tfinfer (tfzapimodelpredict)
		tsDecsuite := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		}))
		defer tsDecsuite.Close()

		tsTfinfer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if strings.Contains(r.URL.Path, "/signed-get/") {
				_, _ = w.Write([]byte(fmt.Sprintf(`{"urls":[{"name":"sqlite.bin","url":%q}]}`, tsStorage.URL)))
			} else {
				_, _ = w.Write([]byte(`{"outputs":"{\"scores\":[0.99]}"}`))
			}
		}))
		defer tsTfinfer.Close()

		c, err := client.New(config.Options{
			Urls: config.UrlOptions{
				Decsuite: tsDecsuite.URL,
				Tfinfer:  tsTfinfer.URL,
			},
			Cache: config.CacheOptions{
				Dir: filepath.Join(tmpDir, "cache"),
			},
		})
		require.NoError(t, err)
		require.NotNil(t, c.Infer())
		require.NotNil(t, c.Funcs())

		fn := &inferzrelevancefashionfbt.Func{
			Req: &inferzrelevancefashionfbt.Req{
				Input: &inferzrelevancefashionfbt.ReqInput{
					Dataset:   lakedsz.MustParse("fashion-ds"),
					SessionId: "session-123",
					UserId:    "user-456",
					Skus: []*inferzrelevancefashionfbt.Sku{
						{SkuId: "cart-sku"},
					},
				},
				Partitions: &inferzrelevancefashionfbt.ReqPartitions{
					StoreId: "1",
				},
			},
		}

		status, err := c.Funcs().ProductRelevanceFashionFbtV1(t.Context(), fn)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)
		require.NotNil(t, fn.Resp)
		require.NotNil(t, fn.Resp.Output)
		require.Len(t, fn.Resp.Output.Results, 1)
		require.Equal(t, 0.99, fn.Resp.Output.Results[0].Rating)
		require.Len(t, fn.Resp.Output.Results[0].Skus, 1)
		require.Equal(t, "rec-item-1", fn.Resp.Output.Results[0].Skus[0].SkuId)
	})
}
