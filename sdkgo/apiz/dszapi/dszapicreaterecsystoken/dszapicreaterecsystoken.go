package dszapicreaterecsystoken

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

type DszCreateRecsysTokenReq struct {
	Dataset    string `json:"dataset"`
	Annotation string `json:"annotation"`
	Channel    string `json:"channel"`
}

type DszCreateRecsysTokenApi struct {
	Req  *DszCreateRecsysTokenReq `json:"req"`
	Resp *DatasetToken            `json:"resp"`
}

func (d *DszCreateRecsysTokenApi) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "dsz-create-recsys-token",
		Method: "POST",
		Path:   "/api/datasetmanager/dataset/recsystoken/create",
	}
}

func (d *DszCreateRecsysTokenApi) GetDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams:  nil,
		QueryParams: nil,
		ReqHeaders:  nil,
		ReqBody: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				if read && d.Req == nil {
					return nil
				}
				if d.Req == nil {
					d.Req = &DszCreateRecsysTokenReq{}
				}
				return map[string]*schemaz.Schema{
					"dataset":    {Str: schemaz.ParseStr(&d.Req.Dataset)},
					"annotation": {Str: schemaz.ParseStr(&d.Req.Annotation)},
					"channel":    {Str: schemaz.ParseStr(&d.Req.Channel)},
				}
			},
		},
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Raw: func() any { return &d.Resp },
		},
	}
}
