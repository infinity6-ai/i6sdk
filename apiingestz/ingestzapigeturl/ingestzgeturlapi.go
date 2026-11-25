//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package ingestzapigeturl

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/ingestztoken"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbnamez"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type ReqParams struct {
	// Name of the dataset that receives the uploaded files.
	Dataset lakedsz.Dataset `json:"dataset"`
}

// Request payload containing the target table, partition key/values, and requested URL amount.
type ReqPayload struct {
	// Name of the table, inside the dataset, that receives the uploaded files.
	Table laketbnamez.TableName `json:"table"`
	// Partition key/value pairs that select the table partition the files are ingested into.
	Partitions map[string]string `json:"partitions,omitempty"`
	// Number of presigned upload URLs to generate.
	Amount int `json:"amount"`
}

type Req struct {
	Params  *ReqParams  `json:"params,omitempty"`
	Payload *ReqPayload `json:"payload,omitempty"`
}

// Presigned upload URLs and the ingestion token that identifies this ingestion.
type RespPayload struct {
	// Presigned PUT URLs, one per requested file; upload each file to its URL.
	Uploads []string `json:"uploads"`
	// Ingestion token that identifies this ingestion and ties the uploaded files together.
	IngestId ingestztoken.IngestToken `json:"ingest_id"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Get Ingestion Upload URLs
// Generates presigned PUT URLs and an ingestion token for uploading files into a target dataset table and partitions.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/ingest/dataset/{dataset}/get-url",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type IngestzGetUrlId struct{}

func (IngestzGetUrlId) IngestzGetUrl(ctx context.Context, api *Api) (int, error) {
	panic("IngestzGetUrl not implemented: run 'go generate ./...'")
}
