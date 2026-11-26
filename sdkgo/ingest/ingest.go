package ingest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/infinity6-ai/gox/commonz/validation/checked"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapigeturl"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbnamez"
)

// UploadProgress represents the progress status of an in-flight upload.
type UploadProgress struct {
	Loaded           int64   `json:"loaded"`
	Total            int64   `json:"total"`
	Percent          float64 `json:"percent"`
	CompletedStreams int     `json:"completed_streams"`
	TotalStreams     int     `json:"total_streams"`
	StreamIndex      int     `json:"stream_index"`
}

// UploadResult represents the result of a completed upload operation.
type UploadResult struct {
	IngestId string `json:"ingest_id"`
	Count    int    `json:"count"`
}

// SignedUrlResult represents the response of an upload URL request.
type SignedUrlResult struct {
	IngestId string   `json:"ingest_id"`
	Url      string   `json:"url"`
	Urls     []string `json:"urls"`
}

// StreamItem represents a readable stream with an optional known byte length.
type StreamItem struct {
	Reader io.Reader
	Size   int64
}

// FromReader wraps an io.Reader into a StreamItem with optional size.
func FromReader(r io.Reader, size ...int64) StreamItem {
	s := int64(0)
	if len(size) > 0 && size[0] > 0 {
		s = size[0]
	}
	return StreamItem{
		Reader: r,
		Size:   s,
	}
}

// FromBytes creates a StreamItem from a raw byte slice.
func FromBytes(b []byte) StreamItem {
	return StreamItem{
		Reader: bytes.NewReader(b),
		Size:   int64(len(b)),
	}
}

// FromString creates a StreamItem from a string.
func FromString(s string) StreamItem {
	return FromBytes([]byte(s))
}

// FromFile creates a StreamItem from a local file path.
func FromFile(path string) (StreamItem, error) {
	file, err := os.Open(path)
	if err != nil {
		return StreamItem{}, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return StreamItem{}, fmt.Errorf("failed to stat file %s: %w", path, err)
	}
	return StreamItem{
		Reader: file,
		Size:   info.Size(),
	}, nil
}

// UploadOptions configures a multi-stream ingestion upload.
type UploadOptions struct {
	Dataset    string
	Table      string
	Partitions map[string]string
	Streams    []StreamItem
	OnProgress func(progress UploadProgress)
	Timeout    time.Duration
	Workers    int
}

// Options configures the Ingest client.
type Options struct {
	Apis       *apicollectionz.SdkApis
	HttpClient *http.Client
}

// Client coordinates presigned URL acquisition and multi-stream parallel uploads.
type Client struct {
	apis       *apicollectionz.SdkApis
	httpClient *http.Client
}

// New creates a new Ingest Client.
func New(opts Options) *Client {
	httpClient := opts.HttpClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 120 * time.Second,
		}
	}
	return &Client{
		apis:       opts.Apis,
		httpClient: httpClient,
	}
}

// GetUrl requests signed upload PUT URLs from the ingestion service.
func (c *Client) GetUrl(
	ctx context.Context,
	dataset string,
	table string,
	partitions map[string]string,
	amount int,
) (*SignedUrlResult, error) {
	if dataset == "" {
		return nil, fmt.Errorf("ingest: dataset is required")
	}
	if table == "" {
		return nil, fmt.Errorf("ingest: table is required")
	}
	if amount <= 0 {
		amount = 1
	}

	if c.apis == nil {
		return nil, fmt.Errorf("ingest: apis client not wired")
	}

	var ds lakedsz.Dataset
	if err := checked.Set(&ds, dataset); err != nil {
		return nil, fmt.Errorf("ingest: invalid dataset %q: %w", dataset, err)
	}

	var tb laketbnamez.TableName
	if err := checked.Set(&tb, table); err != nil {
		return nil, fmt.Errorf("ingest: invalid table %q: %w", table, err)
	}

	api := &ingestzapigeturl.Api{
		Req: &ingestzapigeturl.Req{
			Params: &ingestzapigeturl.ReqParams{
				Dataset: ds,
			},
			Payload: &ingestzapigeturl.ReqPayload{
				Table:      tb,
				Partitions: partitions,
				Amount:     amount,
			},
		},
	}

	status, err := c.apis.IngestzGetUrl(ctx, api)
	if err != nil || status < 200 || status >= 300 || api.Resp == nil || api.Resp.Payload == nil {
		return nil, fmt.Errorf("ingest: failed to get upload URLs (status %d): %w", status, err)
	}

	p := api.Resp.Payload
	firstUrl := ""
	if len(p.Uploads) > 0 {
		firstUrl = p.Uploads[0]
	}

	return &SignedUrlResult{
		IngestId: string(p.IngestId),
		Url:      firstUrl,
		Urls:     p.Uploads,
	}, nil
}

// Upload uploads one or more streams in parallel to the ingest service using signed PUT URLs.
func (c *Client) Upload(ctx context.Context, opts UploadOptions) (*UploadResult, error) {
	if opts.Dataset == "" {
		return nil, fmt.Errorf("ingest: dataset is required")
	}
	if opts.Table == "" {
		return nil, fmt.Errorf("ingest: table is required")
	}
	if len(opts.Streams) == 0 {
		return nil, fmt.Errorf("ingest: streams cannot be empty")
	}

	count := len(opts.Streams)
	var totalBytesAll int64
	for _, s := range opts.Streams {
		if s.Size > 0 {
			totalBytesAll += s.Size
		}
	}

	urlInfo, err := c.GetUrl(ctx, opts.Dataset, opts.Table, opts.Partitions, count)
	if err != nil {
		return nil, err
	}

	if len(urlInfo.Urls) < count {
		return nil, fmt.Errorf("ingest: server returned %d upload URLs for %d streams", len(urlInfo.Urls), count)
	}

	workers := opts.Workers
	if workers <= 0 {
		workers = count
		if workers > 8 {
			workers = 8
		}
	}

	type uploadJob struct {
		index     int
		stream    StreamItem
		targetUrl string
	}

	jobs := make(chan uploadJob, count)
	for i := 0; i < count; i++ {
		jobs <- uploadJob{
			index:     i,
			stream:    opts.Streams[i],
			targetUrl: urlInfo.Urls[i],
		}
	}
	close(jobs)

	var mu sync.Mutex
	bytesPerStream := make([]int64, count)
	completedStreamsCount := 0

	updateProgress := func(streamIndex int, delta int64) {
		if opts.OnProgress == nil {
			return
		}
		mu.Lock()
		defer mu.Unlock()

		bytesPerStream[streamIndex] += delta
		var totalLoaded int64
		for _, b := range bytesPerStream {
			totalLoaded += b
		}

		pct := 100.0
		if totalBytesAll > 0 {
			pct = (float64(totalLoaded) / float64(totalBytesAll)) * 100.0
		}
		if pct > 100.0 {
			pct = 100.0
		}

		totalReport := totalBytesAll
		if totalReport <= 0 {
			totalReport = totalLoaded
		}

		opts.OnProgress(UploadProgress{
			Loaded:           totalLoaded,
			Total:            totalReport,
			Percent:          pct,
			CompletedStreams: completedStreamsCount,
			TotalStreams:     count,
			StreamIndex:      streamIndex,
		})
	}

	uploadCtx, cancel := context.WithCancel(ctx)
	if opts.Timeout > 0 {
		uploadCtx, cancel = context.WithTimeout(ctx, opts.Timeout)
	}
	defer cancel()

	var wg sync.WaitGroup
	var errOnce sync.Once
	var firstErr error

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				if uploadCtx.Err() != nil {
					return
				}

				pr := &progressReader{
					r: job.stream.Reader,
					onRead: func(n int) {
						updateProgress(job.index, int64(n))
					},
					ctx: uploadCtx,
				}

				req, err := http.NewRequestWithContext(uploadCtx, http.MethodPut, job.targetUrl, pr)
				if err != nil {
					errOnce.Do(func() {
						firstErr = fmt.Errorf("failed to create upload request for stream %d: %w", job.index, err)
						cancel()
					})
					return
				}

				if job.stream.Size > 0 {
					req.ContentLength = job.stream.Size
				}

				resp, err := c.httpClient.Do(req)
				if err != nil {
					errOnce.Do(func() {
						firstErr = fmt.Errorf("upload stream %d failed: %w", job.index, err)
						cancel()
					})
					return
				}
				_ = resp.Body.Close()

				if resp.StatusCode < 200 || resp.StatusCode >= 300 {
					errOnce.Do(func() {
						firstErr = fmt.Errorf("upload stream %d returned HTTP %d", job.index, resp.StatusCode)
						cancel()
					})
					return
				}

				mu.Lock()
				completedStreamsCount++
				mu.Unlock()
				updateProgress(job.index, 0)
			}
		}()
	}

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}
	if err := uploadCtx.Err(); err != nil {
		return nil, fmt.Errorf("ingest upload aborted: %w", err)
	}

	return &UploadResult{
		IngestId: urlInfo.IngestId,
		Count:    count,
	}, nil
}

type progressReader struct {
	r      io.Reader
	onRead func(n int)
	ctx    context.Context
}

func (p *progressReader) Read(buf []byte) (int, error) {
	if err := p.ctx.Err(); err != nil {
		return 0, err
	}
	n, err := p.r.Read(buf)
	if n > 0 && p.onRead != nil {
		p.onRead(n)
	}
	return n, err
}
