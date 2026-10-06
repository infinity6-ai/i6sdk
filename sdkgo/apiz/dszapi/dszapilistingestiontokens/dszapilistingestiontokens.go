package dszapilistingestiontokens

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

type DatasetToken struct {
	Id             string   `json:"id"`
	Jti            string   `json:"jti"`
	Dataset        string   `json:"dataset"`
	Annotation     string   `json:"annotation"`
	ExpirationDate int64    `json:"expiration_date"`
	Audience       string   `json:"audience"`
	Grants         []string `json:"grants"`
	Channel        string   `json:"channel,omitempty"`
	Token          string   `json:"token,omitempty"`
}

func (*DatasetToken) JSTypedef() (*jstypez.Typedef, *jstypez.Typedefs) {
	td := &jstypez.Typedef{
		Head:       &jstypez.Tuple{Typ: jstypez.TypObject, Name: ""},
		Properties: &jstypez.Tuples{},
	}
	td.Properties.Add(
		jstypez.Tuple{Typ: jstypez.TypString, Name: "id"},
		jstypez.Tuple{Typ: jstypez.TypString, Name: "jti"},
		jstypez.Tuple{Typ: jstypez.TypString, Name: "dataset"},
		jstypez.Tuple{Typ: jstypez.TypString, Name: "annotation"},
		jstypez.Tuple{Typ: jstypez.TypNumber, Name: "expiration_date"},
		jstypez.Tuple{Typ: jstypez.TypString, Name: "audience"},
		jstypez.Tuple{Typ: "Array<string>", Name: "grants"},
		jstypez.Tuple{Typ: jstypez.TypString, Name: "channel"},
		jstypez.Tuple{Typ: jstypez.TypString, Name: "token"},
	)
	return td, nil
}

type DszListIngestionTokensReq struct {
	Dataset string `json:"dataset"`
}

type DszListIngestionTokensResp []DatasetToken

func (*DszListIngestionTokensResp) JSType() string {
	return jstypez.TypObject
}

type DszListIngestionTokensApi struct {
	Req  DszListIngestionTokensReq  `json:"req"`
	Resp DszListIngestionTokensResp `json:"resp"`
}

func (d *DszListIngestionTokensApi) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-list-ingestion-tokens",
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/ingestiontoken/list",
	}
}

func (d *DszListIngestionTokensApi) GetDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {Str: schemaz.ParseStr(&d.Req.Dataset)},
				}
			},
		},
		QueryParams: nil,
		ReqHeaders:  nil,
		ReqBody:     nil,
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Raw: func() any { return &d.Resp },
		},
	}
}
