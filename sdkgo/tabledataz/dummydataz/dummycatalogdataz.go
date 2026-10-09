package dummydataz

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketbnamez"
)

type DummyCatalogData struct {
	Id          string  `json:"id" bigquery:"id"`
	Product     string  `json:"product" bigquery:"product"`
	Price       float64 `json:"price" bigquery:"price"`
	QtInventory int64   `json:"qt_inventory" bigquery:"qt_inventory"`
	Store       string  `json:"store" bigquery:"store"`
}

func (u *DummyCatalogData) Table() laketbnamez.TableName {
	return laketbnamez.MustParse("dummy_catalog_data")
}

func (u *DummyCatalogData) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"id": {
					Raw: func() any { return &u.Id },
				},
				"product": {
					Raw: func() any { return &u.Product },
				},
				"price": {
					Raw: func() any { return &u.Price },
				},
				"qt_inventory": {
					Raw: func() any { return &u.QtInventory },
				},
				"store": {
					Raw: func() any { return &u.Store },
				},
			}
		},
	}
}
