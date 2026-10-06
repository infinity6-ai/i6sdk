package dsztoken

import (
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

type DatasetToken struct {
	Id             string          `json:"id"`
	Jti            string          `json:"jti"`
	Dataset        lakedsz.Dataset `json:"dataset"`
	Annotation     string          `json:"annotation"`
	ExpirationDate int64           `json:"expiration_date"`
	Audience       string          `json:"audience"`
	Grants         []string        `json:"grants"`
	Channel        string          `json:"channel,omitempty"`
	Token          string          `json:"token,omitempty"`
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
