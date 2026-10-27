package tabledataz

import (
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketbnamez"
)

type TableData interface {
	Schema() *schemaz.Schema
	Table() laketbnamez.TableName
	SetDataset(dataset lakedsz.Dataset)
	SetIngestedAt(ingestedAt string)
	Validate() error
	GetCreatedAt() *time.Time
}

func NewBqNullString(val string) bigquery.NullString {
	return bigquery.NullString{
		StringVal: val,
		Valid:     true,
	}
}
