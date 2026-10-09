package dummydataz

import (
	"cloud.google.com/go/bigquery"
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
	P           string              `json:"p" bigquery:"p"`
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

func (c *Catalog) SetP(p string) {
	c.P = p
}

func (c *Catalog) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"id": {
					Raw: func() any { return &c.Id },
				},
				"product": {
					Raw: func() any { return &c.Product },
				},
				"price": {
					Raw: func() any { return &c.Price },
				},
				"qt_inventory": {
					Raw: func() any { return &c.QtInventory },
				},
				"store_id": {
					Raw: func() any { return &c.StoreId },
				},
			}
		},
	}
}
