package relevancefashiondataz

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketbnamez"
)

type Catalog struct {
	ID        string           `json:"id" bigquery:"id"`
	Dataset   *lakedsz.Dataset `json:"dataset" bigquery:"dataset,nullable"`
	StoreId   string           `json:"store_id" bigquery:"store_id"`
	ProductId string           `json:"product_id" bigquery:"product_id"`
	SkuId     string           `json:"sku_id" bigquery:"sku_id"`
	SkuDs     string           `json:"sku_ds" bigquery:"sku_ds"`
	Inactive  bool             `json:"inactive" bigquery:"inactive"`
	P         string           `json:"p" bigquery:"p"`
}

func (c *Catalog) Table() laketbnamez.TableName {
	return laketbnamez.MustParse("relevance-fashion-catalog")
}

func (c *Catalog) SetDataset(ds *lakedsz.Dataset) {
	c.Dataset = ds
}

func (c *Catalog) SetP(p string) {
	c.P = p
}

func (c *Catalog) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"id": {
					Raw: func() any { return &c.ID },
				},
				"store_id": {
					Raw: func() any { return &c.StoreId },
				},
				"product_id": {
					Raw: func() any { return &c.ProductId },
				},
				"sku_id": {
					Raw: func() any { return &c.SkuId },
				},
				"sku_ds": {
					Raw: func() any { return &c.SkuDs },
				},
				"inactive": {
					Raw: func() any { return &c.Inactive },
				},
			}
		},
	}
}
