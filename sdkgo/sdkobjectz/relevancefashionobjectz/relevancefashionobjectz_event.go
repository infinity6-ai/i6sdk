package relevancefashionobjectz

import "time"

type Event struct {
	Id             string          `json:"id" bigquery:"id"`
	StoreId        string          `json:"store_id" bigquery:"store_id"`
	SessionId      string          `json:"session_id" bigquery:"session_id"`
	UserId         string          `json:"user_id" bigquery:"user_id"`
	Click          *Click          `json:"click,omitempty" bigquery:"click,nullable"`
	AddToCart      *AddToCart      `json:"add_to_cart,omitempty" bigquery:"add_to_cart,nullable"`
	RemoveFromCart *RemoveFromCart `json:"remove_from_cart,omitempty" bigquery:"remove_from_cart,nullable"`
	Search         *Search         `json:"search,omitempty" bigquery:"search,nullable"`
	Purchase       *Purchase       `json:"purchase,omitempty" bigquery:"purchase,nullable"`
	Recraw         *Recraw         `json:"recraw,omitempty" bigquery:"recraw,nullable"`
	Recshow        *Recshow        `json:"recshow,omitempty" bigquery:"recshow,nullable"`
	CreatedAt      time.Time       `json:"created_at" bigquery:"created_at"`
}

type Click struct {
	Price float64 `json:"price" bigquery:"price"`
	SkuId string  `json:"sku_id" bigquery:"sku_id"`
}

type AddToCart struct {
	Price float64 `json:"price" bigquery:"price"`
	SkuId string  `json:"sku_id" bigquery:"sku_id"`
}

type RemoveFromCart struct {
	Price float64 `json:"price" bigquery:"price"`
	SkuId string  `json:"sku_id" bigquery:"sku_id"`
}

type Search struct {
	Text string `json:"text" bigquery:"text"`
}

type Purchase struct {
	OrderId string         `json:"order_id" bigquery:"order_id"`
	Skus    []*PurchaseSku `json:"skus" bigquery:"skus"`
}

type PurchaseSku struct {
	Qtty      int64   `json:"qtty" bigquery:"qtty"`
	SkuId     string  `json:"sku_id" bigquery:"sku_id"`
	UnitPrice float64 `json:"unit_price" bigquery:"unit_price"`
}

type Recraw struct {
	RecId   string `json:"rec_id" bigquery:"rec_id"`
	ModelId string `json:"model_id" bigquery:"model_id"`
	Payload string `json:"payload" bigquery:"payload"`
}

type Recshow struct {
	RecId   string `json:"rec_id" bigquery:"rec_id"`
	Payload string `json:"payload" bigquery:"payload"`
}
