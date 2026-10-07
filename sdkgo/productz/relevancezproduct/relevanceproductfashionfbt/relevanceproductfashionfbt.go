package relevanceproductfashionfbt

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type Sku struct {
	SkuId string `json:"sku_id"`
}

type Result struct {
	SkuId  string `json:"sku_id"`
	Rating string `json:"rating"`
}

type Req struct {
	Skus []*Sku `json:"skus"`
}

type Resp struct {
	Results []*Result `json:"resuts"`
}

type Api struct {
	Req  *Req  `json:"req"`
	Resp *Resp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "product-relevance-fashion-fbt",
		Desc:   nil,
		Method: "GET",
		Path:   "/api/remover/esse-treco",
	}
}

func (d *Api) GetDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		ReqBody:  &schemaz.Schema{},
		RespBody: &schemaz.Schema{},
	}
}
