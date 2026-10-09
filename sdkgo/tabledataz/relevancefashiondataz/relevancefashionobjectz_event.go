package relevancefashiondataz

import (
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketbnamez"
)

type Event struct {
	ID             string              `json:"id" bigquery:"id"`
	Dataset        bigquery.NullString `json:"dataset" bigquery:"dataset"`
	StoreId        string              `json:"store_id" bigquery:"store_id"`
	SessionId      string              `json:"session_id" bigquery:"session_id"`
	UserId         string              `json:"user_id" bigquery:"user_id"`
	Click          *Click              `json:"click,omitempty" bigquery:"click,nullable"`
	AddToCart      *AddToCart          `json:"add_to_cart,omitempty" bigquery:"add_to_cart,nullable"`
	RemoveFromCart *RemoveFromCart     `json:"remove_from_cart,omitempty" bigquery:"remove_from_cart,nullable"`
	Search         *Search             `json:"search,omitempty" bigquery:"search,nullable"`
	Purchase       *Purchase           `json:"purchase,omitempty" bigquery:"purchase,nullable"`
	Recraw         *Recraw             `json:"recraw,omitempty" bigquery:"recraw,nullable"`
	Recshow        *Recshow            `json:"recshow,omitempty" bigquery:"recshow,nullable"`
	CreatedAt      time.Time           `json:"created_at" bigquery:"created_at"`
	P              string              `json:"p" bigquery:"p"`
}

func (e *Event) Table() laketbnamez.TableName {
	return laketbnamez.MustParse("relevance-fashion-event")
}

func (e *Event) SetDataset(ds lakedsz.Dataset) {
	e.Dataset = bigquery.NullString{
		StringVal: ds.String(),
		Valid:     true,
	}
}

func (e *Event) SetP(p string) {
	e.P = p
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

func (u *Click) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"price": {
					Raw: func() any { return &u.Price },
				},
				"sku_id": {
					Raw: func() any { return &u.SkuId },
				},
			}
		},
	}
}

func (u *AddToCart) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"price": {
					Raw: func() any { return &u.Price },
				},
				"sku_id": {
					Raw: func() any { return &u.SkuId },
				},
			}
		},
	}
}

func (u *RemoveFromCart) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"price": {
					Raw: func() any { return &u.Price },
				},
				"sku_id": {
					Raw: func() any { return &u.SkuId },
				},
			}
		},
	}
}

func (u *Search) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"text": {
					Raw: func() any { return &u.Text },
				},
			}
		},
	}
}

func (u *PurchaseSku) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"qtty": {
					Raw: func() any { return &u.Qtty },
				},
				"sku_id": {
					Raw: func() any { return &u.SkuId },
				},
				"unit_price": {
					Raw: func() any { return &u.UnitPrice },
				},
			}
		},
	}
}

func (u *Purchase) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"order_id": {
					Raw: func() any { return &u.OrderId },
				},
				"skus": {
					Array: schemaz.NewArray(schemaz.ArrayOptions[[]*PurchaseSku, *PurchaseSku]{
						NewElement: func(idx int) *PurchaseSku { return &PurchaseSku{} },
						Slice:      func() *[]*PurchaseSku { return &u.Skus },
						Get:        func(idx int, element **PurchaseSku) *schemaz.Schema { return (*element).Schema() },
					}),
				},
			}
		},
	}
}

func (u *Recraw) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"rec_id": {
					Raw: func() any { return &u.RecId },
				},
				"model_id": {
					Raw: func() any { return &u.ModelId },
				},
				"payload": {
					Raw: func() any { return &u.Payload },
				},
			}
		},
	}
}

func (u *Recshow) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"rec_id": {
					Raw: func() any { return &u.RecId },
				},
				"payload": {
					Raw: func() any { return &u.Payload },
				},
			}
		},
	}
}

func (e *Event) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"id": {
					Raw: func() any { return &e.ID },
				},
				"store_id": {
					Raw: func() any { return &e.StoreId },
				},
				"session_id": {
					Raw: func() any { return &e.SessionId },
				},
				"user_id": {
					Raw: func() any { return &e.UserId },
				},
				"click": {
					Object: func(read bool) map[string]*schemaz.Schema {
						if read && e.Click == nil {
							return nil
						}
						if e.Click == nil {
							e.Click = &Click{}
						}
						return e.Click.Schema().Object(read)
					},
				},
				"add_to_cart": {
					Object: func(read bool) map[string]*schemaz.Schema {
						if read && e.AddToCart == nil {
							return nil
						}
						if e.AddToCart == nil {
							e.AddToCart = &AddToCart{}
						}
						return e.AddToCart.Schema().Object(read)
					},
				},
				"remove_from_cart": {
					Object: func(read bool) map[string]*schemaz.Schema {
						if read && e.RemoveFromCart == nil {
							return nil
						}
						if e.RemoveFromCart == nil {
							e.RemoveFromCart = &RemoveFromCart{}
						}
						return e.RemoveFromCart.Schema().Object(read)
					},
				},
				"search": {
					Object: func(read bool) map[string]*schemaz.Schema {
						if read && e.Search == nil {
							return nil
						}
						if e.Search == nil {
							e.Search = &Search{}
						}
						return e.Search.Schema().Object(read)
					},
				},
				"purchase": {
					Object: func(read bool) map[string]*schemaz.Schema {
						if read && e.Purchase == nil {
							return nil
						}
						if e.Purchase == nil {
							e.Purchase = &Purchase{}
						}
						return e.Purchase.Schema().Object(read)
					},
				},
				"recraw": {
					Object: func(read bool) map[string]*schemaz.Schema {
						if read && e.Recraw == nil {
							return nil
						}
						if e.Recraw == nil {
							e.Recraw = &Recraw{}
						}
						return e.Recraw.Schema().Object(read)
					},
				},
				"recshow": {
					Object: func(read bool) map[string]*schemaz.Schema {
						if read && e.Recshow == nil {
							return nil
						}
						if e.Recshow == nil {
							e.Recshow = &Recshow{}
						}
						return e.Recshow.Schema().Object(read)
					},
				},
				"created_at": {
					Raw: func() any { return &e.CreatedAt },
				},
				"p": {
					Raw: func() any { return &e.P },
				},
			}
		},
	}
}
