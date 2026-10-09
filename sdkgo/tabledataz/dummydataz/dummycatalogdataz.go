package dummydataz

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketbnamez"
)

type Catalog struct {
	Id          string  `json:"id" bigquery:"id"`
	Product     string  `json:"product" bigquery:"product"`
	Price       float64 `json:"price" bigquery:"price"`
	QtInventory int64   `json:"qt_inventory" bigquery:"qt_inventory"`
	StoreId     string  `json:"store_id" bigquery:"store_id"`
}

func (u *Catalog) Table() laketbnamez.TableName {
	return laketbnamez.MustParse("dummy-catalog")
}

func (u *Catalog) Schema() *schemaz.Schema {
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
				"store_id": {
					Raw: func() any { return &u.StoreId },
				},
			}
		},
	}
}
