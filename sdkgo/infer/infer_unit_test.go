package infer_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapimodelget"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapieventevent"
	"github.com/infinity6-ai/i6sdk/apitfz/tfzapimodelpredict"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/inferz/inferzrelevancefashionfbt"
	"github.com/infinity6-ai/i6sdk/sdkgo/infer"
	"github.com/infinity6-ai/i6sdk/sdkgo/internal/fszcache"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func createTestModelSqliteDb(t *testing.T, dir string) string {
	t.Helper()
	dbPath := filepath.Join(dir, "model.bin")
	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE namedqueries (name TEXT PRIMARY KEY, query TEXT NOT NULL)")
	require.NoError(t, err)

	_, err = db.Exec(`
		INSERT INTO namedqueries (name, query) VALUES 
		(
			'model-relevance-fashion-fbt-input',
			'SELECT json_object("tensors", ?1) AS ret'
		),
		(
			'model-relevance-fashion-fbt-output',
			'SELECT json_object(
				"results", json_array(
					json_object("sku_id", "sku-rec-1", "rating", 0.95),
					json_object("sku_id", "sku-in-cart", "rating", 0.90),
					json_object("sku_id", "sku-rec-2", "rating", 0.85)
				)
			) AS ret WHERE ?1 IS NOT NULL AND ?2 > 0'
		)
	`)
	require.NoError(t, err)
	return dbPath
}

func TestUnitDefaultSelectResults(t *testing.T) {
	t.Run("Filters cart SKUs and deduplicates", func(t *testing.T) {
		req := map[string]any{
			"input": map[string]any{
				"skus": []any{
					map[string]any{"sku_id": "cart_sku_1"},
				},
			},
		}

		results := map[string]any{
			"results": []any{
				map[string]any{"sku_id": "cart_sku_1", "rating": 0.99}, // in cart: should be filtered
				map[string]any{"sku_id": "rec_sku_1", "rating": 0.88},
				map[string]any{"sku_id": "rec_sku_1", "rating": 0.77},  // duplicate: should be filtered
				map[string]any{"sku_id": "rec_sku_2", "rating": 0.66},
			},
		}

		selected := infer.DefaultSelectResults(req, results)
		resList, ok := selected["results"].([]any)
		require.True(t, ok)
		require.Len(t, resList, 2)

		item1 := resList[0].(map[string]any)
		require.Equal(t, "rec_sku_1", item1["sku_id"])
		require.Equal(t, 0.88, item1["rating"])

		item2 := resList[1].(map[string]any)
		require.Equal(t, "rec_sku_2", item2["sku_id"])
		require.Equal(t, 0.66, item2["rating"])
	})
}

func TestUnitInferResolveSqliteUrl(t *testing.T) {
	ctx := context.Background()

	t.Run("Resolves sqlite.bin URL from model metadata", func(t *testing.T) {
		apis := &apicollectionz.SdkApis{
			Invoker: func(_ context.Context, api apiz.Api) (int, error) {
				if mApi, ok := api.(*dszapimodelget.Api); ok {
					require.Equal(t, lakedsz.MustParse("my-ds"), mApi.Req.Params.Dataset)
					require.Equal(t, "relevance-fashion-fbt", mApi.Req.Params.ModelName)
					mApi.Resp = &dszapimodelget.Resp{
						Payload: &dszapimodelget.RespPayload{
							Urls: []*dszapimodelget.ModelFile{
								{Name: "config.json", Url: "https://example.com/config.json"},
								{Name: "sqlite.bin", Url: "https://example.com/sqlite.bin"},
							},
						},
					}
					return 200, nil
				}
				return 404, nil
			},
		}

		client := infer.New(infer.Options{Apis: apis})
		url, err := client.ResolveSqliteUrl(ctx, infer.QueryOptions{
			Dataset:   "my-ds",
			ModelName: "relevance-fashion-fbt",
		})
		require.NoError(t, err)
		require.Equal(t, "https://example.com/sqlite.bin", url)
	})

	t.Run("Returns error when sqlite.bin file is missing", func(t *testing.T) {
		apis := &apicollectionz.SdkApis{
			Invoker: func(_ context.Context, api apiz.Api) (int, error) {
				if mApi, ok := api.(*dszapimodelget.Api); ok {
					mApi.Resp = &dszapimodelget.Resp{
						Payload: &dszapimodelget.RespPayload{
							Urls: []*dszapimodelget.ModelFile{
								{Name: "config.json", Url: "https://example.com/config.json"},
							},
						},
					}
					return 200, nil
				}
				return 404, nil
			},
		}

		client := infer.New(infer.Options{Apis: apis})
		_, err := client.ResolveSqliteUrl(ctx, infer.QueryOptions{
			Dataset:   "my-ds",
			ModelName: "relevance-fashion-fbt",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "has no sqlite.bin file available")
	})
}

func TestUnitInfer3StepPipeline(t *testing.T) {
	ctx := context.Background()
	tmpDir, err := os.MkdirTemp("", "i6-infer-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dbPath := createTestModelSqliteDb(t, tmpDir)
	dbBytes, err := os.ReadFile(dbPath)
	require.NoError(t, err)

	tsStorage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `"etag-infer-db"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(dbBytes)
	}))
	defer tsStorage.Close()

	var predictCalls int32
	var recrawCalls int32
	var recshowCalls int32

	apis := &apicollectionz.SdkApis{
		Invoker: func(_ context.Context, api apiz.Api) (int, error) {
			switch typed := api.(type) {
			case *dszapimodelget.Api:
				typed.Resp = &dszapimodelget.Resp{
					Payload: &dszapimodelget.RespPayload{
						Urls: []*dszapimodelget.ModelFile{
							{Name: "sqlite.bin", Url: tsStorage.URL},
						},
					},
				}
				return 200, nil
			case *tfzapimodelpredict.Api:
				atomic.AddInt32(&predictCalls, 1)
				require.Equal(t, "relevance-fashion-fbt", typed.Req.Params.ModelName)
				require.NotEmpty(t, typed.Req.Payload.Inputs)
				typed.Resp = &tfzapimodelpredict.Resp{
					Payload: &tfzapimodelpredict.RespPayload{
						Outputs: `{"predictions": [0.95, 0.90, 0.85]}`,
					},
				}
				return 200, nil
			case *ingestzapieventevent.Api:
				if typed.Req.Payload.Recraw != nil {
					atomic.AddInt32(&recrawCalls, 1)
					require.Equal(t, "relevance-fashion-fbt", typed.Req.Payload.Recraw.ModelId)
				}
				if typed.Req.Payload.Recshow != nil {
					atomic.AddInt32(&recshowCalls, 1)
					require.Equal(t, "relevance-fashion-fbt", typed.Req.Payload.Recshow.ModelId)
				}
				return 200, nil
			default:
				return 404, fmt.Errorf("unexpected api: %s", api.ApiSpec().Id)
			}
		},
	}

	fsStore, err := fszcache.NewFs(urlz.MustParse("file://" + filepath.Join(tmpDir, "cache")))
	require.NoError(t, err)
	downloader := fszcache.NewDownloader(fsStore)

	client := infer.New(infer.Options{
		Apis:       apis,
		Downloader: downloader,
	})

	t.Run("Run executes 3-step pipeline and dispatches recraw", func(t *testing.T) {
		opts := infer.QueryOptions{
			Dataset:   "fashion-ds",
			ModelName: "relevance-fashion-fbt",
			SessionId: "sess-123",
			UserId:    "user-456",
			Params: map[string]any{
				"skus": []any{"sku-10"},
			},
			Partitions: map[string]string{
				"store_id": "1",
				"version":  "v1",
			},
		}

		res, err := client.Run(ctx, opts)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotEmpty(t, res["rec_id"])
		require.Equal(t, int32(1), atomic.LoadInt32(&predictCalls))
		require.Equal(t, int32(1), atomic.LoadInt32(&recrawCalls))

		// Step 3 output had 3 items
		resultsList, ok := res["results"].([]any)
		require.True(t, ok)
		require.Len(t, resultsList, 3)

		// Apply SelectResults with cart filtering
		selected, err := client.SelectResults(ctx, infer.SelectOptions{
			Dataset:      opts.Dataset,
			ModelName:    opts.ModelName,
			ModelVersion: opts.Partitions["version"],
			StoreId:      opts.Partitions["store_id"],
			SessionId:    opts.SessionId,
			UserId:       opts.UserId,
			RecId:        res["rec_id"].(string),
			Req: map[string]any{
				"input": map[string]any{
					"skus": []any{"sku-in-cart"}, // in cart!
				},
			},
			Results: res,
		})
		require.NoError(t, err)
		require.Equal(t, int32(1), atomic.LoadInt32(&recshowCalls))

		shownList, ok := selected["results"].([]any)
		require.True(t, ok)
		// sku-in-cart was filtered out: exactly 2 recommendations left
		require.Len(t, shownList, 2)
		require.Equal(t, "sku-rec-1", shownList[0].(map[string]any)["sku_id"])
		require.Equal(t, "sku-rec-2", shownList[1].(map[string]any)["sku_id"])
	})
}

func TestUnitInferExecuteSdkFunc(t *testing.T) {
	ctx := context.Background()
	tmpDir, err := os.MkdirTemp("", "i6-infer-func-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dbPath := createTestModelSqliteDb(t, tmpDir)
	dbBytes, err := os.ReadFile(dbPath)
	require.NoError(t, err)

	tsStorage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `"etag-db-func"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(dbBytes)
	}))
	defer tsStorage.Close()

	apis := &apicollectionz.SdkApis{
		Invoker: func(_ context.Context, api apiz.Api) (int, error) {
			switch typed := api.(type) {
			case *dszapimodelget.Api:
				typed.Resp = &dszapimodelget.Resp{
					Payload: &dszapimodelget.RespPayload{
						Urls: []*dszapimodelget.ModelFile{
							{Name: "sqlite.bin", Url: tsStorage.URL},
						},
					},
				}
				return 200, nil
			case *tfzapimodelpredict.Api:
				typed.Resp = &tfzapimodelpredict.Resp{
					Payload: &tfzapimodelpredict.RespPayload{
						Outputs: `{"predictions": [0.95, 0.90, 0.85]}`,
					},
				}
				return 200, nil
			case *ingestzapieventevent.Api:
				return 200, nil
			default:
				return 404, nil
			}
		},
	}

	fsStore, err := fszcache.NewFs(urlz.MustParse("file://" + filepath.Join(tmpDir, "cache")))
	require.NoError(t, err)
	downloader := fszcache.NewDownloader(fsStore)

	client := infer.New(infer.Options{
		Apis:       apis,
		Downloader: downloader,
	})

	fn := &inferzrelevancefashionfbt.Func{
		Req: &inferzrelevancefashionfbt.Req{
			Input: &inferzrelevancefashionfbt.ReqInput{
				Dataset:   lakedsz.MustParse("fashion-dataset"),
				SessionId: "sess-999",
				UserId:    "user-888",
				Skus: []*inferzrelevancefashionfbt.Sku{
					{SkuId: "sku-in-cart"}, // in cart
				},
			},
			Partitions: &inferzrelevancefashionfbt.ReqPartitions{
				StoreId: "1",
			},
		},
	}

	status, err := client.Execute(ctx, fn)
	require.NoError(t, err)
	require.Equal(t, 200, status)

	require.NotNil(t, fn.Resp)
	require.NotNil(t, fn.Resp.Output)
	// sku-in-cart was filtered out: 2 items returned
	require.Len(t, fn.Resp.Output.Results, 2)
	require.Equal(t, 0.95, fn.Resp.Output.Results[0].Rating)
	require.Len(t, fn.Resp.Output.Results[0].Skus, 1)
	require.Equal(t, "sku-rec-1", fn.Resp.Output.Results[0].Skus[0].SkuId)
}
