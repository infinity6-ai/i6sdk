package dsztoken

import (
	"reflect"

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

func (d *DatasetToken) JSTypedef() (*jstypez.Typedef, *jstypez.Typedefs) {
	t := reflect.TypeOf(*d).String()
	td := &jstypez.Typedef{
		Head:       &jstypez.Tuple{Typ: jstypez.TypObject, GoTyp: t, Name: ""},
		Properties: &jstypez.Tuples{},
	}
	td.Properties.Add(
		jstypez.Tuple{Typ: jstypez.TypString, GoTyp: reflect.TypeOf(d.Id).String(), Name: "id"},
		jstypez.Tuple{Typ: jstypez.TypString, GoTyp: reflect.TypeOf(d.Jti).String(), Name: "jti"},
		jstypez.Tuple{Typ: jstypez.TypString, GoTyp: reflect.TypeOf(d.Dataset).String(), Name: "dataset"},
		jstypez.Tuple{Typ: jstypez.TypString, GoTyp: reflect.TypeOf(d.Annotation).String(), Name: "annotation"},
		jstypez.Tuple{Typ: jstypez.TypNumber, GoTyp: reflect.TypeOf(d.ExpirationDate).String(), Name: "expiration_date"},
		jstypez.Tuple{Typ: jstypez.TypString, GoTyp: reflect.TypeOf(d.Audience).String(), Name: "audience"},
		jstypez.Tuple{Typ: "Array<string>", GoTyp: reflect.TypeOf(d.Grants).String(), Name: "grants"},
		jstypez.Tuple{Typ: jstypez.TypString, GoTyp: reflect.TypeOf(d.Channel).String(), Name: "channel"},
		jstypez.Tuple{Typ: jstypez.TypString, GoTyp: reflect.TypeOf(d.Token).String(), Name: "token"},
	)
	return td, nil
}
