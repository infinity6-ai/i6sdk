package relevancefashionresultz

import (
	"fmt"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketbpathz"
)

type CatalogSimple struct {
	StoreId string `json:"store_id" bigquery:"store_id"`
	SkuId   string `json:"sku_id" bigquery:"sku_id"`
	SkuDs   string `json:"sku_ds" bigquery:"sku_ds"`
}

func (c *CatalogSimple) ResultName() string {
	return "catalog-simple"
}

func (c *CatalogSimple) Path() laketbpathz.Path {
	return laketbpathz.MustParse(fmt.Sprintf("i6tb://results.relevance-fashion-catalog/store_id=%s/part.json.gz", c.StoreId))
}

func (c *CatalogSimple) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"store_id": {
					Field: func() any { return &c.StoreId },
				},
				"sku_id": {
					Field: func() any { return &c.SkuId },
				},
				"sku_ds": {
					Field: func() any { return &c.SkuDs },
				},
			}
		},
	}
}
