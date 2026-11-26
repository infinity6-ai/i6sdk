//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package ingestzapibatchsummary

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbnamez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/storeid"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type ReqParams struct {
	// Name of the dataset that receives the uploaded files.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type ReqQuery struct {
	// Name of the table inside the dataset.
	Table laketbnamez.TableName `json:"table"`
	// Store ID associated with the ingestion.
	StoreId storeid.StoreId `json:"store_id"`
}

type Req struct {
	Params *ReqParams `json:"params,omitempty"`
	Query  *ReqQuery  `json:"query,omitempty"`
}

// Summary of distinct and null counts for a single column.
type ColumnSummary struct {
	// Name of the column.
	ColumnName string `json:"column_name"`
	// Number of distinct values in the column.
	DistinctCount int `json:"distinct_count"`
	// Number of null values in the column.
	NullCount int `json:"null_count"`
}

// Get Last Successful Ingestion Data Summary Response.
type RespPayload struct {
	// Store ID associated with the ingestion.
	StoreId storeid.StoreId `json:"store_id"`
	// Total number of successfully ingested rows.
	RowCount int `json:"row_count"`
	// Timestamp when the batch was successfully ingested.
	IngestedAt string `json:"ingested_at"`
	// Summaries of distinct and null counts per column.
	ColumnSummaries []ColumnSummary `json:"column_summaries"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Get Last Successful Ingestion Data Summary
// Returns the date of ingestion, rows count and distinct and null columns count.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/ingest/dataset/{dataset}/last-batch",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type IngestzLastBatchSummaryId struct{}

func (IngestzLastBatchSummaryId) IngestzLastBatchSummary(ctx context.Context, api *Api) (int, error) {
	panic("IngestzLastBatchSummary not implemented: run 'go generate ./...'")
}
