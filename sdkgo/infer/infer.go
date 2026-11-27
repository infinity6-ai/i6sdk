package infer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/infinity6-ai/gox/commonz/validation/checked"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapimodelget"
	"github.com/infinity6-ai/i6sdk/apiingestz/ingestzapieventevent"
	"github.com/infinity6-ai/i6sdk/apitfz/tfzapimodelpredict"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/tabledataz/relevancefashiondataz"
	"github.com/infinity6-ai/i6sdk/sdkgo/dash"
	"github.com/infinity6-ai/i6sdk/sdkgo/internal/fszcache"
	"github.com/infinity6-ai/i6sdk/sdkgo/internal/sqlitez"
)

// QueryOptions configures model inference queries.
type QueryOptions struct {
	Dataset    string            `json:"dataset"`
	ModelName  string            `json:"model_name"`
	Partitions map[string]string `json:"partitions,omitempty"`
	Params     any               `json:"params,omitempty"`
	SessionId  string            `json:"session_id,omitempty"`
	UserId     string            `json:"user_id,omitempty"`
	Limit      int               `json:"limit,omitempty"`
	SqliteUrl  string            `json:"sqlite_url,omitempty"`
}

// SelectOptions configures result selection and telemetry tracking.
type SelectOptions struct {
	Dataset       string
	ModelName     string
	ModelVersion  string
	StoreId       string
	SessionId     string
	UserId        string
	RecId         string
	Req           any
	Results       any
	SelectResults func(req any, results any) map[string]any
}

// Options configures the Infer client runner.
type Options struct {
	Apis       *apicollectionz.SdkApis
	Downloader *fszcache.Downloader
}

// Client coordinates local candidate generation, remote model prediction, and hydration.
type Client struct {
	apis       *apicollectionz.SdkApis
	downloader *fszcache.Downloader
}

// New creates a new Infer Client.
func New(opts Options) *Client {
	return &Client{
		apis:       opts.Apis,
		downloader: opts.Downloader,
	}
}

// Downloader returns the underlying Downloader.
func (c *Client) Downloader() *fszcache.Downloader {
	return c.downloader
}

// DefaultSelectResults filters out SKUs already in the cart and deduplicates recommendations.
func DefaultSelectResults(req any, results any) map[string]any {
	resMap, ok := results.(map[string]any)
	if !ok {
		return map[string]any{"results": []any{}}
	}

	rawResults, _ := resMap["results"].([]any)
	cartSkus := make(map[string]struct{})

	if reqMap, ok := req.(map[string]any); ok {
		if r, ok := reqMap["req"].(map[string]any); ok {
			reqMap = r
		}
		if inputMap, ok := reqMap["input"].(map[string]any); ok {
			if skusList, ok := inputMap["skus"].([]any); ok {
				for _, s := range skusList {
					if skuMap, ok := s.(map[string]any); ok {
						if skuId, ok := skuMap["sku_id"].(string); ok && skuId != "" {
							cartSkus[skuId] = struct{}{}
						}
					} else if skuStr, ok := s.(string); ok && skuStr != "" {
						cartSkus[skuStr] = struct{}{}
					}
				}
			}
		}
	}

	var selected []any
	seenSkus := make(map[string]struct{})

	for _, item := range rawResults {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}

		var itemSkus []any
		if skus, ok := itemMap["skus"].([]any); ok {
			itemSkus = skus
		} else if skuId, ok := itemMap["sku_id"]; ok {
			itemSkus = []any{map[string]any{"sku_id": fmt.Sprintf("%v", skuId)}}
		}

		for _, skuObj := range itemSkus {
			skuId := ""
			if sMap, ok := skuObj.(map[string]any); ok {
				if id, ok := sMap["sku_id"].(string); ok {
					skuId = id
				}
			} else if sStr, ok := skuObj.(string); ok {
				skuId = sStr
			}

			if skuId != "" {
				if _, inCart := cartSkus[skuId]; !inCart {
					if _, seen := seenSkus[skuId]; !seen {
						seenSkus[skuId] = struct{}{}
						entry := map[string]any{
							"sku_id": skuId,
							"rating": itemMap["rating"],
							"skus":   []any{map[string]any{"sku_id": skuId}},
						}
						if skus, ok := itemMap["skus"].([]any); ok && len(skus) > 0 {
							entry["skus"] = skus
						}
						if enriched, ok := itemMap["enriched"]; ok {
							entry["enriched"] = enriched
						}
						selected = append(selected, entry)
						break
					}
				}
			}
		}
	}

	return map[string]any{"results": selected}
}

// ResolveSqliteUrl resolves the presigned download URL for the model's sqlite.bin file.
func (c *Client) ResolveSqliteUrl(ctx context.Context, opts QueryOptions) (string, error) {
	if opts.SqliteUrl != "" {
		return opts.SqliteUrl, nil
	}

	if opts.Dataset == "" {
		return "", fmt.Errorf("infer: dataset is required to resolve sqlite url")
	}
	if opts.ModelName == "" {
		return "", fmt.Errorf("infer: model_name is required to resolve sqlite url")
	}

	if c.apis == nil {
		return "", fmt.Errorf("infer: cannot resolve model sqlite url: apis client not wired")
	}

	var ds lakedsz.Dataset
	if err := checked.Set(&ds, opts.Dataset); err != nil {
		return "", fmt.Errorf("infer: invalid dataset %q: %w", opts.Dataset, err)
	}

	modelVer := ""
	if opts.Partitions != nil {
		modelVer = opts.Partitions["version"]
	}

	api := &dszapimodelget.Api{
		Req: &dszapimodelget.Req{
			Params: &dszapimodelget.ReqParams{
				Dataset:   ds,
				ModelName: opts.ModelName,
			},
			Payload: &dszapimodelget.ReqPayload{
				ModelVersion:   modelVer,
				ModelPartitions: opts.Partitions,
			},
		},
	}

	status, err := c.apis.DszModelGet(ctx, api)
	if err != nil || status < 200 || status >= 300 || api.Resp == nil || api.Resp.Payload == nil {
		return "", fmt.Errorf("infer: failed to get model metadata for %s (status %d): %w", opts.ModelName, status, err)
	}

	for _, file := range api.Resp.Payload.Urls {
		if file.Name == "sqlite.bin" && file.Url != "" {
			return file.Url, nil
		}
	}

	return "", fmt.Errorf("infer: model %s has no sqlite.bin file available in metadata", opts.ModelName)
}

// Update preloads or revalidates the model SQLite bundle in the local cache.
func (c *Client) Update(ctx context.Context, opts QueryOptions) (string, error) {
	if opts.ModelName == "" {
		return "", fmt.Errorf("infer: model_name is required")
	}

	pkey := dash.PartitionsKey(opts.Partitions)
	cacheName := fmt.Sprintf("infer-%s-%s-sqlite", opts.ModelName, pkey)

	sqliteUrl, err := c.ResolveSqliteUrl(ctx, opts)
	if err != nil {
		return "", err
	}

	return c.downloader.Update(ctx, cacheName, sqliteUrl, 300.0)
}

// Run executes the 3-step inference pipeline:
// Step 1: Local SQLite query 'model-{name}-input' generates candidate tensors.
// Step 2: Server-side TfzModelPredict evaluates candidates.
// Step 3: Local SQLite query 'model-{name}-output' hydrates and ranks results.
func (c *Client) Run(ctx context.Context, opts QueryOptions) (map[string]any, error) {
	if opts.ModelName == "" {
		return nil, fmt.Errorf("infer: model_name is required")
	}
	if opts.Dataset == "" {
		return nil, fmt.Errorf("infer: dataset is required")
	}
	if opts.SessionId == "" || opts.UserId == "" {
		return nil, fmt.Errorf("infer: session_id and user_id are required")
	}
	if opts.Params == nil {
		return nil, fmt.Errorf("infer: params is required")
	}

	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}

	// 1. Download or get cached model SQLite database
	sqlitePath, err := c.Update(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("infer: failed to update model cache: %w", err)
	}

	conn, err := sqlitez.Open(sqlitePath)
	if err != nil {
		return nil, fmt.Errorf("infer: failed to open model sqlite db: %w", err)
	}
	defer conn.Close()

	// 2. Run Step 1: Input query
	queryInput := fmt.Sprintf("model-%s-input", opts.ModelName)
	inputsStr, err := conn.NamedStringQuery(ctx, queryInput, opts.Params)
	if err != nil || inputsStr == nil || *inputsStr == "" {
		return nil, fmt.Errorf("infer: input query %q returned no inputs: %w", queryInput, err)
	}

	// 3. Run Step 2: Call server-side TensorFlow predict API
	if c.apis == nil {
		return nil, fmt.Errorf("infer: apis client not wired for tfzModelPredict")
	}

	predictVersion := ""
	if opts.Partitions != nil {
		predictVersion = opts.Partitions["version"]
	}
	if predictVersion == "" {
		return nil, fmt.Errorf("infer: partitions['version'] is required")
	}

	predictApi := &tfzapimodelpredict.Api{
		Req: &tfzapimodelpredict.Req{
			Params: &tfzapimodelpredict.ReqParams{
				Dataset:   opts.Dataset,
				ModelName: opts.ModelName,
				Version:   predictVersion,
			},
			Payload: &tfzapimodelpredict.ReqPayload{
				Inputs:          *inputsStr,
				ModelPartitions: opts.Partitions,
			},
		},
	}

	status, err := c.apis.TfzModelPredict(ctx, predictApi)
	if err != nil || status < 200 || status >= 300 || predictApi.Resp == nil || predictApi.Resp.Payload == nil {
		return nil, fmt.Errorf("infer: model prediction API failed (status %d): %w", status, err)
	}

	var modelOutputs any
	rawOutputs := predictApi.Resp.Payload.Outputs
	if err := json.Unmarshal([]byte(rawOutputs), &modelOutputs); err != nil {
		modelOutputs = rawOutputs
	}

	// 4. Run Step 3: Output query
	queryOutput := fmt.Sprintf("model-%s-output", opts.ModelName)
	rawResult, err := conn.NamedJsonQuery(ctx, queryOutput, modelOutputs, limit)
	if err != nil {
		return nil, fmt.Errorf("infer: output query %q failed: %w", queryOutput, err)
	}

	var results map[string]any
	if len(rawResult) > 0 {
		if err := json.Unmarshal(rawResult, &results); err != nil {
			return nil, fmt.Errorf("infer: failed to unmarshal output query results: %w", err)
		}
	}
	if results == nil {
		results = make(map[string]any)
	}

	recId := uuid.NewString()
	results["rec_id"] = recId

	// 5. Send recraw telemetry event (non-fatal)
	storeId := "1"
	if s, ok := opts.Partitions["store_id"]; ok && s != "" {
		storeId = s
	} else if s, ok := opts.Partitions["storeId"]; ok && s != "" {
		storeId = s
	}

	c.sendRecEvent(
		ctx,
		opts.Dataset,
		opts.ModelName,
		predictVersion,
		storeId,
		opts.SessionId,
		opts.UserId,
		recId,
		"recraw",
		results,
	)

	return results, nil
}

// SelectResults applies recommendation selection, dispatches recshow telemetry, and returns final recommendations.
func (c *Client) SelectResults(ctx context.Context, opts SelectOptions) (map[string]any, error) {
	if opts.RecId == "" {
		return nil, fmt.Errorf("infer: rec_id is required")
	}
	if opts.SessionId == "" || opts.UserId == "" {
		return nil, fmt.Errorf("infer: session_id and user_id are required")
	}

	callback := opts.SelectResults
	if callback == nil {
		callback = DefaultSelectResults
	}

	selected := callback(opts.Req, opts.Results)
	rawResults, _ := selected["results"].([]any)

	// Collect shown SKUs for recshow event
	var skusList []any
	for _, item := range rawResults {
		if itemMap, ok := item.(map[string]any); ok {
			if skuId, ok := itemMap["sku_id"]; ok && skuId != nil {
				skusList = append(skusList, map[string]any{"sku_id": fmt.Sprintf("%v", skuId)})
			}
		}
	}

	shown := map[string]any{"skus": skusList}

	storeId := opts.StoreId
	if storeId == "" {
		storeId = "1"
	}

	c.sendRecEvent(
		ctx,
		opts.Dataset,
		opts.ModelName,
		opts.ModelVersion,
		storeId,
		opts.SessionId,
		opts.UserId,
		opts.RecId,
		"recshow",
		shown,
	)

	return selected, nil
}

// Execute runs an inference workflow function through the 3-step pipeline.
func (c *Client) Execute(ctx context.Context, fn sdkfuncz.Func) (int, error) {
	spec := fn.SdkFuncSpec()
	if spec.Id == "" {
		return 400, fmt.Errorf("infer: sdkfunc spec id is empty")
	}

	modelName := strings.TrimPrefix(spec.Id, "product-")

	rawFn, err := json.Marshal(fn)
	if err != nil {
		return 400, fmt.Errorf("infer: failed to marshal sdkfunc request: %w", err)
	}

	var envelope struct {
		Req struct {
			Input struct {
				Dataset   string `json:"dataset"`
				SessionId string `json:"session_id"`
				UserId    string `json:"user_id"`
				Skus      any    `json:"skus"`
				Params    any    `json:"params"`
			} `json:"input"`
			Partitions map[string]string `json:"partitions"`
		} `json:"req"`
	}
	if err := json.Unmarshal(rawFn, &envelope); err != nil {
		return 400, fmt.Errorf("infer: failed to parse sdkfunc request envelope: %w", err)
	}

	params := envelope.Req.Input.Params
	if params == nil {
		params = envelope.Req.Input.Skus
	}

	partitions := envelope.Req.Partitions
	if partitions == nil {
		partitions = make(map[string]string)
	}
	if partitions["version"] == "" && spec.Version != "" {
		partitions["version"] = spec.Version
	}

	opts := QueryOptions{
		Dataset:    envelope.Req.Input.Dataset,
		ModelName:  modelName,
		Partitions: partitions,
		Params:     params,
		SessionId:  envelope.Req.Input.SessionId,
		UserId:     envelope.Req.Input.UserId,
	}

	results, err := c.Run(ctx, opts)
	if err != nil {
		return 500, err
	}

	recId, _ := results["rec_id"].(string)

	storeId := "1"
	if s, ok := partitions["store_id"]; ok && s != "" {
		storeId = s
	}

	var reqGeneric any
	if err := json.Unmarshal(rawFn, &reqGeneric); err != nil {
		return 500, fmt.Errorf("infer: failed to parse request payload: %w", err)
	}
	var reqPayload any = reqGeneric
	if m, ok := reqGeneric.(map[string]any); ok {
		if r, ok := m["req"]; ok {
			reqPayload = r
		}
	}

	selected, err := c.SelectResults(ctx, SelectOptions{
		Dataset:      envelope.Req.Input.Dataset,
		ModelName:    modelName,
		ModelVersion: partitions["version"],
		StoreId:      storeId,
		SessionId:    envelope.Req.Input.SessionId,
		UserId:       envelope.Req.Input.UserId,
		RecId:        recId,
		Req:          reqPayload,
		Results:      results,
	})
	if err != nil {
		return 500, fmt.Errorf("infer: select_results failed: %w", err)
	}

	wrapper := map[string]any{
		"resp": map[string]any{
			"output": selected,
		},
	}
	wrappedBytes, err := json.Marshal(wrapper)
	if err != nil {
		return 500, fmt.Errorf("infer: failed to wrap response: %w", err)
	}

	if err := json.Unmarshal(wrappedBytes, fn); err != nil {
		return 500, fmt.Errorf("infer: failed to populate sdkfunc response: %w", err)
	}

	return 200, nil
}

func (c *Client) sendRecEvent(
	ctx context.Context,
	dataset string,
	modelName string,
	modelVersion string,
	storeId string,
	sessionId string,
	userId string,
	recId string,
	kind string,
	payload any,
) {
	if c.apis == nil {
		return
	}

	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return
	}

	evt := &relevancefashiondataz.Event{
		ID:        uuid.NewString(),
		StoreId:   storeId,
		SessionId: sessionId,
		CreatedAt: time.Now().UTC(),
	}
	if userId != "" {
		evt.UserId = &userId
	}

	if kind == "recraw" {
		evt.Recraw = &relevancefashiondataz.Recraw{
			RecId:        recId,
			ModelId:      modelName,
			ModelVersion: modelVersion,
			Payload:      string(rawPayload),
		}
	} else if kind == "recshow" {
		evt.Recshow = &relevancefashiondataz.Recshow{
			RecId:        recId,
			ModelId:      modelName,
			ModelVersion: modelVersion,
			Payload:      string(rawPayload),
		}
	}

	var ds lakedsz.Dataset
	if err := checked.Set(&ds, dataset); err != nil {
		return
	}

	api := &ingestzapieventevent.Api{
		Req: &ingestzapieventevent.Req{
			Params: &ingestzapieventevent.ReqParams{
				Dataset: ds,
			},
			Payload: evt,
		},
	}

	_, _ = c.apis.IngestzEventRelevanceFashionEvent(ctx, api)
}
