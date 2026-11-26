package dash

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/infinity6-ai/gox/fsz/fsz/fszsignmethod"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapisignedurl"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbpathz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
	"github.com/infinity6-ai/i6sdk/sdkgo/internal/fszcache"
	"github.com/infinity6-ai/i6sdk/sdkgo/internal/sqlitez"
)

// QueryOptions holds execution options for analytical SQLite dash queries.
type QueryOptions struct {
	Dataset    string            `json:"dataset"`
	DashName   string            `json:"dash_name"`
	Partitions map[string]string `json:"partitions,omitempty"`
	Params     any               `json:"params,omitempty"`
	SqliteUrl  string            `json:"sqlite_url,omitempty"`
}

// PartitionsKey builds a deterministic cache partition key from a partitions map.
func PartitionsKey(partitions map[string]string) string {
	if len(partitions) == 0 {
		return "default"
	}
	keys := make([]string, 0, len(partitions))
	for k := range partitions {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = fmt.Sprintf("%s_%s", k, partitions[k])
	}
	return strings.Join(parts, "-")
}

// Options configures the Dash client runner.
type Options struct {
	Apis       *apicollectionz.SdkApis
	Downloader *fszcache.Downloader
}

// Client executes dashboard analytical queries on pre-aggregated SQLite bundles.
type Client struct {
	apis       *apicollectionz.SdkApis
	downloader *fszcache.Downloader
}

// New constructs a new Dash Client.
func New(opts Options) *Client {
	return &Client{
		apis:       opts.Apis,
		downloader: opts.Downloader,
	}
}

// Downloader returns the underlying Downloader instance.
func (c *Client) Downloader() *fszcache.Downloader {
	return c.downloader
}

// ResolveSqliteUrl resolves the presigned download URL for the dash SQLite database bundle.
func (c *Client) ResolveSqliteUrl(ctx context.Context, opts QueryOptions) (string, error) {
	if opts.SqliteUrl != "" {
		return opts.SqliteUrl, nil
	}

	if opts.Dataset == "" {
		return "", fmt.Errorf("dash: dataset is required to resolve sqlite url")
	}

	storeId := "1"
	if s, ok := opts.Partitions["store_id"]; ok && s != "" {
		storeId = s
	} else if s, ok := opts.Partitions["storeId"]; ok && s != "" {
		storeId = s
	}

	pathStr := fmt.Sprintf("i6tb://results.relevance-fashion-store-dash/store_id=%s/sqlite.bin", storeId)

	if c.apis == nil {
		return "", fmt.Errorf("dash: cannot resolve signed url: apis client not wired")
	}

	signApi := &dszapisignedurl.Api{
		Req: &dszapisignedurl.Req{
			Params: &dszapisignedurl.ReqParams{
				Dataset: lakedsz.MustParse(opts.Dataset),
			},
			Payload: &dszapisignedurl.ReqPayload{
				Paths: []*dszapisignedurl.Path{
					{
						Method: fszsignmethod.SignMethodGet,
						Path:   laketbpathz.MustParse(pathStr),
					},
				},
			},
		},
	}

	status, err := c.apis.DszSignedUrl(ctx, signApi)
	if err != nil || status < 200 || status >= 300 || signApi.Resp == nil || signApi.Resp.Payload == nil || len(signApi.Resp.Payload.Urls) == 0 {
		return "", fmt.Errorf("dash: failed to resolve signed url for %s (status %d): %w", pathStr, status, err)
	}

	firstUrl := signApi.Resp.Payload.Urls[0].Url
	if strings.TrimSpace(firstUrl) == "" {
		return "", fmt.Errorf("dash: received empty signed url for %s", pathStr)
	}

	return firstUrl, nil
}

// QueryRaw downloads the SQLite database bundle, executes the named query, and returns the raw JSON bytes of the 'ret' column.
func (c *Client) QueryRaw(ctx context.Context, opts QueryOptions) ([]byte, error) {
	if opts.DashName == "" {
		return nil, fmt.Errorf("dash: dash_name is required")
	}

	if c.downloader == nil {
		return nil, fmt.Errorf("dash: downloader not configured")
	}

	pkey := PartitionsKey(opts.Partitions)
	cacheName := fmt.Sprintf("%s-%s-sqlite", opts.DashName, pkey)

	sqliteUrl, err := c.ResolveSqliteUrl(ctx, opts)
	if err != nil {
		return nil, err
	}

	sqlitePath, err := c.downloader.Update(ctx, cacheName, sqliteUrl, 300.0)
	if err != nil {
		return nil, fmt.Errorf("dash: failed to download/update sqlite bundle: %w", err)
	}

	conn, err := sqlitez.Open(sqlitePath)
	if err != nil {
		return nil, fmt.Errorf("dash: failed to open sqlite connection at %s: %w", sqlitePath, err)
	}
	defer conn.Close()

	var queryParams []any
	if opts.Params != nil {
		queryParams = append(queryParams, opts.Params)
	}

	return conn.NamedJsonQuery(ctx, opts.DashName, queryParams...)
}

// QueryInto executes QueryRaw and unmarshals the JSON result directly into target.
func (c *Client) QueryInto(ctx context.Context, target any, opts QueryOptions) error {
	raw, err := c.QueryRaw(ctx, opts)
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return fmt.Errorf("dash: failed to unmarshal query result for %s: %w", opts.DashName, err)
	}
	return nil
}

// Execute runs an analytical sdkfuncz.Func workflow function through the Dash SQLite execution pipeline.
func (c *Client) Execute(ctx context.Context, fn sdkfuncz.Func) (int, error) {
	spec := fn.SdkFuncSpec()
	if spec.Id == "" {
		return 400, fmt.Errorf("dash: sdkfunc spec id is empty")
	}

	rawFn, err := json.Marshal(fn)
	if err != nil {
		return 400, fmt.Errorf("dash: failed to marshal sdkfunc request: %w", err)
	}

	var envelope struct {
		Req struct {
			Input struct {
				Dataset string `json:"dataset"`
				Params  any    `json:"params"`
			} `json:"input"`
			Partitions map[string]string `json:"partitions"`
		} `json:"req"`
	}
	if err := json.Unmarshal(rawFn, &envelope); err != nil {
		return 400, fmt.Errorf("dash: failed to parse sdkfunc request envelope: %w", err)
	}

	opts := QueryOptions{
		Dataset:    envelope.Req.Input.Dataset,
		DashName:   spec.Id,
		Partitions: envelope.Req.Partitions,
		Params:     envelope.Req.Input.Params,
	}

	rawRet, err := c.QueryRaw(ctx, opts)
	if err != nil {
		return 500, err
	}

	if len(rawRet) == 0 {
		return 200, nil
	}

	// Wrap rawRet into {"resp": {"output": <rawRet>}} and unmarshal back into the typed sdkfunc struct
	wrapper := map[string]any{
		"resp": map[string]json.RawMessage{
			"output": rawRet,
		},
	}
	wrappedBytes, err := json.Marshal(wrapper)
	if err != nil {
		return 500, fmt.Errorf("dash: failed to wrap response: %w", err)
	}

	if err := json.Unmarshal(wrappedBytes, fn); err != nil {
		return 500, fmt.Errorf("dash: failed to populate sdkfunc response: %w", err)
	}

	return 200, nil
}
