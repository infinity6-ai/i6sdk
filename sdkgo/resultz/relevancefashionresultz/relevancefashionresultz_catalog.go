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

type CatalogSimpleQuery struct {
	StoreId string `json:"store_id"`
}

func (c *CatalogSimple) NewQuery() any {
	return &CatalogSimpleQuery{}
}

func (c *CatalogSimple) QuerySchema(q any) *schemaz.Schema {
	query := q.(*CatalogSimpleQuery)
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"store_id": {
					Desc:     &schemaz.Desc{Summary: "Store ID to query."},
					UrlValue: schemaz.ParseStr(&query.StoreId),
				},
			}
		},
	}
}

func (c *CatalogSimple) Path(q any) laketbpathz.Path {
	query := q.(*CatalogSimpleQuery)
	return laketbpathz.MustParse(fmt.Sprintf("i6tb://results.relevance-fashion-catalog/store_id=%s/part.json.gz", query.StoreId))
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
