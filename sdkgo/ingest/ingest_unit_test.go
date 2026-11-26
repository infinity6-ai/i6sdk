package ingest_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapigeturl"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/ingestztoken"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbnamez"
	"github.com/infinity6-ai/i6sdk/sdkgo/ingest"
	"github.com/stretchr/testify/require"
)

func TestUnitStreamItemFactories(t *testing.T) {
	t.Run("FromBytes creates stream with correct size", func(t *testing.T) {
		data := []byte("hello stream")
		item := ingest.FromBytes(data)
		require.Equal(t, int64(len(data)), item.Size)

		content, err := io.ReadAll(item.Reader)
		require.NoError(t, err)
		require.Equal(t, data, content)
	})

	t.Run("FromString creates stream with string content", func(t *testing.T) {
		str := "string payload"
		item := ingest.FromString(str)
		require.Equal(t, int64(len(str)), item.Size)

		content, err := io.ReadAll(item.Reader)
		require.NoError(t, err)
		require.Equal(t, []byte(str), content)
	})

	t.Run("FromReader wraps reader", func(t *testing.T) {
		buf := bytes.NewBufferString("buffered data")
		item := ingest.FromReader(buf, int64(buf.Len()))
		require.Equal(t, int64(13), item.Size)

		content, err := io.ReadAll(item.Reader)
		require.NoError(t, err)
		require.Equal(t, "buffered data", string(content))
	})

	t.Run("FromFile loads file with size", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "ingest-file-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		filePath := filepath.Join(tmpDir, "test.csv")
		fileContent := "id,name\n1,alpha\n2,beta\n"
		err = os.WriteFile(filePath, []byte(fileContent), 0644)
		require.NoError(t, err)

		item, err := ingest.FromFile(filePath)
		require.NoError(t, err)
		require.Equal(t, int64(len(fileContent)), item.Size)

		readBytes, err := io.ReadAll(item.Reader)
		require.NoError(t, err)
		require.Equal(t, fileContent, string(readBytes))
	})
}

func TestUnitIngestGetUrl(t *testing.T) {
	ctx := context.Background()

	t.Run("Requests signed URLs from ingest service", func(t *testing.T) {
		var called bool
		apis := &apicollectionz.SdkApis{
			Invoker: func(_ context.Context, api apiz.Api) (int, error) {
				if gApi, ok := api.(*ingestzapigeturl.Api); ok {
					called = true
					require.Equal(t, lakedsz.MustParse("my-dataset"), gApi.Req.Params.Dataset)
					require.Equal(t, laketbnamez.MustParse("events"), gApi.Req.Payload.Table)
					require.Equal(t, 2, gApi.Req.Payload.Amount)
					require.Equal(t, "2026-11", gApi.Req.Payload.Partitions["month"])

					gApi.Resp = &ingestzapigeturl.Resp{
						Payload: &ingestzapigeturl.RespPayload{
							IngestId: ingestztoken.MustParse("ingest-token-123"),
							Uploads: []string{
								"https://storage.example.com/put-1",
								"https://storage.example.com/put-2",
							},
						},
					}
					return 200, nil
				}
				return 404, nil
			},
		}

		client := ingest.New(ingest.Options{Apis: apis})
		res, err := client.GetUrl(ctx, "my-dataset", "events", map[string]string{"month": "2026-11"}, 2)
		require.NoError(t, err)
		require.True(t, called)
		require.Equal(t, "ingest-token-123", res.IngestId)
		require.Equal(t, "https://storage.example.com/put-1", res.Url)
		require.Len(t, res.Urls, 2)
	})

	t.Run("Validation errors on missing fields", func(t *testing.T) {
		client := ingest.New(ingest.Options{})
		_, err := client.GetUrl(ctx, "", "events", nil, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "dataset is required")

		_, err = client.GetUrl(ctx, "my-dataset", "", nil, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "table is required")
	})
}

func TestUnitIngestUpload(t *testing.T) {
	ctx := context.Background()

	t.Run("Uploads multiple streams in parallel with progress tracking", func(t *testing.T) {
		var mu sync.Mutex
		receivedData := make(map[string][]byte)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			mu.Lock()
			receivedData[r.URL.Path] = body
			mu.Unlock()

			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		apis := &apicollectionz.SdkApis{
			Invoker: func(_ context.Context, api apiz.Api) (int, error) {
				if gApi, ok := api.(*ingestzapigeturl.Api); ok {
					count := gApi.Req.Payload.Amount
					urls := make([]string, count)
					for i := 0; i < count; i++ {
						urls[i] = fmt.Sprintf("%s/upload-%d", ts.URL, i)
					}
					gApi.Resp = &ingestzapigeturl.Resp{
						Payload: &ingestzapigeturl.RespPayload{
							IngestId: ingestztoken.MustParse("ingest-456"),
							Uploads:  urls,
						},
					}
					return 200, nil
				}
				return 404, nil
			},
		}

		client := ingest.New(ingest.Options{
			Apis:       apis,
			HttpClient: ts.Client(),
		})

		stream1 := ingest.FromString("stream 1 content")
		stream2 := ingest.FromString("stream 2 content larger")
		stream3 := ingest.FromString("stream 3 final")

		totalExpectedBytes := stream1.Size + stream2.Size + stream3.Size

		var progressCalls int32
		var lastProgress ingest.UploadProgress

		res, err := client.Upload(ctx, ingest.UploadOptions{
			Dataset: "fashion-dataset",
			Table:   "catalog",
			Partitions: map[string]string{
				"category": "shoes",
			},
			Streams: []ingest.StreamItem{stream1, stream2, stream3},
			OnProgress: func(p ingest.UploadProgress) {
				atomic.AddInt32(&progressCalls, 1)
				mu.Lock()
				lastProgress = p
				mu.Unlock()
			},
			Workers: 3,
		})

		require.NoError(t, err)
		require.Equal(t, "ingest-456", res.IngestId)
		require.Equal(t, 3, res.Count)

		// Verify data arrived at mock server
		mu.Lock()
		require.Len(t, receivedData, 3)
		require.Equal(t, []byte("stream 1 content"), receivedData["/upload-0"])
		require.Equal(t, []byte("stream 2 content larger"), receivedData["/upload-1"])
		require.Equal(t, []byte("stream 3 final"), receivedData["/upload-2"])

		require.Greater(t, atomic.LoadInt32(&progressCalls), int32(0))
		require.Equal(t, totalExpectedBytes, lastProgress.Total)
		require.Equal(t, totalExpectedBytes, lastProgress.Loaded)
		require.Equal(t, 100.0, lastProgress.Percent)
		require.Equal(t, 3, lastProgress.CompletedStreams)
		mu.Unlock()
	})

	t.Run("Upload supports context cancellation", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond) // slow server
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		apis := &apicollectionz.SdkApis{
			Invoker: func(_ context.Context, api apiz.Api) (int, error) {
				if gApi, ok := api.(*ingestzapigeturl.Api); ok {
					gApi.Resp = &ingestzapigeturl.Resp{
						Payload: &ingestzapigeturl.RespPayload{
							IngestId: ingestztoken.MustParse("ingest-cancel"),
							Uploads:  []string{ts.URL + "/cancel-stream"},
						},
					}
					return 200, nil
				}
				return 404, nil
			},
		}

		client := ingest.New(ingest.Options{
			Apis:       apis,
			HttpClient: ts.Client(),
		})

		cancelCtx, cancel := context.WithCancel(ctx)
		// Cancel immediately after starting
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		largePayload := make([]byte, 1024*1024)
		_, err := client.Upload(cancelCtx, ingest.UploadOptions{
			Dataset: "fashion-dataset",
			Table:   "catalog",
			Streams: []ingest.StreamItem{ingest.FromBytes(largePayload)},
		})

		require.Error(t, err)
	})

	t.Run("Server error on PUT returns actionable error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		apis := &apicollectionz.SdkApis{
			Invoker: func(_ context.Context, api apiz.Api) (int, error) {
				if gApi, ok := api.(*ingestzapigeturl.Api); ok {
					gApi.Resp = &ingestzapigeturl.Resp{
						Payload: &ingestzapigeturl.RespPayload{
							IngestId: ingestztoken.MustParse("ingest-err"),
							Uploads:  []string{ts.URL + "/upload-err"},
						},
					}
					return 200, nil
				}
				return 404, nil
			},
		}

		client := ingest.New(ingest.Options{
			Apis:       apis,
			HttpClient: ts.Client(),
		})

		_, err := client.Upload(ctx, ingest.UploadOptions{
			Dataset: "fashion-dataset",
			Table:   "catalog",
			Streams: []ingest.StreamItem{ingest.FromString("content")},
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "HTTP 500")
	})
}
