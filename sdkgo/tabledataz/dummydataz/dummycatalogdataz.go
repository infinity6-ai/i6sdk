package dummydataz

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
	Id          string              `json:"id" bigquery:"id"`
	Dataset     bigquery.NullString `json:"dataset" bigquery:"dataset"`
	Product     string              `json:"product" bigquery:"product"`
	Price       float64             `json:"price" bigquery:"price"`
	QtInventory int64               `json:"qt_inventory" bigquery:"qt_inventory"`
	StoreId     string              `json:"store_id" bigquery:"store_id"`
	IngestedAt  string              `json:"ingested_at" bigquery:"ingested_at"`
}

func (c *Catalog) Validate() error {
	if err := validation.StrNotEmpty(c.Id, "id is required"); err != nil {
		return fmt.Errorf("failed to validate id: %w", err)
	}
	if err := validation.StrNotEmpty(c.Product, "product is required"); err != nil {
		return fmt.Errorf("failed to validate product: %w", err)
	}
	if err := validation.StrNotEmpty(c.StoreId, "store_id is required"); err != nil {
		return fmt.Errorf("failed to validate store_id: %w", err)
	}
	if err := validation.GreaterOrEqual(c.Price, 0.0, "price must be non-negative"); err != nil {
		return fmt.Errorf("failed to validate price: %w", err)
	}
	if err := validation.GreaterOrEqual(c.QtInventory, int64(0), "qt_inventory must be non-negative"); err != nil {
		return fmt.Errorf("failed to validate qt_inventory: %w", err)
	}
	return nil
}

func (c *Catalog) GetCreatedAt() *time.Time {
	return nil
}

func (c *Catalog) Table() laketbnamez.TableName {
	return laketbnamez.MustParse("dummy-catalog")
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
			Summary: "Dummy Catalog represents a product in the dummy store. Id, Product, and StoreId are required. Price and QtInventory must be non-negative.",
		},
		Object: func(read bool) map[string]*schemaz.Schema {
			m := map[string]*schemaz.Schema{
				"id": {
					Desc:  &schemaz.Desc{Summary: "Catalog item ID (required)."},
					Field: func() any { return &c.Id },
				},
				"product": {
					Desc:  &schemaz.Desc{Summary: "Product name or identifier (required)."},
					Field: func() any { return &c.Product },
				},
				"price": {
					Desc:  &schemaz.Desc{Summary: "Product price (non-negative)."},
					Field: func() any { return &c.Price },
				},
				"qt_inventory": {
					Desc:  &schemaz.Desc{Summary: "Quantity in inventory (non-negative)."},
					Field: func() any { return &c.QtInventory },
				},
				"store_id": {
					Desc:  &schemaz.Desc{Summary: "Store ID (required)."},
					Field: func() any { return &c.StoreId },
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
