package dash_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapisignedurl"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashfunnel"
	"github.com/infinity6-ai/i6sdk/sdkgo/dash"
	"github.com/infinity6-ai/i6sdk/sdkgo/internal/fszcache"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func createTestSqliteDb(t *testing.T, dir string) string {
	t.Helper()
	dbPath := filepath.Join(dir, "bundle.bin")
	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE namedqueries (name TEXT PRIMARY KEY, query TEXT NOT NULL)")
	require.NoError(t, err)

	_, err = db.Exec(`
		INSERT INTO namedqueries (name, query) VALUES 
		(
			'product-relevance-fashion-dash-funnel', 
			'SELECT json_object(
				"results", json_object(
					"rows", json_array(
						json_object("event_type", "view", "from_rec", "true", "qt_events", "100"),
						json_object("event_type", "cart", "from_rec", "false", "qt_events", "25")
					)
				)
			) AS ret'
		),
		(
			'custom-dash',
			'SELECT json_object("summary", "total 50", "count", json_extract(?1, "$.limit")) AS ret'
		)
	`)
	require.NoError(t, err)
	return dbPath
}

func TestUnitPartitionsKey(t *testing.T) {
	t.Run("Empty partitions returns default", func(t *testing.T) {
		require.Equal(t, "default", dash.PartitionsKey(nil))
		require.Equal(t, "default", dash.PartitionsKey(map[string]string{}))
	})

	t.Run("Deterministic alphabetical sorting of partitions", func(t *testing.T) {
		input := map[string]string{
			"store_id": "99",
			"category": "fashion",
			"region":   "us-east",
		}
		expected := "category_fashion-region_us-east-store_id_99"
		require.Equal(t, expected, dash.PartitionsKey(input))
	})
}

func TestUnitDashQueryRawAndInto(t *testing.T) {
	ctx := context.Background()
	tmpDir, err := os.MkdirTemp("", "i6-dash-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dbPath := createTestSqliteDb(t, tmpDir)
	dbBytes, err := os.ReadFile(dbPath)
	require.NoError(t, err)

	// Local HTTP server serving the bundle.bin
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `"etag-db-1"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(dbBytes)
	}))
	defer ts.Close()

	cacheBaseUrl := urlz.MustParse("file://" + filepath.Join(tmpDir, "cache"))
	fsStore, err := fszcache.NewFs(cacheBaseUrl)
	require.NoError(t, err)
	downloader := fszcache.NewDownloader(fsStore)

	client := dash.New(dash.Options{
		Downloader: downloader,
	})

	t.Run("QueryRaw executes named query and returns JSON bytes", func(t *testing.T) {
		raw, err := client.QueryRaw(ctx, dash.QueryOptions{
			DashName:  "custom-dash",
			SqliteUrl: ts.URL,
			Params: map[string]any{
				"limit": 50,
			},
		})
		require.NoError(t, err)
		require.Contains(t, string(raw), `"summary":"total 50"`)
		require.Contains(t, string(raw), `"count":50`)
	})

	t.Run("QueryInto unmarshals result into typed struct", func(t *testing.T) {
		type CustomOutput struct {
			Summary string `json:"summary"`
			Count   int    `json:"count"`
		}

		var out CustomOutput
		err := client.QueryInto(ctx, &out, dash.QueryOptions{
			DashName:  "custom-dash",
			SqliteUrl: ts.URL,
			Params: map[string]any{
				"limit": 50,
			},
		})
		require.NoError(t, err)
		require.Equal(t, "total 50", out.Summary)
		require.Equal(t, 50, out.Count)
	})
}

func TestUnitDashSignedUrlResolution(t *testing.T) {
	ctx := context.Background()
	tmpDir, err := os.MkdirTemp("", "i6-dash-signedurl-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	var signedUrlCalled bool
	apis := &apicollectionz.SdkApis{
		Invoker: func(_ context.Context, api apiz.Api) (int, error) {
			if sApi, ok := api.(*dszapisignedurl.Api); ok {
				signedUrlCalled = true
				require.Equal(t, lakedsz.MustParse("my-dataset"), sApi.Req.Params.Dataset)
				require.Equal(t, "i6tb://results.relevance-fashion-store-dash/store_id=42/sqlite.bin", sApi.Req.Payload.Paths[0].Path.String())

				sApi.Resp = &dszapisignedurl.Resp{
					Payload: &dszapisignedurl.RespPayload{
						Urls: []*dszapisignedurl.Url{
							{Url: "https://storage.example.com/signed-sqlite.bin?token=xyz"},
						},
					},
				}
				return 200, nil
			}
			return 404, fmt.Errorf("unexpected api: %s", api.ApiSpec().Id)
		},
	}

	client := dash.New(dash.Options{
		Apis: apis,
	})

	url, err := client.ResolveSqliteUrl(ctx, dash.QueryOptions{
		Dataset: "my-dataset",
		Partitions: map[string]string{
			"store_id": "42",
		},
	})
	require.NoError(t, err)
	require.True(t, signedUrlCalled)
	require.Equal(t, "https://storage.example.com/signed-sqlite.bin?token=xyz", url)
}

func TestUnitDashExecuteSdkFunc(t *testing.T) {
	ctx := context.Background()
	tmpDir, err := os.MkdirTemp("", "i6-dash-exec-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dbPath := createTestSqliteDb(t, tmpDir)
	dbBytes, err := os.ReadFile(dbPath)
	require.NoError(t, err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `"etag-db-exec"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(dbBytes)
	}))
	defer ts.Close()

	apis := &apicollectionz.SdkApis{
		Invoker: func(_ context.Context, api apiz.Api) (int, error) {
			if sApi, ok := api.(*dszapisignedurl.Api); ok {
				sApi.Resp = &dszapisignedurl.Resp{
					Payload: &dszapisignedurl.RespPayload{
						Urls: []*dszapisignedurl.Url{
							{Url: ts.URL},
						},
					},
				}
				return 200, nil
			}
			return 404, nil
		},
	}

	fsStore, err := fszcache.NewFs(urlz.MustParse("file://" + filepath.Join(tmpDir, "cache")))
	require.NoError(t, err)
	downloader := fszcache.NewDownloader(fsStore)

	client := dash.New(dash.Options{
		Apis:       apis,
		Downloader: downloader,
	})

	fn := &relevancefashiondashfunnel.Func{
		Req: &relevancefashiondashfunnel.Req{
			Input: &relevancefashiondashfunnel.ReqInput{
				Dataset: lakedsz.MustParse("fashion-dataset"),
				Params: &relevancefashiondashfunnel.Params{
					MinYearMonth: "2026-07",
					MaxYearMonth: "2026-11",
				},
			},
			Partitions: &relevancefashiondashfunnel.ReqPartitions{
				StoreId: "1",
			},
		},
	}

	status, err := client.Execute(ctx, fn)
	require.NoError(t, err)
	require.Equal(t, 200, status)

	// Verify fn.Resp.Output was successfully populated from SQLite execution
	require.NotNil(t, fn.Resp)
	require.NotNil(t, fn.Resp.Output)
	require.Len(t, fn.Resp.Output.Results.Rows, 2)
	require.Equal(t, "view", fn.Resp.Output.Results.Rows[0].EventType)
	require.Equal(t, "true", fn.Resp.Output.Results.Rows[0].FromRec)
	require.Equal(t, "100", fn.Resp.Output.Results.Rows[0].QtEvents)
	require.Equal(t, "cart", fn.Resp.Output.Results.Rows[1].EventType)
	require.Equal(t, "25", fn.Resp.Output.Results.Rows[1].QtEvents)
}
