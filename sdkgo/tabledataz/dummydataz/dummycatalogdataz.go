package dummydataz

import (
	"fmt"

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
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"id": {
					Field: func() any { return &c.Id },
				},
				"product": {
					Field: func() any { return &c.Product },
				},
				"price": {
					Field: func() any { return &c.Price },
				},
				"qt_inventory": {
					Field: func() any { return &c.QtInventory },
				},
				"store_id": {
					Field: func() any { return &c.StoreId },
				},
			}
		},
	}
}
