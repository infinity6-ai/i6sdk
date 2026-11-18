package targetingfinancedataz

import (
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/infinity6-ai/gox/commonz/validation"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketbnamez"
	"go.code.infinity6.ai/sdkgo/sdkgo/tabledataz"
)

type Weight struct {
	ID         string              `json:"id" bigquery:"id"`
	Dataset    bigquery.NullString `json:"dataset" bigquery:"dataset"`
	IngestedAt string              `json:"ingested_at" bigquery:"ingested_at"`
	Weight     float64             `json:"weight" bigquery:"weight"`
}

func (w *Weight) Validate() error {
	if err := validation.StrNotEmpty(w.ID, "id is required"); err != nil {
		return fmt.Errorf("failed to validate id: %w", err)
	}
	return nil
}

func (w *Weight) GetCreatedAt() *time.Time {
	return nil
}

func (w *Weight) Table() laketbnamez.TableName {
	return laketbnamez.MustParse("finance-predictivecampaign-weights")
}

func (w *Weight) SetDataset(ds lakedsz.Dataset) {
	w.Dataset = tabledataz.NewBqNullString(ds.Get())
}

func (w *Weight) SetIngestedAt(p string) {
	w.IngestedAt = p
}

func (w *Weight) SchemaQuery() *schemaz.Schema {
	if w == nil {
		w = &Weight{}
	}
	return &schemaz.Schema{
		Desc: &schemaz.Desc{Summary: "Weight is a row of the finance-predictivecampaign-weights table. The id is provided by the sender."},
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"id": {
					Desc:     &schemaz.Desc{Summary: "Unique row identifier, provided by the sender (required)."},
					UrlValue: schemaz.ParseStr(&w.ID),
				},
				"weight": {
					Desc:     &schemaz.Desc{Summary: "weight (required)."},
					UrlValue: schemaz.ParseStrNumber(&w.Weight),
				},
			}
		},
	}
}

func (w *Weight) SchemaBody() *schemaz.Schema {
	if w == nil {
		w = &Weight{}
	}
	return &schemaz.Schema{
		Desc: &schemaz.Desc{Summary: "Weight is a row of the finance-predictivecampaign-weights table. The id is provided by the sender."},
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"id": {
					Desc:  &schemaz.Desc{Summary: "Unique row identifier, provided by the sender (required)."},
					Field: func() any { return &w.ID },
				},
				"weight": {
					Desc:  &schemaz.Desc{Summary: "weight (required)."},
					Field: func() any { return &w.Weight },
				},
			}
		},
	}
}
