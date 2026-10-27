package relevancefashiondataz

import (
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/infinity6-ai/gox/commonz/validation"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketbnamez"
)

type Catalog struct {
	ID         string              `json:"id" bigquery:"id"`
	Dataset    bigquery.NullString `json:"dataset" bigquery:"dataset"`
	StoreId    string              `json:"store_id" bigquery:"store_id"`
	ProductId  string              `json:"product_id" bigquery:"product_id"`
	SkuId      string              `json:"sku_id" bigquery:"sku_id"`
	SkuDs      string              `json:"sku_ds" bigquery:"sku_ds"`
	Inactive   bool                `json:"inactive" bigquery:"inactive"`
	IngestedAt string              `json:"ingested_at" bigquery:"ingested_at"`
}

func (c *Catalog) Validate() error {
	if err := validation.StrNotEmpty(c.ID, "id is required"); err != nil {
		return fmt.Errorf("failed to validate id: %w", err)
	}
	if err := validation.StrNotEmpty(c.StoreId, "store_id is required"); err != nil {
		return fmt.Errorf("failed to validate store_id: %w", err)
	}
	if err := validation.StrNotEmpty(c.ProductId, "product_id is required"); err != nil {
		return fmt.Errorf("failed to validate product_id: %w", err)
	}
	if err := validation.StrNotEmpty(c.SkuId, "sku_id is required"); err != nil {
		return fmt.Errorf("failed to validate sku_id: %w", err)
	}
	return nil
}

func (c *Catalog) GetCreatedAt() *time.Time {
	return nil
}

func (c *Catalog) Table() laketbnamez.TableName {
	return laketbnamez.MustParse("relevance-fashion-catalog")
}

func (c *Catalog) SetDataset(ds lakedsz.Dataset) {
	c.Dataset = bigquery.NullString{
		StringVal: ds.String(),
		Valid:     true,
	}
}

func (c *Catalog) SetIngestedAt(p string) {
	c.IngestedAt = p
}

func (c *Catalog) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Desc: &schemaz.Desc{
			Summary: "Catalog represents a product SKU in the store. ID, StoreId, ProductId, and SkuId are required.",
		},
		Object: func(read bool) map[string]*schemaz.Schema {
			m := map[string]*schemaz.Schema{
				"id": {
					Desc:  &schemaz.Desc{Summary: "Catalog item ID (required)."},
					Field: func() any { return &c.ID },
				},
				"store_id": {
					Desc:  &schemaz.Desc{Summary: "Store ID (required)."},
					Field: func() any { return &c.StoreId },
				},
				"product_id": {
					Desc:  &schemaz.Desc{Summary: "Product ID (required)."},
					Field: func() any { return &c.ProductId },
				},
				"sku_id": {
					Desc:  &schemaz.Desc{Summary: "SKU ID (required)."},
					Field: func() any { return &c.SkuId },
				},
				"sku_ds": {
					Desc:  &schemaz.Desc{Summary: "Optional SKU description."},
					Field: func() any { return &c.SkuDs },
				},
				"inactive": {
					Desc:  &schemaz.Desc{Summary: "Indicates if the SKU is inactive."},
					Field: func() any { return &c.Inactive },
				},
			}
			if read {
				m["dataset"] = &schemaz.Schema{
					Desc:  &schemaz.Desc{Summary: "Dataset name (read-only)."},
					Field: func() any { return &c.Dataset },
				}
				m["ingested_at"] = &schemaz.Schema{
					Desc:  &schemaz.Desc{Summary: "Ingestion timestamp (read-only)."},
					Field: func() any { return &c.IngestedAt },
				}
			}
			return m
		},
	}
}
