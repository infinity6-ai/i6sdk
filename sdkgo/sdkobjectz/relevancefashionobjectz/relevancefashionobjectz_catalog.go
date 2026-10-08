package relevancefashionobjectz

type Catalog struct {
	Id        string `json:"id" bigquery:"id"`
	StoreId   string `json:"store_id" bigquery:"store_id"`
	ProductId string `json:"product_id" bigquery:"product_id"`
	SkuId     string `json:"sku_id" bigquery:"sku_id"`
	SkuDs     string `json:"sku_ds" bigquery:"sku_ds"`
	Inactive  bool   `json:"inactive" bigquery:"inactive"`
}
